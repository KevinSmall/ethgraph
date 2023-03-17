package cmd

import (
	"errors"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/services"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// byblockCmd represents the byblock command to build by a block range
var byblockCmd = &cobra.Command{
	Use:   "byblock <url>",
	Short: "Builds GraphML from data selected by a range of blocks",
	Long: `Builds GraphML from data selected by a range of blocks. For example:

    1) select Transfer events by block range:
       ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 

    2) select Transfer events by block range, only for token USDT:
       ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 -o 0xdAC17F958D2ee523a2206206994597C13D831ec7

    3) select Transfer events by block range, clear the local token cache file first (deletes file .tokens_*_cache.csv):
       ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 -c

	4) select Transfer events by block range, fetching all data in serial, capping the number of HTTP requests to 10 per second:
	   ethgraph byblock "https://chain-rpc-endpoint" -f 16670050 -t 16670150 -s`,

	Args: cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate block from and to
		from, err := cmd.Flags().GetUint64("block-from")
		if err != nil {
			return err
		}
		to, err := cmd.Flags().GetUint64("block-to")
		if err != nil {
			return err
		}
		if from > to {
			return errors.New("the --block-from flag must be less than or equal to the --block-to flag")
		}
		// Validate "only this address" if it exists
		onlyThisAddress, err := cmd.Flags().GetString("only-token-address")
		if err != nil {
			return err
		}
		if onlyThisAddress != "" && !common.IsHexAddress(onlyThisAddress) {
			return errors.New("the --only-token-address value is not a valid hex address. Use for example 0xdAC17F958D2ee523a2206206994597C13D831ec7 for USDT")
		}
		// validation successful
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if *flagIsVerboseOutputRequested {
			logr.SetVerbosity(true)
		} else {
			logr.SetVerbosity(false)
		}
		services.BuildByBlockRange(args[0],
			flagBlockFrom, flagBlockTo,
			*flagOnlyThisTokenAddress,
			*flagDoNotFetchMissingMasterData,
			*flagForceSerialExecution,
			*flagClearTokenCache)
	},
	Aliases: []string{"byb"},
}

func init() {

	rootCmd.AddCommand(byblockCmd)

	byblockCmd.PersistentFlags().Uint64VarP(&flagBlockFrom, "block-from", "f", 0, "Block number from eg 16667050")
	byblockCmd.MarkPersistentFlagRequired("block-from")

	byblockCmd.PersistentFlags().Uint64VarP(&flagBlockTo, "block-to", "t", 1, "Block number to eg 16667150")
	byblockCmd.MarkPersistentFlagRequired("block-to")

	flagOnlyThisTokenAddress = byblockCmd.PersistentFlags().StringP("only-token-address", "o", "", "Only select events for the specified token address eg USDT is 0xdAC17F958D2ee523a2206206994597C13D831ec7")

	flagDoNotFetchMissingMasterData = byblockCmd.PersistentFlags().BoolP("no-fetch-master-data", "n", false, "If set with -n then no fetch of master data for unknown tokens (faster runtime). If omitted (which is the default) then master data is fetched (longer runtime).")

	flagForceSerialExecution = byblockCmd.PersistentFlags().BoolP("force-serial-execution", "s", false, "If set with -s then serial execution is forced, PLUS a cap is set on HTTP requests to 10 per second (longer runtime).")

	flagClearTokenCache = byblockCmd.PersistentFlags().BoolP("clear-token-cache", "c", false, "If set with -c then the token cache file .tokens_*_cache.csv is deleted (longer runtime). The * in the filename is the chainId see https://chainlist.org/, so 1 for Ethereum.")

	flagIsVerboseOutputRequested = byblockCmd.PersistentFlags().BoolP("verbose-output", "v", false, "If set with -v then detailed logging information written to stdout.")
}
