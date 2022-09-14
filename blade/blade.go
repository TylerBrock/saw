package blade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/TylerBrock/saw/config"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/fatih/color"
)

// A Blade is a Saw execution instance
type Blade struct {
	config *config.Configuration
	aws    *config.AWSConfiguration
	output *config.OutputConfiguration
	cwl    *cloudwatchlogs.Client
}

// NewBlade creates a new Blade with CloudWatchLogs instance from provided config
func NewBlade(
	ctx context.Context,
	config *config.Configuration,
	awsConfig *config.AWSConfiguration,
	outputConfig *config.OutputConfiguration,
) (*Blade, error) {
	blade := Blade{}

	var opts []func(*awsconfig.LoadOptions) error
	if awsConfig.Profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(awsConfig.Profile))
	}
	if awsConfig.Region != "" {
		opts = append(opts, awsconfig.WithRegion(awsConfig.Region))
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)

	blade.cwl = cloudwatchlogs.NewFromConfig(awsCfg)
	blade.config = config
	blade.output = outputConfig

	return &blade, err
}

// GetLogGroups gets the log groups from AWS given the blade configuration
func (b *Blade) GetLogGroups(ctx context.Context) (groups []types.LogGroup, err error) {
	input := b.config.DescribeLogGroupsInput()
	logGroupsPaginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(b.cwl, input)
	var page *cloudwatchlogs.DescribeLogGroupsOutput
	for logGroupsPaginator.HasMorePages() {
		page, err = logGroupsPaginator.NextPage(ctx)
		if err != nil {
			return
		}
		groups = append(groups, page.LogGroups...)
	}
	return
}

func groupNameMatches(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (b *Blade) ResolveFuzzyGroupName(ctx context.Context) (err error) {
	if !b.config.Fuzzy {
		return
	}
	b.config.Fuzzy = false
	groups, err := b.GetLogGroups(ctx)
	if err != nil {
		return
	}
	if len(groups) == 0 {
		return errors.New("no log groups found")
	}
	filtered := filterGroupNames(groups, b.config.Group)
	if len(filtered) > 1 {
		return fmt.Errorf("too many results for log group fuzzy search\n%s", strings.Join(filtered, "\n"))
	}
	if len(filtered) == 0 {
		return fmt.Errorf("no results for log group fuzzy search in %d groups\n%s", len(groups), strings.Join(getGroupNames(groups), "\n"))
	}
	b.config.Group = filtered[0]
	return
}

func getGroupNames(groups []types.LogGroup) (op []string) {
	op = make([]string, len(groups))
	for i := 0; i < len(groups); i++ {
		op[i] = *groups[i].LogGroupName
	}
	return
}

func filterGroupNames(groups []types.LogGroup, group string) (op []string) {
	for i := 0; i < len(groups); i++ {
		if groupNameMatches(*groups[i].LogGroupName, group) {
			op = append(op, *groups[i].LogGroupName)
		}
	}
	return
}

// GetLogStreams gets the log streams from AWS given the blade configuration
func (b *Blade) GetLogStreams(ctx context.Context) (streams []types.LogStream, err error) {
	if err := b.ResolveFuzzyGroupName(ctx); err != nil {
		return nil, err
	}
	input := b.config.DescribeLogStreamsInput()
	logStreamsPaginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(b.cwl, input)
	var page *cloudwatchlogs.DescribeLogStreamsOutput
	for logStreamsPaginator.HasMorePages() {
		page, err = logStreamsPaginator.NextPage(ctx)
		if err != nil {
			return
		}
		streams = append(streams, page.LogStreams...)
	}
	return
}

// GetEvents gets events from AWS given the blade configuration
func (b *Blade) GetEvents(ctx context.Context) (err error) {
	if err := b.ResolveFuzzyGroupName(ctx); err != nil {
		return err
	}
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()
	logEventsPaginator := cloudwatchlogs.NewFilterLogEventsPaginator(b.cwl, input)
	var page *cloudwatchlogs.FilterLogEventsOutput
	for logEventsPaginator.HasMorePages() {
		page, err = logEventsPaginator.NextPage(ctx)
		if err != nil {
			return
		}
		for _, event := range page.Events {
			if b.output.Pretty {
				fmt.Println(strings.TrimRight(formatEvent(formatter, event), "\n"))
			} else {
				fmt.Println(strings.TrimRight(*event.Message, "\n"))
			}
		}
	}
	return
}

// StreamEvents continuously prints log events to the console
func (b *Blade) StreamEvents(ctx context.Context) (err error) {
	var lastSeenTime *int64
	var seenEventIDs map[string]bool
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()

	clearSeenEventIds := func() {
		seenEventIDs = make(map[string]bool, 0)
	}

	addSeenEventIDs := func(id *string) {
		seenEventIDs[*id] = true
	}

	updateLastSeenTime := func(ts *int64) {
		if lastSeenTime == nil || *ts > *lastSeenTime {
			lastSeenTime = ts
			clearSeenEventIds()
		}
	}

	for {
		logEventsPaginator := cloudwatchlogs.NewFilterLogEventsPaginator(b.cwl, input)
		var page *cloudwatchlogs.FilterLogEventsOutput
		for logEventsPaginator.HasMorePages() {
			page, err = logEventsPaginator.NextPage(ctx)
			if err != nil {
				return
			}
			for _, event := range page.Events {
				updateLastSeenTime(event.Timestamp)
				if _, seen := seenEventIDs[*event.EventId]; !seen {
					var message string
					if b.output.Raw {
						message = *event.Message
					} else {
						message = formatEvent(formatter, event)
					}
					message = strings.TrimRight(message, "\n")
					fmt.Println(message)
					addSeenEventIDs(event.EventId)
				}
			}
		}
		if lastSeenTime != nil {
			input.StartTime = lastSeenTime
		}
		time.Sleep(1 * time.Second)
	}
}

// formatEvent returns a CloudWatch log event as a formatted string using the provided formatter
func formatEvent(formatter *colorjson.Formatter, event types.FilteredLogEvent) string {
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	str := *event.Message
	bytes := []byte(*event.Message)
	dateStr := time.UnixMilli(*event.Timestamp).Format(time.RFC3339)
	streamStr := *event.LogStreamName
	jl := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &jl); err != nil {
		return fmt.Sprintf("[%s] (%s) %s", red(dateStr), white(streamStr), str)
	}

	output, _ := formatter.Marshal(jl)
	return fmt.Sprintf("[%s] (%s) %s", red(dateStr), white(streamStr), output)
}
