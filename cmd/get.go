package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var getConfigGlobal config.Configuration
var getOutputConfig config.OutputConfiguration

var getCommand = &cobra.Command{
	Use:   "get <log group>",
	Short: "Get log events",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("getting events requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := runMultiGroup(args[0], func(group string) {
			getConfig := getConfigGlobal
			getConfig.Group = group
			b := blade.NewBlade(&getConfig, &awsConfig, &getOutputConfig)
			if getConfig.Prefix != "" {
				streams := b.GetLogStreams()
				if len(streams) == 0 {
					fmt.Printf("No streams found in %s with prefix %s\n", getConfig.Group, getConfig.Prefix)
					fmt.Printf("To view available streams: `saw streams %s`\n", getConfig.Group)
					os.Exit(3)
				}
				getConfig.Streams = streams
			}
			b.GetEvents()
		})

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	getCommand.Flags().StringVar(&getConfigGlobal.Prefix, "prefix", "", "log group prefix filter")
	getCommand.Flags().StringVar(
		&getConfigGlobal.Start,
		"start",
		"",
		`start getting the logs from this point
Takes an absolute timestamp in RFC3339 format, or a relative time (eg. -2h).
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	getCommand.Flags().StringVar(
		&getConfigGlobal.End,
		"stop",
		"now",
		`stop getting the logs at this point
Takes an absolute timestamp in RFC3339 format, or a relative time (eg. -2h).
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	getCommand.Flags().StringVar(&getConfigGlobal.Filter, "filter", "", "event filter pattern")
	getCommand.Flags().BoolVar(&getOutputConfig.Pretty, "pretty", false, "print timestamp and stream name prefix")
	getCommand.Flags().BoolVar(&getOutputConfig.Expand, "expand", false, "indent JSON log messages")
	getCommand.Flags().BoolVar(&getOutputConfig.Invert, "invert", false, "invert colors for light terminal themes")
	getCommand.Flags().BoolVar(&getOutputConfig.RawString, "rawString", false, "print JSON strings without escaping")
}
