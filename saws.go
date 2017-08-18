package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

var lastSeenTime *int64
var formatter *colorjson.Formatter

func updateLastSeenTime(ts *int64) {
	if lastSeenTime == nil || *ts > *lastSeenTime {
		lastSeenTime = ts
	}
}

func handlePage(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	for _, event := range page.Events {
		str := aws.StringValue(event.Message)
		bytes := []byte(str)
		date := aws.MillisecondsTimeValue(event.IngestionTime)
		dateStr := date.Format(time.RFC3339)
		streamStr := aws.StringValue(event.LogStreamName)
		jl := map[string]interface{}{}
		if err := json.Unmarshal(bytes, &jl); err != nil {
			fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), str)
		} else {
			output, _ := formatter.Marshal(jl)
			fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), output)
		}
		updateLastSeenTime(event.Timestamp)
	}
	return !lastPage
}

func filterLogStreams(
	cwl *cloudwatchlogs.CloudWatchLogs,
	logGroupName *string,
	logStreamPrefix *string,
) []*string {
	input := cloudwatchlogs.DescribeLogStreamsInput{}
	input.SetLogGroupName(*logGroupName)
	input.SetLogStreamNamePrefix(*logStreamPrefix)

	streams := make([]*cloudwatchlogs.LogStream, 0)
	cwl.DescribeLogStreamsPages(&input, func(
		out *cloudwatchlogs.DescribeLogStreamsOutput,
		lastPage bool,
	) bool {
		for _, stream := range out.LogStreams {
			streams = append(streams, stream)
		}
		return !lastPage
	})

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

func getTime(timeStr *string) (time.Time, error) {
	relative, err := time.ParseDuration(*timeStr)
	if err == nil {
		return time.Now().Add(relative), nil
	}

	absolute, err := time.Parse(time.RFC3339, *timeStr)
	if err == nil {
		return absolute, nil
	}

	return time.Time{}, errors.New("Could not parse relative or absolute time")
}

func configure(cw *cloudwatchlogs.CloudWatchLogs) *cloudwatchlogs.FilterLogEventsInput {
	logGroupName := flag.String("group", "", "Log group to stream")
	logStreamPrefix := flag.String("prefix", "", "Log stream prefix")
	startTime := flag.String("start", "", "Start time")
	endTime := flag.String("end", "", "End time")
	filterPattern := flag.String("filter", "", "Filter pattern")
	flag.Parse()

	if *logGroupName == "" {
		fmt.Println("Error: Must provide a logGroup!")
		os.Exit(1)
	}

	input := cloudwatchlogs.FilterLogEventsInput{}
	input.SetInterleaved(true)
	input.SetLogGroupName(*logGroupName)

	absoluteStartTime := time.Now()
	if *startTime != "" {
		st, err := getTime(startTime)
		if err == nil {
			absoluteStartTime = st
		}
	}
	input.SetStartTime(aws.TimeUnixMilli(absoluteStartTime))

	if *endTime != "" {
		et, err := getTime(startTime)
		if err == nil {
			input.SetEndTime(aws.TimeUnixMilli(et))
		}
	}

	if len(*logStreamPrefix) > 0 {
		streamNamePointers := filterLogStreams(cw, logGroupName, logStreamPrefix)
		if len(streamNamePointers) > 0 {
			input.SetLogStreamNames(streamNamePointers)
		} else {
			fmt.Printf("Cannot find any log streams with prefix \"%s\"\n", *logStreamPrefix)
			os.Exit(3)
		}
	}

	if len(*filterPattern) != 0 {
		input.SetFilterPattern(*filterPattern)
	}

	return &input
}

func configureFormatter() *colorjson.Formatter {
	expand := flag.Bool("expand", false, "Expand JSON log lines")
	//raw := flag.Bool("raw", false, "Disable all color and adornment of log lines")
	rawString := flag.Bool("rawString", false, "Write raw JSON strings")
	invert := flag.Bool("invert-color", false, "Inverts key color from white to black")
	noColor := flag.Bool("no-color", false, "Disable color output")
	flag.Parse()

	var formatter = colorjson.NewFormatter()

	if *expand {
		formatter.Indent = 4
	}

	if *rawString {
		formatter.RawStrings = true
	}

	if *invert {
		formatter.KeyColor = color.New(color.FgWhite)
	}

	if *noColor {
		color.NoColor = true // disables colorized output
	}

	return formatter
}

func main() {
	config := aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}
	sess := session.Must(session.NewSession(&config))
	cwl := cloudwatchlogs.New(sess)
	input := configure(cwl)
	formatter = configureFormatter()

	for {
		err := cwl.FilterLogEventsPages(input, handlePage)
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
