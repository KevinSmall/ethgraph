package services

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/masterdata/tokens"
	"strconv"
	"time"
)

const destinationFilenameTemplate string = "masterdata/tokens/data/tokens_%s_working.csv"

func BuildMasterData(url string, blockFrom uint64, blockTo uint64) {
	start := time.Now()

	evmChain, err := chain.CreateEvmClient(url)
	if err != nil {
		logr.Error.Panicln(err)
	}
	logr.Info.Printf("Connecting to: %s with ChainId: %s", evmChain.Name, evmChain.ChainId)

	// Prepare master data lookups
	tokens.Init(evmChain.ChainId)

	//-------------------------------------------------------------------------
	// GET EVENTS
	//-------------------------------------------------------------------------
	var allEvents []*chain.TransferEvent
	blockCount := int(blockTo-blockFrom) + 1
	for i := 0; i < blockCount; i++ {
		events, err := chain.GetTransferEventsByBlock(evmChain.Client, blockFrom, "")
		if err != nil {
			logr.Error.Panicln(err)
		}
		blockFrom++
		allEvents = append(allEvents, events...)
	}
	logr.Info.Printf("Total events: %v", len(allEvents))

	//-------------------------------------------------------------------------
	// BUILD UNIQUE ADDRESSES
	//-------------------------------------------------------------------------
	uniqueAddressesMap := tokens.BuildUniqueTokenAddressesFromEvents(allEvents)
	// Optional, maybe no issue with volumes
	// uniqueAddressesMap = filterAddressesRemoveNoise(uniqueAddressesMap, 25)

	// Sort addresses by usage count descending, makes final file more intelligible
	uniqueTokenAddresses := tokens.SortAddressesByCount(uniqueAddressesMap)

	//-------------------------------------------------------------------------
	// PROCESS UNIQUE ADDRESS, BUILDING MASTER DATA
	//-------------------------------------------------------------------------
	// Create a slice to store the token information
	var tokenInfo [][]string

	// Iterate through the unique addresses and get the token information for each address
	logr.Info.Printf("Unique token addresses: %d\n", len(uniqueTokenAddresses))
	chunkSize := (len(uniqueTokenAddresses) / 100) + 1
	i := 0
	for _, a := range uniqueTokenAddresses {
		if i%chunkSize == 0 {
			// Calculate the current progress as a percentage
			progress := ((i * 100) / len(uniqueTokenAddresses)) + 1

			// PrintSummary out the progress as a percentage
			logr.Info.Printf("Progress: %d%%\n", progress)
		}

		tokenData := tokens.GetTokenFromChain(evmChain.ChainId, evmChain.Client, a.Address, a.TransferType)
		tokenInfo = append(tokenInfo, []string{
			evmChain.ChainId,
			tokenData.Name,
			tokenData.Symbol,
			strconv.Itoa(tokenData.Decimals),
			tokenData.TokenAddress})
		i++
	}
	// tokenInfo is built

	//-------------------------------------------------------------------------
	// WRITE THE FILE
	//-------------------------------------------------------------------------
	filename := fmt.Sprintf(destinationFilenameTemplate, evmChain.ChainId)
	tokens.WriteTokenInfoToFile(filename, tokenInfo)

	logr.Info.Println("Tokens written: ", len(tokenInfo), " to file: ", filename)
	elapsed := time.Since(start)
	logr.Info.Printf("Runtime: %s", elapsed)
}
