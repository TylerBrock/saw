package config

import (
	"errors"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Configuration struct {
	Group      string
	Fuzzy      bool
	Prefix     string
	Start      string
	End        string
	Filter     string
	Streams    []types.LogStream
	Descending bool
	OrderBy    string
}

// Define the order of time formats to attempt to use to parse our input absolute time
var absoluteTimeFormats = []string{
	time.RFC3339,

	"2006-01-02",          // Simple date
	"2006-01-02 15:04:05", // Simple date & time
}

// Parse the input string into a time.Time object.
// Provide the currentTime as a parameter to support relative time.
func getTime(timeStr string, currentTime time.Time) (time.Time, error) {
	relative, err := time.ParseDuration(timeStr)
	if err == nil {
		return currentTime.Add(relative), nil
	}

	// Iterate over available absolute time formats until we find one that works
	for _, timeFormat := range absoluteTimeFormats {
		absolute, err := time.Parse(timeFormat, timeStr)

		if err == nil {
			return absolute, err
		}
	}

	return time.Time{}, errors.New("Could not parse relative or absolute time")
}

func (c *Configuration) DescribeLogGroupsInput() *cloudwatchlogs.DescribeLogGroupsInput {
	input := cloudwatchlogs.DescribeLogGroupsInput{}
	if c.Prefix != "" {
		input.LogGroupNamePrefix = &c.Prefix
	}
	return &input
}

func (c *Configuration) DescribeLogStreamsInput() *cloudwatchlogs.DescribeLogStreamsInput {
	input := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &c.Group,
		Descending:   &c.Descending,
	}
	if c.OrderBy != "" {
		input.OrderBy = types.OrderBy(c.OrderBy)
	}
	if c.Prefix != "" {
		input.LogStreamNamePrefix = &c.Prefix
	}
	return &input
}

func (c *Configuration) FilterLogEventsInput() *cloudwatchlogs.FilterLogEventsInput {
	input := cloudwatchlogs.FilterLogEventsInput{
		Interleaved:  aws.Bool(true),
		LogGroupName: &c.Group,
	}

	if len(c.Streams) != 0 {
		input.LogStreamNames = c.TopStreamNames()
	}

	currentTime := time.Now()
	absoluteStartTime := currentTime
	if c.Start != "" {
		st, err := getTime(c.Start, currentTime)
		if err == nil {
			absoluteStartTime = st
		}
	}
	input.StartTime = aws.Int64(absoluteStartTime.UnixMilli())

	if c.End != "" {
		et, err := getTime(c.End, currentTime)
		if err == nil {
			input.EndTime = aws.Int64(et.UnixMilli())
		}
	}

	if c.Filter != "" {
		input.FilterPattern = &c.Filter
	}

	return &input
}

func (c *Configuration) TopStreamNames() []string {
	// FilerLogEvents can only take 100 streams so lets sort by LastEventTimestamp
	// (descending) and take only the names of the most recent 100.
	sort.Slice(c.Streams, func(i int, j int) bool {
		return *c.Streams[i].LastEventTimestamp > *c.Streams[j].LastEventTimestamp
	})

	numStreams := 100
	if len(c.Streams) < 100 {
		numStreams = len(c.Streams)
	}

	streamNames := make([]string, numStreams)
	for i := 0; i < numStreams; i++ {
		streamNames[i] = *c.Streams[i].LogStreamName
	}

	return streamNames
}
