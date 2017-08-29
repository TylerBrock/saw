package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.2"

var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the version string",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", version)
	},
}
