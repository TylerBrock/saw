package blade

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/TylerBrock/colorjson"
	"github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

// A Blade is a Saw execution instance
type Blade struct {
	config *config.Configuration
	aws    *config.AWSConfiguration
	output *config.OutputConfiguration
	cwl    *cloudwatchlogs.CloudWatchLogs
}

// NewBlade creates a new Blade with CloudWatchLogs instance from provided config
func NewBlade(
	config *config.Configuration,
	awsConfig *config.AWSConfiguration,
	outputConfig *config.OutputConfiguration,
) *Blade {
	blade := Blade{}
	awsCfg := aws.Config{}

	if awsConfig.Region != "" {
		awsCfg.Region = &awsConfig.Region
	}

	awsSessionOpts := session.Options{
		Config:                  awsCfg,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
	}

	if awsConfig.Profile != "" {
		awsSessionOpts.Profile = awsConfig.Profile
	}

	sess := session.Must(session.NewSessionWithOptions(awsSessionOpts))

	blade.cwl = cloudwatchlogs.New(sess)
	blade.config = config
	blade.output = outputConfig

	return &blade
}

// GetLogGroups gets the log groups from AWS given the blade configuration
func (b *Blade) GetLogGroups() []*cloudwatchlogs.LogGroup {
	input := b.config.DescribeLogGroupsInput()
	groups := make([]*cloudwatchlogs.LogGroup, 0)
	b.cwl.DescribeLogGroupsPages(input, func(
		out *cloudwatchlogs.DescribeLogGroupsOutput,
		lastPage bool,
	) bool {
		for _, group := range out.LogGroups {
			groups = append(groups, group)
		}
		return !lastPage
	})
	return groups
}

// GetLogStreams gets the log streams from AWS given the blade configuration
func (b *Blade) GetLogStreams() []*cloudwatchlogs.LogStream {
	input := b.config.DescribeLogStreamsInput()
	streams := make([]*cloudwatchlogs.LogStream, 0)
	b.cwl.DescribeLogStreamsPages(input, func(
		out *cloudwatchlogs.DescribeLogStreamsOutput,
		lastPage bool,
	) bool {
		for _, stream := range out.LogStreams {
			streams = append(streams, stream)
		}
		return !lastPage
	})

	return streams
}

// GetEvents gets events from AWS given the blade configuration
func (b *Blade) GetEvents() {
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()

	handlePage := func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range page.Events {
			if b.output.Pretty {
				fmt.Println(formatEvent(formatter, event))
			} else {
				fmt.Println(*event.Message)
			}
		}
		return !lastPage
	}
	err := b.cwl.FilterLogEventsPages(input, handlePage)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(2)
	}
}

// StreamEvents continuously prints log events to the console
func (b *Blade) StreamEvents() {
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

	handlePage := func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
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
		return !lastPage
	}

	for {
		err := b.cwl.FilterLogEventsPages(input, handlePage)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}
		if lastSeenTime != nil {
			input.SetStartTime(*lastSeenTime)
		}
		time.Sleep(1 * time.Second)
	}
}

// RunQuery runs a CloudWatch Logs Insights query and prints the results to the console
func (b *Blade) RunQuery() {
	input := b.config.StartQueryInput()
	startQueryOutput, err := b.cwl.StartQuery(input)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(2)
	}

	queryId := startQueryOutput.QueryId
	queryResultsOutput := &cloudwatchlogs.GetQueryResultsOutput{}
	for {
		queryResultsOutput, err = b.cwl.GetQueryResults(&cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryId,
		})
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(2)
		}

		if *queryResultsOutput.Status == "Complete" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	// We're assuming the field names are always in the same order
	// and that the order is the same as they appear in the query
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)
	table.SetAutoWrapText(false)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetNoWhiteSpace(true)
	table.SetRowSeparator("")
	table.SetTablePadding("\t")

	firstRow := queryResultsOutput.Results[0]
	headers := make([]string, len(firstRow)-1)
	for i := 0; i < len(firstRow) - 1; i++ {
		headers[i] = *firstRow[i].Field
	}
	if !b.output.NoHeaders {
		table.SetHeader(headers)
	}

	logEntries := make([][]string, len(queryResultsOutput.Results))
	for index, logRecord := range queryResultsOutput.Results {
		logEntry := make([]string, len(logRecord)-1)
		for i := 0; i < len(logRecord)-1; i++ {
			logEntry[i] = *logRecord[i].Value
		}
		logEntries[index] = logEntry
	}
	table.AppendBulk(logEntries)
	table.Render()
}

// formatEvent returns a CloudWatch log event as a formatted string using the provided formatter
func formatEvent(formatter *colorjson.Formatter, event *cloudwatchlogs.FilteredLogEvent) string {
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	str := aws.StringValue(event.Message)
	bytes := []byte(str)
	date := aws.MillisecondsTimeValue(event.Timestamp)
	dateStr := date.Format(time.RFC3339)
	streamStr := aws.StringValue(event.LogStreamName)
	jl := map[string]interface{}{}

	if err := json.Unmarshal(bytes, &jl); err != nil {
		return fmt.Sprintf("[%s] (%s) %s", red(dateStr), white(streamStr), str)
	}

	output, _ := formatter.Marshal(jl)
	return fmt.Sprintf("[%s] (%s) %s", red(dateStr), white(streamStr), output)
}
