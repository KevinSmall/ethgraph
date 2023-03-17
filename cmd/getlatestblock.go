package cmd

import (
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/services"
	"github.com/spf13/cobra"
)

// getlatestblockCmd represents the getlatestblock command for getting the most recent block
var getlatestblockCmd = &cobra.Command{
	Use:   "getblock <url>",
	Short: "Gets latest block number",
	Long: `Gets the most recent block number from the given chain:

    ethgraph getblock "https://chain-rpc-endpoint"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logr.SetVerbosity(false)
		services.GetLatestBlockNumber(args[0])
	},
	Aliases: []string{"glb"},
}

func init() {

	rootCmd.AddCommand(getlatestblockCmd)
}
