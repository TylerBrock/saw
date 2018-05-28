package cmd

import (
	"errors"

	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var getConfig config.Configuration

var GetCommand = &cobra.Command{
	Use:   "get",
	Short: "Get log events",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("listing streams requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		//b := blade.NewBlade(&getConfig)
	},
}

func init() {
	GetCommand.Flags().StringVar(&groupsConfig.Prefix, "prefix", "", "log group prefix filter")
	GetCommand.Flags().StringVar(&groupsConfig.Start, "start", "", "start getting the logs from this point")
	GetCommand.Flags().StringVar(&groupsConfig.End, "end", "now", "stop getting the logs at this point")
}
