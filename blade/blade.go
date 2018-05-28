package blade

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

type Blade struct {
	config *config.Configuration
	output *config.OutputConfiguration
	cwl    *cloudwatchlogs.CloudWatchLogs
}

func NewBlade(config *config.Configuration, outputConfig *config.OutputConfiguration) *Blade {
	blade := Blade{}
	region := endpoints.UsEast1RegionID
	awsConfig := aws.Config{Region: &region}
	sess := session.Must(session.NewSession(&awsConfig))

	blade.cwl = cloudwatchlogs.New(sess)
	blade.config = config
	blade.output = outputConfig

	return &blade
}

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

func (b *Blade) StreamEvents() {
	var lastSeenTime *int64
	formatter := b.output.Formatter()
	input := b.config.FilterLogEventsInput()

	updateLastSeenTime := func(ts *int64) {
		if lastSeenTime == nil || *ts > *lastSeenTime {
			lastSeenTime = ts
		}
	}

	handlePage := func(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
		for _, event := range page.Events {
			printEvent(formatter, event)
			updateLastSeenTime(event.Timestamp)
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

func printEvent(formatter *colorjson.Formatter, event *cloudwatchlogs.FilteredLogEvent) {
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	str := aws.StringValue(event.Message)
	bytes := []byte(str)
	date := aws.MillisecondsTimeValue(event.Timestamp)
	dateStr := date.Format(time.RFC3339)
	streamStr := aws.StringValue(event.LogStreamName)
	jl := map[string]interface{}{}
	if err := json.Unmarshal(bytes, &jl); err != nil {
		fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), str)
	} else {
		output, _ := formatter.Marshal(jl)
		fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), output)
	}
}

func topStreamNames(streams []*cloudwatchlogs.LogStream) []*string {
	// FilerLogEvents can only take 100 streams so lets sort by LastEventTimestamp
	// (descending) and take only the names of the most recent 100.
	sort.Slice(streams, func(i int, j int) bool {
		return *streams[i].LastEventTimestamp > *streams[j].LastEventTimestamp
	})

	var streamNames = make([]*string, 0)

	for _, stream := range streams[:100] {
		fmt.Println(stream.LogStreamName, aws.MillisecondsTimeValue(stream.LastEventTimestamp))
		streamNames = append(streamNames, stream.LogStreamName)
	}

	return streamNames
}
