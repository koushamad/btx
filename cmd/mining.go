package cmd

import (
	"github.com/koushamad/btx/src/app/mining"

	"github.com/spf13/cobra"
)

// miningCmd represents the mining command
var miningCmd = &cobra.Command{
	Use:   "mining",
	Short: "Create Bitcoin addresses and save them to a file",
	Run: func(cmd *cobra.Command, args []string) {
		chunkSize := 1000000
		mining.Mine(chunkSize)
	},
}

func init() {
	// add address chunk size flag
	miningCmd.Flags().IntP("chunk-size", "c", 1000000, "Chunk size for generating addresses")

	rootCmd.AddCommand(miningCmd)
}
