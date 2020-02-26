package cmd

import (
	"strings"

	"github.com/TylerBrock/saw/blade"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"gopkg.in/ini.v1"
)

func awsRegions() []string {
	regions := endpoints.AwsPartition().Regions()
	keys := make([]string, 0, len(regions))
	for k := range regions {
		keys = append(keys, k)
	}
	return keys
}

func awsProfiles() (profiles []string, err error) {
	profiles = []string{"default"}

	var cfg *ini.File
	if cfg, err = ini.Load(defaults.SharedConfigFilename()); err == nil {
		for _, section := range cfg.SectionStrings() {
			if strings.HasPrefix(section, "profile ") {
				profiles = append(profiles, strings.TrimPrefix(section, "profile "))
			}
		}
	}
	return
}

func groupNames() []string {
	b := blade.NewBlade(&groupsConfig, &awsConfig, nil)
	logGroups := b.GetLogGroups()
	names := make([]string, len(logGroups))
	for index, group := range logGroups {
		names[index] = *group.LogGroupName
	}
	return names
}

func streams(group string) []string {
	streamsConfig.Group = group
	streamsConfig.OrderBy = "LastEventTime"
	streamsConfig.Descending = true
	b := blade.NewBlade(&streamsConfig, &awsConfig, nil)

	logStreams := b.GetLogStreams(1)

	streams := make([]string, len(logStreams))
	for index, stream := range logStreams {
		streams[index] = *stream.LogStreamName
	}
	return streams
}

func streamPrefixes(group string) []string {
	prefixes := make(map[string]bool)
	for _, s := range streams(group) {
		parts := strings.Split(s, "/")
		withoutId := parts[:len(parts)-1]
		prefixes[strings.Join(withoutId, "/")] = true
	}

	values := make([]string, len(prefixes))
	i := 0
	for key, _ := range prefixes {
		values[i] = key
		i++
	}
	return values
}
