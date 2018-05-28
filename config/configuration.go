package config

import (
	"errors"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

const maxLimit = 50

type Configuration struct {
	Group      string
	Prefix     string
	Start      string
	End        string
	Filter     string
	Region     string
	Streams    []*cloudwatchlogs.LogStream
	Descending bool
	OrderBy    string
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

func (c *Configuration) DescribeLogGroupsInput() *cloudwatchlogs.DescribeLogGroupsInput {
	input := cloudwatchlogs.DescribeLogGroupsInput{}
	if c.Prefix != "" {
		input.SetLogGroupNamePrefix(c.Prefix)
	}
	return &input
}

func (c *Configuration) DescribeLogStreamsInput() *cloudwatchlogs.DescribeLogStreamsInput {
	input := cloudwatchlogs.DescribeLogStreamsInput{}
	input.SetLogGroupName(c.Group)
	input.SetDescending(c.Descending)

	if c.OrderBy != "" {
		input.SetOrderBy(c.OrderBy)
	}

	if c.Prefix != "" {
		input.SetLogStreamNamePrefix(c.Prefix)
	}
	return &input
}

func (c *Configuration) FilterLogEventsInput() *cloudwatchlogs.FilterLogEventsInput {
	input := cloudwatchlogs.FilterLogEventsInput{}
	input.SetInterleaved(true)
	input.SetLogGroupName(c.Group)

	if len(c.Streams) != 0 {
		input.SetLogStreamNames(c.TopStreamNames())
	}

	absoluteStartTime := time.Now()
	if c.Start != "" {
		st, err := getTime(c.Start)
		if err == nil {
			absoluteStartTime = st
		}
	}
	input.SetStartTime(aws.TimeUnixMilli(absoluteStartTime))

	if c.End != "" {
		et, err := getTime(c.End)
		if err == nil {
			input.SetEndTime(aws.TimeUnixMilli(et))
		}
	}

	if c.Filter != "" {
		input.SetFilterPattern(c.Filter)
	}

	return &input
}

func (c *Configuration) TopStreamNames() []*string {
	// FilerLogEvents can only take 100 streams so lets sort by LastEventTimestamp
	// (descending) and take only the names of the most recent 100.
	sort.Slice(c.Streams, func(i int, j int) bool {
		return *c.Streams[i].LastEventTimestamp > *c.Streams[j].LastEventTimestamp
	})

	numStreams := 100

	if len(c.Streams) < 100 {
		numStreams = len(c.Streams)
	}

	streamNames := make([]*string, 0)

	for _, stream := range c.Streams[:numStreams] {
		streamNames = append(streamNames, stream.LogStreamName)
	}

	return streamNames
}
