package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TylerBrock/saws/jdad"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

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
			fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), jdad.Serialize(jl))
		}
	}
	return !lastPage
}

func main() {
	logGroupName := flag.String("group", "", "the log group to stream")
	logStreamPrefix := flag.String("prefix", "", "the log stream prefix")
	filterPattern := flag.String("filter", "", "the filter pattern")
	noColor := flag.Bool("no-color", false, "Disable color output")

	flag.Parse()

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
	input.SetStartTime(aws.TimeUnixMilli(time.Now().Add(-10 * time.Second)))

	if len(*logStreamPrefix) > 0 {
		streams := strings.Split(*logStreamPrefix, ",")
		streamNamePointers := make([]*string, len(streams))
		for i, stream := range streams {
			streamNamePointers[i] = &stream
		}
		input.SetLogStreamNames(streamNamePointers)
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
		time.Sleep(1 * time.Second)
	}
}
