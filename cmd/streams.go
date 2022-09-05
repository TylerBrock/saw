package cmd

import (
	"errors"
	"fmt"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var streamsConfig config.Configuration

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
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		streamsConfig.Group = args[0]
		b := blade.NewBlade(&streamsConfig, &awsConfig, nil)

		logStreams, err := b.GetLogStreams()
		if err != nil {
			return err
		}
		for _, stream := range logStreams {
			fmt.Println(*stream.LogStreamName)
		}
		return
	},
}

func init() {
	streamsCommand.Flags().BoolVar(&streamsConfig.Fuzzy, "fuzzy", false, "log group fuzzy match")
	streamsCommand.Flags().StringVar(&streamsConfig.Prefix, "prefix", "", "stream prefix filter")
	streamsCommand.Flags().StringVar(&streamsConfig.OrderBy, "orderBy", "LogStreamName", "order streams by LogStreamName or LastEventTime")
	streamsCommand.Flags().BoolVar(&streamsConfig.Descending, "descending", false, "order streams descending")
}
