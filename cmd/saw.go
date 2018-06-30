package cmd

import (
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var sawCommand = &cobra.Command{
	Use:   "saw <command>",
	Short: "A fast, multipurpose tool for AWS CloudWatch Logs",
	Long:  "Saw is a fast, multipurpose tool for cutting through AWS CloudWatch Logs.",
	Example: `  saw version
  saw groups --prefix prod
  saw streams production --prefix api
  saw watch production --prefix api --filter error`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var awsConfig config.AWSConfiguration

func init() {
	sawCommand.AddCommand(groupsCommand)
	sawCommand.AddCommand(streamsCommand)
	sawCommand.AddCommand(versionCommand)
	sawCommand.AddCommand(watchCommand)
	sawCommand.AddCommand(getCommand)
	sawCommand.PersistentFlags().StringVar(&awsConfig.Region, "region", "", "override profile AWS region")
	sawCommand.PersistentFlags().StringVar(&awsConfig.Profile, "profile", "", "override default AWS profile")
}
