package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.6"

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the version string",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", version)
	},
}
