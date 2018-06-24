package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var watchConfig config.Configuration
var outputConfig config.OutputConfiguration

var WatchCommand = &cobra.Command{
	Use:   "watch <log group>",
	Short: "Continously stream log events",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("watching streams requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		watchConfig.Group = args[0]
		b := blade.NewBlade(&watchConfig, &awsConfig, &outputConfig)
		if watchConfig.Prefix != "" {
			streams := b.GetLogStreams()
			if len(streams) == 0 {
				fmt.Printf("No streams found in %s with prefix %s\n", watchConfig.Group, watchConfig.Prefix)
				fmt.Printf("To view available streams: `saw streams %s`\n", watchConfig.Group)
				os.Exit(3)
			}
			watchConfig.Streams = streams
		}
		b.StreamEvents()
	},
}

func init() {
	WatchCommand.Flags().StringVar(&watchConfig.Prefix, "prefix", "", "log stream prefix filter")
	WatchCommand.Flags().StringVar(&watchConfig.Filter, "filter", "", "event filter pattern")
	WatchCommand.Flags().BoolVar(&outputConfig.Expand, "expand", false, "indent JSON log messages")
	WatchCommand.Flags().BoolVar(&outputConfig.Invert, "invert", false, "invert colors for light terminal themes")
	WatchCommand.Flags().BoolVar(&outputConfig.RawString, "rawString", false, "print JSON strings without escaping")
	WatchCommand.Flags().BoolVar(&outputConfig.HideDate, "hideDate", false, "omit the date from log output")
	WatchCommand.Flags().BoolVar(&outputConfig.HideStreamName, "hideName", false, "omit the stream name from log output")
}
