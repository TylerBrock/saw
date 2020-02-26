package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	zsh "github.com/rsteube/cobra-zsh-gen"
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
	Run: func(cmd *cobra.Command, args []string) {
		watchConfig.Group = args[0]
		b := blade.NewBlade(&watchConfig, &awsConfig, &watchOutputConfig)
		if watchConfig.Prefix != "" {
			streams := b.GetLogStreams(0)
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
	watchCommand.Flags().StringVar(&watchConfig.Prefix, "prefix", "", "log stream prefix filter")
	watchCommand.Flags().StringVar(&watchConfig.Filter, "filter", "", "event filter pattern")
	watchCommand.Flags().BoolVar(&watchOutputConfig.Raw, "raw", false, "print raw log event without timestamp or stream prefix")
	watchCommand.Flags().BoolVar(&watchOutputConfig.Expand, "expand", false, "indent JSON log messages")
	watchCommand.Flags().BoolVar(&watchOutputConfig.Invert, "invert", false, "invert colors for light terminal themes")
	watchCommand.Flags().BoolVar(&watchOutputConfig.RawString, "rawString", false, "print JSON strings without escaping")
	SawCommand.AddCommand(watchCommand)

	zsh.Gen(watchCommand).FlagCompletion(zsh.ActionMap{
		"prefix": zsh.ActionCallback(func(args []string) zsh.Action {
			if len(args) == 0 {
				return zsh.ActionMessage("missing log group argument")
			}
			return zsh.ActionMultiParts('/', streamPrefixes(args[0])...)
		}),
	})

	zsh.Gen(watchCommand).PositionalCompletion(
		zsh.ActionCallback(func(args []string) zsh.Action {
			return zsh.ActionValues(groupNames()...)
		}),
	)
}
