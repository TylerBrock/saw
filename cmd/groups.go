package cmd

import (
	"fmt"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

// TODO: colorize based on logGroup prefix (/aws/lambda, /aws/kinesisfirehose, etc...)
var groupsConfig config.Configuration

var groupsCommand = &cobra.Command{
	Use:   "groups",
	Short: "List log groups",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		b, err := blade.NewBlade(cmd.Context(), &groupsConfig, &awsConfig, nil)
		if err != nil {
			return
		}
		logGroups, err := b.GetLogGroups(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to get log groups: %w", err)
		}
		for _, group := range logGroups {
			fmt.Println(*group.LogGroupName)
		}
		return
	},
}

func init() {
	groupsCommand.Flags().StringVar(&groupsConfig.Prefix, "prefix", "", "log group prefix filter")
}
