package services

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/conv"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/work"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

// getEventsFromBlocksSerial is the non-concurrent version
func getEventsFromBlocksSerial(evmChain chain.EvmClient, blockFrom uint64, blockTo uint64, onlyThisTokenAddress string) []*chain.TransferEvent {
	var allEvents []*chain.TransferEvent
	fmt.Printf("Getting blocks ")
	blockCount := int(blockTo-blockFrom) + 1
	for i := 0; i < blockCount; i++ {
		// Arbitrary throttle in serial mode, some chain providers can throttle calls
		time.Sleep(throttleHttpDelayMilliseconds)

		events, err := chain.GetTransferEventsByBlock(evmChain.Client, blockFrom, onlyThisTokenAddress)
		if err != nil {
			logr.Error.Panicln(err)
		}
		blockFrom++
		allEvents = append(allEvents, events...)
		fmt.Printf(".")
	}
	fmt.Printf("done.\n")
	return allEvents
}

// getBlockWorker is to hold the work that needs done
type getBlockWorker struct {
	url                  string
	blockNumber          uint64
	onlyThisTokenAddress string
	resultChan           chan []*chain.TransferEvent
}

// Task is the work that needs done and fulfills the Pool's Worker interface
func (w *getBlockWorker) Task() {

	// connect to the client
	client, err := ethclient.Dial(w.url)
	if err != nil {
		logr.Info.Println("TASK ERROR ", err)
		return
	}
	events, err := chain.GetTransferEventsByBlock(client, w.blockNumber, w.onlyThisTokenAddress)
	if err != nil {
		logr.Error.Panicln(err)
	}
	w.resultChan <- events
}

func getEventsFromBlocksConcurrent(evmChain chain.EvmClient, blockFrom uint64, blockTo uint64, onlyThisTokenAddress string) []*chain.TransferEvent {

	allEvents := make([]*chain.TransferEvent, 0)

	blockCount := int(blockTo-blockFrom) + 1
	pool := work.New(1_000)
	resultsChan := make(chan []*chain.TransferEvent, blockCount)

	fmt.Printf("Getting blocks...")

	for i := 0; i < blockCount; i++ {
		blockNumber := blockFrom + uint64(i)
		worker := &getBlockWorker{
			url:                  evmChain.Url,
			blockNumber:          blockNumber,
			onlyThisTokenAddress: onlyThisTokenAddress,
			resultChan:           resultsChan,
		}
		pool.Run(worker) // blocks main thread if nobody able to pick up the work
	}

	// Wait for workers to finish
	for i := 0; i < blockCount; i++ {
		events := <-resultsChan
		allEvents = append(allEvents, events...)
	}

	pool.Shutdown()

	fmt.Printf("done.\n")
	return allEvents
}

func GetLatestBlockNumber(url string) {
	start := time.Now()

	evmChain, err := chain.CreateEvmClient(url)
	if err != nil {
		logr.Error.Panicln(err)
	}
	logr.Info.Printf("Connecting to: %s with ChainId: %s", evmChain.Name, evmChain.ChainId)

	// PrintSummary the latest block number.
	logr.Info.Println("Latest block number: ", conv.PrettyBlockNumberWithUnderscores(evmChain.LatestBlockNumber))

	elapsed := time.Since(start)
	logr.Info.Printf("Runtime: %s", elapsed)
}
