package services

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/blocks"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/work"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

// getBlockMasterDataSerial updates uniqueBlocksMap with block timestamp (non-concurrent version)
func getBlockMasterDataSerial(evmChain chain.EvmClient, uniqueBlocksMap blocks.BlockMap) {
	fmt.Printf("Getting block times...")

	for blockMapKey, _ := range uniqueBlocksMap {
		// Arbitrary throttle in serial mode, some chain providers can throttle calls
		time.Sleep(throttleHttpDelayMilliseconds)

		blockData, err := blocks.GetBlockFromChain(evmChain.Client, blockMapKey)
		if err != nil {
			logr.Error.Panicln(err)
		} else {
			uniqueBlocksMap[blockMapKey] = blocks.BlockMapValue{
				BlockTimestamp:   blockData.BlockTimestamp,
				TransactionCount: blockData.TransactionCount}
		}
		fmt.Printf(".")
	}
	fmt.Printf("done.\n")
}

// getBlockAttrWorker is to hold the work that needs done
type getBlockAttrWorker struct {
	chainId     string
	url         string
	blockNumber uint64
	resultChan  chan blocks.BlockDataFromSource
}

// Task is the work that needs done and fulfills the Pool's Worker interface
func (w *getBlockAttrWorker) Task() {

	// connect to the client
	client, err := ethclient.Dial(w.url)
	if err != nil {
		logr.Warning.Println("TASK ERROR getBlockAttrWorker ", err)
		return
	}
	blockData, err := blocks.GetBlockFromChain(client, blocks.BlockKey{BlockNumber: w.blockNumber})
	if err != nil {
		logr.Warning.Println("TASK ERROR GetBlockFromChain ", err)
		return
	}
	w.resultChan <- blockData
}

func getBlockMasterDataConcurrent(evmChain chain.EvmClient, uniqueBlocksMap blocks.BlockMap) blocks.BlockMap {

	updatedUniqueBlocksMap := make(blocks.BlockMap, 0)
	pool := work.New(10_000)
	resultsChan := make(chan blocks.BlockDataFromSource, len(uniqueBlocksMap))

	fmt.Printf("Getting block times...")

	for blockKey, _ := range uniqueBlocksMap {
		worker := &getBlockAttrWorker{
			chainId:     evmChain.ChainId,
			url:         evmChain.Url,
			blockNumber: blockKey.BlockNumber,
			resultChan:  resultsChan,
		}
		pool.Run(worker) // blocks main thread if nobody able to pick up the work
	}

	// Wait for workers to finish
	for i := 0; i < len(uniqueBlocksMap); i++ {
		blockAttrFromChain := <-resultsChan
		updatedUniqueBlocksMap[blocks.BlockKey{BlockNumber: blockAttrFromChain.BlockNumber}] =
			blocks.BlockMapValue{
				BlockTimestamp:   blockAttrFromChain.BlockTimestamp,
				TransactionCount: blockAttrFromChain.TransactionCount,
			}
	}

	pool.Shutdown()

	fmt.Printf("done.\n")
	return updatedUniqueBlocksMap
}
