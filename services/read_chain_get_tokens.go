package services

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/masterdata/tokens"
	"github.com/KevinSmall/ethgraph/work"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

// getTokenMasterDataForMissingTokensSerial is the non-concurrent version
func getTokenMasterDataForMissingTokensSerial(
	evmChain chain.EvmClient,
	tokensWithoutMasterData map[common.Address]*tokens.AddressMapValue) (
	tokenMapToAdd map[string]tokens.TokenDataFromSource) {

	tokenMapToAdd = make(map[string]tokens.TokenDataFromSource, 0)
	fmt.Printf("Getting tokens ")
	for address, addressMapValue := range tokensWithoutMasterData {
		fmt.Printf(".")
		// Arbitrary throttle in serial mode, some chain providers can throttle calls
		time.Sleep(throttleHttpDelayMilliseconds)
		tokenDataFromChain := tokens.GetTokenFromChain(evmChain.ChainId, evmChain.Client, address, addressMapValue.TransferType)
		tokenMapToAdd[tokenDataFromChain.TokenAddress] = tokenDataFromChain
	}
	fmt.Printf("done.\n")
	return tokenMapToAdd
}

// getTokenWorker is to hold the work that needs done
type getTokenWorker struct {
	chainId      string
	url          string
	address      common.Address
	transferType string
	resultChan   chan tokens.TokenDataFromSource
}

// Task is the work that needs done and fulfills the Pool's Worker interface
func (w *getTokenWorker) Task() {

	// connect to the client
	client, err := ethclient.Dial(w.url)
	if err != nil {
		logr.Info.Println("TASK ERROR ", err)
		return
	}
	tokenDataFromChain := tokens.GetTokenFromChain(w.chainId, client, w.address, w.transferType)
	w.resultChan <- tokenDataFromChain
}

func getTokenMasterDataForMissingTokensConcurrent(
	evmChain chain.EvmClient,
	tokensWithoutMasterData map[common.Address]*tokens.AddressMapValue) (
	tokenMapToAdd map[string]tokens.TokenDataFromSource) {

	tokenMapToAdd = make(map[string]tokens.TokenDataFromSource, 0)
	pool := work.New(10_000)
	resultsChan := make(chan tokens.TokenDataFromSource, len(tokensWithoutMasterData))

	fmt.Printf("Getting tokens...")

	for address, addressMapValue := range tokensWithoutMasterData {

		worker := &getTokenWorker{
			chainId:      evmChain.ChainId,
			url:          evmChain.Url,
			address:      address,
			transferType: addressMapValue.TransferType,
			resultChan:   resultsChan,
		}
		pool.Run(worker) // blocks main thread if nobody able to pick up the work
	}

	// Wait for workers to finish
	for i := 0; i < len(tokensWithoutMasterData); i++ {
		tokenDataFromChain := <-resultsChan
		tokenMapToAdd[tokenDataFromChain.TokenAddress] = tokenDataFromChain
	}

	pool.Shutdown()

	fmt.Printf("done.\n")
	return tokenMapToAdd
}
