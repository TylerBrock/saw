package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	prettyjson "github.com/TylerBrock/go-prettyjson"
	color "github.com/fatih/color"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func printLogs() {

}

func main() {
	// TODO: use StringVar to store directly into DescribeLogStreamsInput struct
	logGroupName := flag.String("logGroup", "", "the log group to stream")
	logStreamPrefix := flag.String("prefix", "", "the log stream prefix")
	flag.Parse()

	config := aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}
	sess := session.Must(session.NewSession(&config))
	cw := cloudwatchlogs.New(sess)
	lsd := cloudwatchlogs.DescribeLogStreamsInput{}
	lsd.SetLogGroupName(*logGroupName)
	if *logStreamPrefix != "" {
		lsd.SetLogStreamNamePrefix(*logStreamPrefix)
	} else {
		lsd.SetOrderBy("LastEventTime")
		lsd.SetDescending(true)
	}
	lsd.SetLimit(10)
	logStreamsDescriptions, err := cw.DescribeLogStreams(&lsd)
	if err != nil {
		fmt.Println("error", err)
	}
	for _, stream := range logStreamsDescriptions.LogStreams {
		name := aws.StringValue(stream.LogStreamName)
		firstEvent := aws.MillisecondsTimeValue(stream.FirstEventTimestamp)
		lastEvent := aws.MillisecondsTimeValue(stream.LastEventTimestamp)
		storedBytes := aws.Int64Value(stream.StoredBytes)
		fmt.Println(name, firstEvent, lastEvent, storedBytes)
	}
	stream := logStreamsDescriptions.LogStreams[0]
	gle := cloudwatchlogs.GetLogEventsInput{}
	gle.SetLogGroupName("production")
	gle.SetLogStreamName(aws.StringValue(stream.LogStreamName))
	gle.SetStartTime(aws.Int64Value(stream.FirstEventTimestamp))
	gle.SetEndTime(aws.Int64Value(stream.LastEventTimestamp))
	logEvents, err := cw.GetLogEvents(&gle)
	if err != nil {
		fmt.Println("error", err)
	}
	formatter := prettyjson.NewFormatter()
	var dat map[string]interface{}
	red := color.New(color.FgRed).SprintFunc()
	for _, event := range logEvents.Events {
		str := aws.StringValue(event.Message)
		byt := []byte(str)
		date := aws.MillisecondsTimeValue(event.IngestionTime)
		dateStr := date.Format(time.RFC3339)
		if err := json.Unmarshal(byt, &dat); err != nil {
			fmt.Printf("[%s] %s\n", red(dateStr), str)
		} else {
			s, _ := formatter.Marshal(dat)
			fmt.Printf("[%s] %s\n", red(dateStr), string(s))
		}
	}
}
