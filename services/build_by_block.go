package services

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/blocks"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/graph"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/masterdata/addresses"
	"github.com/KevinSmall/ethgraph/masterdata/tokens"
	"github.com/ethereum/go-ethereum/common"
	"sort"
	"time"
)

const throttleHttpDelayMilliseconds = time.Millisecond * 100

// BuildByBlockRange is entry point for building graph based on block selection.
// Note injection of url string not a client, because later on in concurrent execution
// we want to create many clients
func BuildByBlockRange(
	url string,
	blockFrom uint64,
	blockTo uint64,
	onlyThisTokenAddress string,
	doNotFetchMissingMasterData bool,
	forceSerialExecution bool,
	clearTokenCache bool) {
	start := time.Now()

	// Client
	evmChain, err := chain.CreateEvmClient(url)
	if err != nil {
		logr.Error.Panicln(err)
	}
	logr.Info.Printf("Connecting to: %s with ChainId: %s\n", evmChain.Name, evmChain.ChainId)

	// Prepare []allEvents
	// Does do:      data cleansing, time field enrichment, ERC1155 decompose
	// Does not do:  business logic, no master data reads
	allEvents := getTransferEvents(evmChain, blockFrom, blockTo, forceSerialExecution, onlyThisTokenAddress)

	// Prepare token and address master data
	if clearTokenCache {
		tokens.DeleteTokenCache(evmChain.ChainId)
	}
	tokens.Init(evmChain.ChainId)
	addresses.Init(evmChain.ChainId)
	logr.Trace.Printf("Loaded %v popular addresses and names\n", addresses.GetAddressesLoadedCount())
	logr.Trace.Printf("Loaded %v token addresses and symbols\n", tokens.GetGlobalTokensLoadedCount())

	// Fetch missing token master data from chain, if desired
	if !doNotFetchMissingMasterData {
		fetchMissingTokenMasterData(evmChain, allEvents, forceSerialExecution)
	}

	// Prepare Graph, most business logic inc master data lookups is here
	ethGraph, creationResult := graph.CreateGraph(evmChain.Name, allEvents)
	creationResult.PrintSummary()

	// Write Graph
	filename := fmt.Sprintf("%s.graphml", evmChain.Name)
	err = graph.WriteGraph(filename, ethGraph)
	if err != nil {
		logr.Error.Panicln(err)
	}
	elapsed := time.Since(start)
	logr.Info.Printf("Runtime: %.3f seconds\n", elapsed.Seconds())
	logr.Info.Printf("File created: %s\n", filename)
}

func fetchMissingTokenMasterData(evmChain chain.EvmClient, allEvents []*chain.TransferEvent, forceSerialExecution bool) {

	// Augment the token master data (cache file and global map) based on allEvents
	uniqueAddressesMap := tokens.BuildUniqueTokenAddressesFromEvents(allEvents)

	// Reduce to a UNIQUE list of UNKNOWN tokens (those without master data)
	uniqueAddressesMap = tokens.RemoveUnknownTokens(uniqueAddressesMap)
	logr.Info.Println("Tokens not seen before:", len(uniqueAddressesMap))
	//tokens.PrintContents(uniqueAddressesMap)

	// Populate the missing master data
	var tokenMapToAdd map[string]tokens.TokenDataFromSource
	if forceSerialExecution {
		tokenMapToAdd = getTokenMasterDataForMissingTokensSerial(evmChain, uniqueAddressesMap)
	} else {
		tokenMapToAdd = getTokenMasterDataForMissingTokensConcurrent(evmChain, uniqueAddressesMap)
	}

	// Merge tokenMapToAdd entries into the tokenMap global data
	tokens.MergeTokensIntoGlobalTokenMap(tokenMapToAdd)

	// Merge tokenMapAdd into the local .csv cache file
	tokens.WriteGlobalTokenMapToCache(evmChain.ChainId)
}

func getTransferEvents(evmChain chain.EvmClient, blockFrom uint64, blockTo uint64, forceSerialExecution bool, onlyThisTokenAddress string) []*chain.TransferEvent {
	var allEvents []*chain.TransferEvent
	if forceSerialExecution {
		allEvents = getEventsFromBlocksSerial(evmChain, blockFrom, blockTo, onlyThisTokenAddress)
	} else {
		allEvents = getEventsFromBlocksConcurrent(evmChain, blockFrom, blockTo, onlyThisTokenAddress)
	}

	// Prepare block master data, first find out what blocks need master data
	uniqueBlocksMap := blocks.BuildUniqueBlocksFromEvents(allEvents)
	// then populate the master data for each block
	if forceSerialExecution {
		getBlockMasterDataSerial(evmChain, uniqueBlocksMap)
	} else {
		uniqueBlocksMap = getBlockMasterDataConcurrent(evmChain, uniqueBlocksMap)
	}
	allEvents = enrichAllEventsWithTimeEstimates(allEvents, uniqueBlocksMap)
	return allEvents
}

// enrichAllEventsWithTimeEstimates returns the enriched slice with all time fields added,
// plus the maximum timestamp seen on any transaction
func enrichAllEventsWithTimeEstimates(events []*chain.TransferEvent,
	uniqueBlocksMap blocks.BlockMap) []*chain.TransferEvent {

	addressTimeFirstSeenMap := make(map[common.Address]time.Time, 0)

	for _, event := range events {
		txTimestamp := blocks.GetEstimatedBlockTimeForTransaction(
			uniqueBlocksMap,
			event.BlockNumber,
			event.TxIndex).Round(time.Millisecond)
		event.TransactionTimestampEstimate = txTimestamp

		// Store the first seen times
		// For the FROM address
		addressTimeFirstSeen, exists := addressTimeFirstSeenMap[event.LogAddressFrom]
		if !exists {
			addressTimeFirstSeenMap[event.LogAddressFrom] = event.TransactionTimestampEstimate
		} else {
			if addressTimeFirstSeen.After(event.TransactionTimestampEstimate) {
				addressTimeFirstSeenMap[event.LogAddressFrom] = event.TransactionTimestampEstimate
			}
		}
		// For the TO address
		addressTimeFirstSeen, exists = addressTimeFirstSeenMap[event.LogAddressTo]
		if !exists {
			addressTimeFirstSeenMap[event.LogAddressTo] = event.TransactionTimestampEstimate
		} else {
			if addressTimeFirstSeen.After(event.TransactionTimestampEstimate) {
				addressTimeFirstSeenMap[event.LogAddressTo] = event.TransactionTimestampEstimate
			}
		}
	}

	// Write back the first seen times for the FROM and TO addresses
	for _, event := range events {
		event.LogAddressFromFirstSeen, _ = addressTimeFirstSeenMap[event.LogAddressFrom]
		event.LogAddressToFirstSeen, _ = addressTimeFirstSeenMap[event.LogAddressTo]
	}

	// For preparing test data
	//fmt.Print("[]*Event{")
	//for i, event := range events {
	//	if i > 0 {
	//		fmt.Print(", ")
	//	}
	//	fmt.Printf("&Event{%#v}", *event)
	//}
	//fmt.Println("}")

	// Add index-equivalents of the timestamps
	timeToIndexMap := GetTimeToIndexMap(events)
	for _, event := range events {
		event.TransactionTimestampEstimateIndex, _ = timeToIndexMap[event.TransactionTimestampEstimate]
		event.LogAddressFromFirstSeenIndex, _ = timeToIndexMap[event.LogAddressFromFirstSeen]
		event.LogAddressToFirstSeenIndex, _ = timeToIndexMap[event.LogAddressToFirstSeen]
	}
	return events
}

// GetTimeToIndexMap takes the events.TransactionTimestampEstimate field, makes a unique list of them, and assigns an
// index counter to them 0, 1, 2 etc.
func GetTimeToIndexMap(events []*chain.TransferEvent) (timeToIndexMap map[time.Time]uint) {
	// Create a slice to store unique event.TransactionTimestampEstimate values
	uniqueTimestamps := make([]time.Time, 0)

	// Create a map to store the index of each unique timestamp
	timestampIndexMap := make(map[time.Time]uint)

	// Loop through the events to collect unique timestamps
	for _, event := range events {
		if _, ok := timestampIndexMap[event.TransactionTimestampEstimate]; !ok {
			uniqueTimestamps = append(uniqueTimestamps, event.TransactionTimestampEstimate)
			timestampIndexMap[event.TransactionTimestampEstimate] = uint(len(uniqueTimestamps) - 1)
		}
	}

	// Sort the unique timestamps slice in ascending order
	sort.Slice(uniqueTimestamps, func(i, j int) bool {
		return uniqueTimestamps[i].Before(uniqueTimestamps[j])
	})

	// Create a map to store the index of each unique timestamp
	timeToIndexMap = make(map[time.Time]uint)
	for _, timestamp := range uniqueTimestamps {
		timeToIndexMap[timestamp] = timestampIndexMap[timestamp]
	}

	return timeToIndexMap
}
