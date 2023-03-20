package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

//------------------------------------------------------------------------------
//  Global Flags, these can get written to by any command
//------------------------------------------------------------------------------

var flagBlockFrom uint64
var flagBlockTo uint64
var flagOnlyThisTokenAddress *string
var flagDoNotFetchMissingMasterData *bool
var flagForceSerialExecution *bool
var flagClearTokenCache *bool
var flagIsVerboseOutputRequested *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "1.0.1",
	Use:     "ethgraph <command> <url> [flags]",
	Short:   "Ethgraph builds GraphML files of ERC20, ERC721 and ERC1155 token movements from EVM-compatible chains",
	Long: `Ethgraph is a CLI tool that queries Ethereum or other EVM-compatible chains to
extract the logs that contracts emit when a Transfer event occurs. Each Transfer event
shows the from and to address of a token movement. Ethgraph then builds a GraphML file with 
addresses as nodes and token movements as either edges (examples 1 and 2) or nodes (example 3).
Examples:

    1) select Transfer events by block range:
       ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150

    2) select Transfer events only for token USDT by block range:
       ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 -o 0xdAC17F958D2ee523a2206206994597C13D831ec7
	 
    3) select Transfer events by block range, clear the local token cache file first (deletes file .tokens_*_cache.csv):
       ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 -c

    4) select Transfer events by block range, fetching all data in serial, capping the number of HTTP requests to 10 per second:
	   ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 -s

    5) display the latest block number for a chain:
       ethgraph getblock "https://chain-rpc-endpoint"`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	//will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goethgraph01.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
