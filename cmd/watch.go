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

var watchOutputConfig config.OutputConfiguration

var watchCommand = &cobra.Command{
	Use:   "watch <log group>",
	Short: "Continuously stream log events",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("watching streams requires log group argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		watchConfig.Group = args[0]
		b, err := blade.NewBlade(cmd.Context(), &watchConfig, &awsConfig, &watchOutputConfig)
		if err != nil {
			return
		}
		if watchConfig.Prefix != "" {
			streams, err := b.GetLogStreams(cmd.Context())
			if err != nil {
				return fmt.Errorf("failed to get log streams: %w", err)
			}
			if len(streams) == 0 {
				fmt.Printf("No streams found in %s with prefix %s\n", watchConfig.Group, watchConfig.Prefix)
				fmt.Printf("To view available streams: `saw streams %s`\n", watchConfig.Group)
				os.Exit(3)
			}
			watchConfig.Streams = streams
		}
		return b.StreamEvents(cmd.Context())
	},
}

func init() {
	watchCommand.Flags().BoolVar(&watchConfig.Fuzzy, "fuzzy", false, "log group fuzzy match")
	watchCommand.Flags().StringVar(&watchConfig.Prefix, "prefix", "", "log stream prefix filter")
	watchCommand.Flags().StringVar(&watchConfig.Filter, "filter", "", "event filter pattern")
	watchCommand.Flags().BoolVar(&watchOutputConfig.Raw, "raw", false, "print raw log event without timestamp or stream prefix")
	watchCommand.Flags().BoolVar(&watchOutputConfig.Expand, "expand", false, "indent JSON log messages")
	watchCommand.Flags().BoolVar(&watchOutputConfig.Invert, "invert", false, "invert colors for light terminal themes")
	watchCommand.Flags().BoolVar(&watchOutputConfig.RawString, "rawString", false, "print JSON strings without escaping")
}
