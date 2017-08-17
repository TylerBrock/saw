package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

var formatter = colorjson.NewFormatter()
var lastSeenTime *int64

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
	cw *cloudwatchlogs.CloudWatchLogs,
	logGroupName *string,
	logStreamPrefix *string,
) []*string {
	input := cloudwatchlogs.DescribeLogStreamsInput{}
	input.SetLogGroupName(*logGroupName)
	input.SetLogStreamNamePrefix(*logStreamPrefix)
	var streamNames = make([]*string, 0)
	cw.DescribeLogStreamsPages(&input, func(out *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range out.LogStreams {
			streamNames = append(streamNames, stream.LogStreamName)
		}
		return lastPage
	})
	return streamNames
}

func main() {
	logGroupName := flag.String("group", "", "Log group to stream")
	logStreamPrefix := flag.String("prefix", "", "Log stream prefix")
	filterPattern := flag.String("filter", "", "Filter pattern")
	expand := flag.Bool("expand", false, "Expand JSON log lines")
	//raw := flag.Bool("raw", false, "Disable all color and adornment of log lines")
	rawString := flag.Bool("rawString", false, "Write raw JSON strings")
	invert := flag.Bool("invert-color", false, "Inverts key color from white to black")
	noColor := flag.Bool("no-color", false, "Disable color output")

	flag.Parse()

	if *expand {
		formatter.Indent = 4
	}

	if *rawString {
		formatter.RawStrings = true
	}

	if *invert {
		formatter.KeyColor = color.New(color.FgWhite)
	}

	if *logGroupName == "" {
		fmt.Println("Error: Must provide a logGroup!")
		os.Exit(1)
	}

	if *noColor {
		color.NoColor = true // disables colorized output
	}

	config := aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}
	sess := session.Must(session.NewSession(&config))
	cw := cloudwatchlogs.New(sess)

	input := cloudwatchlogs.FilterLogEventsInput{}
	input.SetInterleaved(true)
	input.SetLogGroupName(*logGroupName)
	input.SetStartTime(aws.TimeUnixMilli(time.Now()))

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

	for {
		err := cw.FilterLogEventsPages(&input, handlePage)
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
