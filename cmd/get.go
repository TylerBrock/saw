package cmd

import (
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var getConfig config.Configuration

var Get = &cobra.Command{
	Use:   "get",
	Short: "Get log events",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		//b := blade.NewBlade(&getConfig)
	},
}
