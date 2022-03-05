package cmd

import (
	"fmt"
	"os"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var queryConfig config.Configuration
var queryOutputConfig config.OutputConfiguration

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query CloudWatch Logs Insights",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		b := blade.NewBlade(&queryConfig, &awsConfig, &queryOutputConfig)
		if queryConfig.Query == "" {
			fmt.Println("--query must be set")
			os.Exit(2)
		}
		if len(queryConfig.Groups) == 0 {
			fmt.Println("--groups must be defined at least once")
		}

		b.RunQuery()
	},
}

func init() {
	queryCommand.Flags().StringSliceVar(&queryConfig.Groups, "groups", []string{}, "CloudWatch Log Groups to query against")
	queryCommand.Flags().StringVar(&queryConfig.Query, "query", "", "CloudWatch Logs Insights query to run")
	queryCommand.Flags().StringVar(
		&queryConfig.Start,
		"start",
		"",
		`start getting the logs from this point
Takes an absolute timestamp in RFC3339 format, or a relative time (eg. -2h).
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	queryCommand.Flags().StringVar(
		&queryConfig.End,
		"stop",
		"now",
		`stop getting the logs at this point
Takes an absolute timestamp in RFC3339 format, or a relative time (eg. -2h).
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
	)
	queryCommand.Flags().BoolVar(&queryOutputConfig.NoHeaders, "no-headers", false, "Whether or not to print field headers")
}
