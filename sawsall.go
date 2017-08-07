package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/fatih/color"
)

type jsonLog map[string]interface{}

type formatter struct {
	KeyColor        *color.Color
	StringColor     *color.Color
	BoolColor       *color.Color
	NumberColor     *color.Color
	NullColor       *color.Color
	StringMaxLength int
	DisabledColor   bool
	Indent          int
}

var fmter = formatter{
	KeyColor:        color.New(color.FgWhite),
	StringColor:     color.New(color.FgGreen),
	BoolColor:       color.New(color.FgYellow),
	NumberColor:     color.New(color.FgCyan),
	NullColor:       color.New(color.FgMagenta),
	StringMaxLength: 0,
	DisabledColor:   false,
	Indent:          0,
}

func (jl *jsonLog) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	var next string
	for _, v := range *jl {
		switch v.(type) {
		case string:
			next = fmter.StringColor.SprintFunc()(v)
		case map[string]interface{}:
			next = "{}"
		case float64:
			next = "1"
		case bool:
			next = "true"
		case nil:
			next = "null"
		case []interface{}:
			next = "[]"
		default:
			next = "XXX"
		}
		fmt.Println("next:", next, "val", v)
		buffer.WriteString(next)
	}
	return buffer.Bytes(), nil
}

func handlePage(page *cloudwatchlogs.FilterLogEventsOutput, lastPage bool) bool {
	//formatter := prettyjson.NewFormatter()
	red := color.New(color.FgRed).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()

	for _, event := range page.Events {
		str := aws.StringValue(event.Message)
		bytes := []byte(str)
		date := aws.MillisecondsTimeValue(event.IngestionTime)
		dateStr := date.Format(time.RFC3339)
		streamStr := aws.StringValue(event.LogStreamName)
		jl := jsonLog{}
		if err := json.Unmarshal(bytes, &jl); err != nil {
			fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), str)
		} else {
			prettyStr, _ := json.Marshal(&jl)
			fmt.Printf("[%s] (%s) %s\n", red(dateStr), white(streamStr), prettyStr)
		}
	}
	return !lastPage
}

func main() {
	logGroupName := flag.String("logGroup", "", "the log group to stream")
	logStreams := flag.String("logStreams", "", "the log stream prefix")
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

	if len(*logStreams) > 0 {
		streams := strings.Split(*logStreams, ",")
		streamNamePointers := make([]*string, len(streams))
		for i, stream := range streams {
			streamNamePointers[i] = &stream
		}
		input.SetLogStreamNames(streamNamePointers)
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
