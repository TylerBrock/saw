package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

var formatter = colorjson.NewFormatter()

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
	}
	return !lastPage
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
