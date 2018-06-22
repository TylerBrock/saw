package cmd

import (
	"fmt"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

// TODO: colorize based on logGroup prefix (/aws/lambda, /aws/kinesisfirehose, etc...)
var groupsConfig config.Configuration

var GroupsCommand = &cobra.Command{
	Use:   "groups",
	Short: "List log groups",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		b := blade.NewBlade(&groupsConfig, &awsConfig, nil)
		logGroups := b.GetLogGroups()
		for _, group := range logGroups {
			fmt.Println(*group.LogGroupName)
		}
	},
}

func init() {
	GroupsCommand.Flags().StringVar(&groupsConfig.Prefix, "prefix", "", "log group prefix filter")
}
