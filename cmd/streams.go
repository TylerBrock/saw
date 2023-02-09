package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var streamsConfigGlobal config.Configuration

var streamsCommand = &cobra.Command{
	Use:   "streams <log group>",
	Short: "List streams in log group",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("listing streams requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runMultiGroup(args[0], func(group string) {
			streamsConfig := streamsConfigGlobal
			streamsConfig.Group = group
			b := blade.NewBlade(&streamsConfig, &awsConfig, nil)

			logStreams := b.GetLogStreams()
			for _, stream := range logStreams {
				fmt.Println(*stream.LogStreamName)
			}
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	streamsCommand.Flags().StringVar(&streamsConfigGlobal.Prefix, "prefix", "", "stream prefix filter")
	streamsCommand.Flags().StringVar(&streamsConfigGlobal.OrderBy, "orderBy", "LogStreamName", "order streams by LogStreamName or LastEventTime")
	streamsCommand.Flags().BoolVar(&streamsConfigGlobal.Descending, "descending", false, "order streams descending")
}
