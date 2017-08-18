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

const version = "v0.0.1"

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
	logGroupName string,
	logStreamPrefix string,
) []*string {
	input := cloudwatchlogs.DescribeLogStreamsInput{}
	input.SetLogGroupName(logGroupName)
	input.SetLogStreamNamePrefix(logStreamPrefix)

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

func getTime(timeStr string) (time.Time, error) {
	relative, err := time.ParseDuration(timeStr)
	if err == nil {
		return time.Now().Add(relative), nil
	}

	absolute, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return absolute, nil
	}

	return time.Time{}, errors.New("Could not parse relative or absolute time")
}

type configuration struct {
	group     string
	prefix    string
	start     string
	end       string
	filter    string
	region    string
	expand    bool
	raw       bool
	rawString bool
	invert    bool
	noColor   bool
}

func configure() *configuration {
	config := configuration{}
	flag.StringVar(&config.group, "group", "", "Log group to stream")
	flag.StringVar(&config.prefix, "prefix", "", "Log stream prefix")
	flag.StringVar(&config.start, "start", "", "Start time")
	flag.StringVar(&config.end, "end", "", "End time")
	flag.StringVar(&config.filter, "filter", "", "Filter pattern")
	flag.StringVar(&config.region, "region", endpoints.UsEast1RegionID, "AWS region")
	flag.BoolVar(&config.expand, "expand", false, "Expand JSON log lines")
	flag.BoolVar(&config.raw, "raw", false, "Disable all color and adornment of log lines")
	flag.BoolVar(&config.rawString, "rawString", false, "Write raw JSON strings")
	flag.BoolVar(&config.invert, "invert", false, "Inverts key color from white to black")
	flag.BoolVar(&config.noColor, "no-color", false, "Disable color output")
	showVersion := flag.Bool("version", false, "Print version string")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return &config
}

func (c *configuration) getFilterEventsInput(cwl *cloudwatchlogs.CloudWatchLogs) *cloudwatchlogs.FilterLogEventsInput {
	if c.group == "" {
		fmt.Println("Error: Must provide a CloudWatchLogs log group with --group!")
		os.Exit(1)
	}

	input := cloudwatchlogs.FilterLogEventsInput{}
	input.SetInterleaved(true)
	input.SetLogGroupName(c.group)

	absoluteStartTime := time.Now()
	if c.start != "" {
		st, err := getTime(c.start)
		if err == nil {
			absoluteStartTime = st
		}
	}
	input.SetStartTime(aws.TimeUnixMilli(absoluteStartTime))

	if c.end != "" {
		et, err := getTime(c.end)
		if err == nil {
			input.SetEndTime(aws.TimeUnixMilli(et))
		}
	}

	if len(c.prefix) > 0 {
		streamNamePointers := filterLogStreams(cwl, c.group, c.prefix)
		if len(streamNamePointers) > 0 {
			input.SetLogStreamNames(streamNamePointers)
		} else {
			fmt.Printf("Cannot find any log streams with prefix \"%s\"\n", c.prefix)
			os.Exit(3)
		}
	}

	if len(c.filter) != 0 {
		input.SetFilterPattern(c.filter)
	}

	return &input
}

func (c *configuration) getFormatter() *colorjson.Formatter {
	var formatter *colorjson.Formatter = colorjson.NewFormatter()

	if c.expand {
		formatter.Indent = 4
	}

	if c.rawString {
		formatter.RawStrings = true
	}

	if c.invert {
		formatter.KeyColor = color.New(color.FgBlack)
	}

	if c.noColor {
		color.NoColor = true
	}

	return formatter
}

func main() {
	config := configure()
	awsConfig := aws.Config{Region: &config.region}
	sess := session.Must(session.NewSession(&awsConfig))
	cwl := cloudwatchlogs.New(sess)
	input := config.getFilterEventsInput(cwl)
	formatter = config.getFormatter()

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
