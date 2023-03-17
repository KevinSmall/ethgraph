package cmd

import (
	"errors"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/services"
	"github.com/spf13/cobra"
)

// buildmdCmd represents the buildmd command for building .csv token master data
var buildmdCmd = &cobra.Command{
	Use:   "buildmd <url>",
	Short: "For internal use to build cache of token names and decimals.",
	Long: `For internal use. Creates a .csv file of all tokens in the given block range
that have valid ERC20 master data of name and decimals. This .csv is built into
new releases of ethgraph. For example:

    ethgraph buildmd "https://chain-rpc-endpoint" -f 16670050 -t 16670150 

Will scan the block range from and to and all valid token master data 
will get written to a .csv file.`,
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
		// validation successful
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		logr.SetVerbosity(true)
		services.BuildMasterData(args[0], flagBlockFrom, flagBlockTo)
	},
}

func init() {

	rootCmd.AddCommand(buildmdCmd)

	buildmdCmd.PersistentFlags().Uint64VarP(&flagBlockFrom, "block-from", "f", 0, "Block number from eg 16667050")
	buildmdCmd.MarkPersistentFlagRequired("block-from")

	buildmdCmd.PersistentFlags().Uint64VarP(&flagBlockTo, "block-to", "t", 1, "Block number to eg 16667150")
	buildmdCmd.MarkPersistentFlagRequired("block-to")
}
