package cmd

import (
	"fmt"
	"github.com/TylerBrock/saw/config"
	zsh "github.com/rsteube/cobra-zsh-gen"
	"github.com/spf13/cobra"
)

// TODO: colorize based on logGroup prefix (/aws/lambda, /aws/kinesisfirehose, etc...)
var groupsConfig config.Configuration

var groupsCommand = &cobra.Command{
	Use:   "groups",
	Short: "List log groups",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		for _, name := range groupNames() {
			fmt.Println(name)
		}
	},
}

func init() {
	groupsCommand.Flags().StringVar(&groupsConfig.Prefix, "prefix", "", "log group prefix filter")
	SawCommand.AddCommand(groupsCommand)

	zsh.Gen(groupsCommand).FlagCompletion(zsh.ActionMap{
		"prefix": zsh.ActionCallback(func(args []string) zsh.Action {
			return zsh.ActionMultiParts('/', groupNames()...)
		}),
	})
}
