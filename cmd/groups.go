package cmd

import (
	"errors"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if groupsConfig.Fuzzy {
			if len(args) < 1 {
				return errors.New("listing groups with fuzzy search requires log group argument")
			}
			groupsConfig.Group = args[0]
		}
		b := blade.NewBlade(&groupsConfig, &awsConfig, nil)
		logGroups := b.GetLogGroups()
		for _, group := range logGroups {
			fmt.Println(*group.LogGroupName)
		}
		return nil
	},
}

func init() {
	groupsCommand.Flags().BoolVar(&groupsConfig.Fuzzy, "fuzzy", false, "log group fuzzy match")
	groupsCommand.Flags().StringVar(&groupsConfig.Prefix, "prefix", "", "log group prefix filter")
}
