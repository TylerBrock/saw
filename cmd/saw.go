package cmd

import (
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var SawCommand = &cobra.Command{
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
	SawCommand.AddCommand(GroupsCommand)
	SawCommand.AddCommand(StreamsCommand)
	SawCommand.AddCommand(VersionCommand)
	SawCommand.AddCommand(WatchCommand)
	SawCommand.AddCommand(GetCommand)
	//Saw.AddCommand(Delete)
	SawCommand.PersistentFlags().StringVar(&awsConfig.Region, "region", "", "override profile AWS region")
	SawCommand.PersistentFlags().StringVar(&awsConfig.Profile, "profile", "", "override default AWS profile")
}
