package cmd

import (
	"github.com/koushamad/btx/src/app/collector"
	"github.com/spf13/cobra"
)

// miningCmd represents the mining command
var collectorCmd = &cobra.Command{
	Use:   "collector",
	Short: "Collector the BTC addresses from the blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		return collector.Collect()
	},
}

func init() {
	rootCmd.AddCommand(collectorCmd)
}
