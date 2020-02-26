package cmd

import (
	"github.com/TylerBrock/saw/config"
	zsh "github.com/rsteube/cobra-zsh-gen"
	"github.com/spf13/cobra"
)

// SawCommand is the main top-level command
var SawCommand = &cobra.Command{
	Use:   "saw <command>",
	Short: "A fast, multipurpose tool for AWS CloudWatch Logs",
	Long:  "Saw is a fast, multipurpose tool for AWS CloudWatch Logs.",
	Example: `  saw groups
  saw streams production
  saw watch production`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var awsConfig config.AWSConfiguration

func init() {
	SawCommand.PersistentFlags().StringVar(&awsConfig.Region, "region", "", "override profile AWS region")
	SawCommand.PersistentFlags().StringVar(&awsConfig.Profile, "profile", "", "override default AWS profile")

	zsh.Gen(SawCommand).FlagCompletion(zsh.ActionMap{
		"region": zsh.ActionValues(awsRegions()...),
		"profile": zsh.ActionCallback(func(args []string) zsh.Action {
			if profiles, err := awsProfiles(); err != nil {
				return zsh.ActionMessage(err.Error())
			} else {
				return zsh.ActionValues(profiles...)
			}
		}),
	})
}
