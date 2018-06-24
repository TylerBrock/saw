package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:                   "completion SHELL",
	Short:                 "Generate shell completion scripts",
	Long:                  "Generate shell completion scripts",
	Hidden:                true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify argument of either 'bash' or 'zsh'")
			os.Exit(4)
		}

		if args[0] == "bash" {
			SawCommand.GenBashCompletion(os.Stdout)
		} else if args[0] == "zsh" {
			SawCommand.GenZshCompletion(os.Stdout)
		} else {
			fmt.Println("Not able to generate completions for", args[0])
			os.Exit(5)
		}
	},
}
