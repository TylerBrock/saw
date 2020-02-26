package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.2.2"

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the version string",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", version)
	},
}

func init() {
	SawCommand.AddCommand(versionCommand)
}
