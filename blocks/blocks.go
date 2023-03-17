package blocks

import (
	"context"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum"
	"math/big"
	"time"
)

func GetEstimatedBlockTimeForTransaction(uniqueBlocksMap BlockMap, blockNumber uint64, txIndex uint) time.Time {
	blockMapValue, exists := uniqueBlocksMap[BlockKey{BlockNumber: blockNumber}]
	if !exists {
		logr.Error.Panicf("Block %v missing in internal uniqueBlocksMap", blockNumber)
	}
	blockDurationSeconds := float64(12) // arbitrary estimate
	proportionThroughBlockInSeconds := (float64(txIndex) / float64(blockMapValue.TransactionCount)) * blockDurationSeconds
	timeToAdd := time.Duration(proportionThroughBlockInSeconds * float64(time.Second))
	estimatedTime := blockMapValue.BlockTimestamp.Add(timeToAdd)
	return estimatedTime
}

func GetBlockFromChain(client ethereum.ChainReader, blockKey BlockKey) (
	blockData BlockDataFromSource, err error) {
	blockNumber := big.NewInt(int64(blockKey.BlockNumber))

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return BlockDataFromSource{}, err
	}

	bn := block.Number().Uint64()
	bts := time.Unix(int64(block.Time()), 0)
	txc := uint(len(block.Transactions()))

	return BlockDataFromSource{
		BlockNumber:      bn,
		BlockTimestamp:   bts,
		TransactionCount: txc,
	}, nil
}

func BuildUniqueBlocksFromEvents(allEvents []*chain.TransferEvent) (uniqueBlocksMap BlockMap) {
	uniqueBlocksMap = make(BlockMap)
	// Iterate through the events and get unique blocks
	for _, event := range allEvents {
		_, exists := uniqueBlocksMap[BlockKey{BlockNumber: event.BlockNumber}]
		if !exists {
			// New block not seen before
			uniqueBlocksMap[BlockKey{BlockNumber: event.BlockNumber}] = BlockMapValue{
				BlockTimestamp:   time.Time{},
				TransactionCount: 0,
			}
		}
	}
	return uniqueBlocksMap
}
