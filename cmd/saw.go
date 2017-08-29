package cmd

import "github.com/spf13/cobra"

var SawCommand = &cobra.Command{
	Use:   "saw",
	Short: "A fast, multipurpose tool for AWS CloudWatch Logs",
	Long:  "Saw is a fast, multipurpose tool for cutting through AWS CloudWatch Logs.",
	Example: `  saw version
  saw groups --prefix "/aws"
  saw streams production --prefix "api"`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func init() {
	SawCommand.AddCommand(GroupsCommand)
	SawCommand.AddCommand(StreamsCommand)
	SawCommand.AddCommand(VersionCommand)
	SawCommand.AddCommand(WatchCommand)
	//Saw.AddCommand(Get)
	//Saw.AddCommand(Delete)
}
