// Package chain reads an EVM chain and returns ERC20, ERC721 and ERC1155 transfer event logs.
// The data returned is []*TransferEvent, and it is cleansed and enhanced with TransferType to
// distinguish the token types.
//   - No graph-related logic is applied here, this is pure event log handling.
//   - For performance reasons, no additional chain reads are allowed here. It must only be
//     a single log query hitting the chain.
package chain

import (
	"context"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetTransferEventsByBlock(client *ethclient.Client, blockNumberInt uint64,
	onlyThisTokenAddress string) ([]*TransferEvent, error) {

	// Notes on how to use FilterQuery
	// Filter info:
	//
	//	    BlockHash *common.Hash     // used by eth_getLogs, return logs only from block with this hash
	//		FromBlock *big.Int         // beginning of the queried range, nil means genesis block
	//		ToBlock   *big.Int         // end of the range, nil means latest block
	//		Addresses []common.Address // restricts matches to events created by specific contracts
	//
	//		The Topic list restricts matches to particular event topics. Each event has a list
	//		of topics. Topics match a prefix of that list. An empty element slice matches any
	//		topic. Non-empty elements represent an alternative that matches any of the
	//		contained topics.
	//
	//		Examples:
	//		{} or nil          matches any topic list
	//		{{A}}              matches topic A in first position
	//		{{}, {B}}          matches any topic in first position AND B in second position
	//		{{A}, {B}}         matches topic A in first position AND B in second position
	//		{{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position

	// the block number to retrieve transactions from
	blockNumberBig := big.NewInt(int64(blockNumberInt))

	// create a filter query for the specified block
	var addresses []common.Address
	if onlyThisTokenAddress != "" {
		addresses = []common.Address{common.HexToAddress(onlyThisTokenAddress)}
	}
	query := ethereum.FilterQuery{
		Addresses: addresses,
		FromBlock: blockNumberBig,
		ToBlock:   blockNumberBig,
		Topics: [][]common.Hash{{common.HexToHash(transferEventKeccakTokens),
			common.HexToHash(transferEventKeccakHybridSingle),
			common.HexToHash(transferEventKeccakHybridBatch),
		}},
	}
	// retrieve the logs matching the filter query
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}
	logr.Trace.Println("Block: ", blockNumberInt, "Logs found: ", len(logs))

	events := logsToEvents(logs)
	return events, nil
}

func logsToEvents(logs []types.Log) []*TransferEvent {

	var events []*TransferEvent

	// iterate through the logs
	for _, log := range logs {
		// log.Topics contain only indexed logs
		if len(log.Topics) == 3 {
			//-----------------------------------------------------------------
			// ERC20 (log.Data contains value uint256)
			//-----------------------------------------------------------------
			from := common.HexToAddress(log.Topics[1].Hex())
			to := common.HexToAddress(log.Topics[2].Hex())
			value := new(big.Int).SetBytes(log.Data)
			event := TransferEvent{
				BlockNumber:       log.BlockNumber,
				TxHash:            log.TxHash,
				TxIndex:           log.TxIndex,
				TransferType:      ERC20,
				LogIndex:          log.Index,
				LogAddressFrom:    from,
				LogAddressTo:      to,
				LogTokenValue:     *value,
				LogNftId:          "",
				LogOperator:       zeroAddress,
				LogEmitterAddress: log.Address,
			}
			events = append(events, &event)

		} else if len(log.Topics) == 4 {
			topic0 := log.Topics[0].Hex()
			if topic0 == transferEventKeccakTokens {
				//-----------------------------------------------------------------
				// ERC721 (log.Data is empty)
				//-----------------------------------------------------------------
				from := common.HexToAddress(log.Topics[1].Hex())
				to := common.HexToAddress(log.Topics[2].Hex())
				nftId := log.Topics[3].Big().String()
				event := TransferEvent{
					BlockNumber:       log.BlockNumber,
					TxHash:            log.TxHash,
					TxIndex:           log.TxIndex,
					TransferType:      ERC721,
					LogIndex:          log.Index,
					LogAddressFrom:    from,
					LogAddressTo:      to,
					LogTokenValue:     *big.NewInt(0),
					LogNftId:          nftId,
					LogOperator:       zeroAddress,
					LogEmitterAddress: log.Address,
				}
				events = append(events, &event)

			} else if topic0 == transferEventKeccakHybridSingle {
				//-----------------------------------------------------------------
				// ERC1155 Single (log.Data contains id, value pair (both uint256), where value represents quantity)
				//-----------------------------------------------------------------
				id, value := DecodeDataForErc1155Single(log.Data)
				operator := common.HexToAddress(log.Topics[1].Hex())
				from := common.HexToAddress(log.Topics[2].Hex())
				to := common.HexToAddress(log.Topics[3].Hex())
				event := TransferEvent{
					BlockNumber:       log.BlockNumber,
					TxHash:            log.TxHash,
					TxIndex:           log.TxIndex,
					TransferType:      ERC1155_SINGLE,
					LogIndex:          log.Index,
					LogAddressFrom:    from,
					LogAddressTo:      to,
					LogTokenValue:     *value,
					LogNftId:          id.String(),
					LogOperator:       operator,
					LogEmitterAddress: log.Address,
				}
				events = append(events, &event)

			} else if topic0 == transferEventKeccakHybridBatch {
				//-----------------------------------------------------------------
				// ERC1155 Batch (log.Data contains []id, []value lists of uint256, where values represents quantities)
				//-----------------------------------------------------------------
				ids, values := DecodeDataForErc1155Batch(log.Data)
				for i := 0; i < len(ids); i++ {
					operator := common.HexToAddress(log.Topics[1].Hex())
					from := common.HexToAddress(log.Topics[2].Hex())
					to := common.HexToAddress(log.Topics[3].Hex())
					event := TransferEvent{
						BlockNumber:       log.BlockNumber,
						TxHash:            log.TxHash,
						TxIndex:           log.TxIndex,
						TransferType:      ERC1155_BATCH,
						LogIndex:          log.Index,
						LogAddressFrom:    from,
						LogAddressTo:      to,
						LogTokenValue:     *values[i],
						LogNftId:          ids[i].String(),
						LogOperator:       operator,
						LogEmitterAddress: log.Address,
					}
					events = append(events, &event)
				}

			} else {
				// Unknown topic0 value for a 4-topic log
				logr.Trace.Printf("Transaction %s\n has unrecognised log topic0 %v", log.TxHash.Hex(), topic0)
				continue
			}
		} else {
			// Unknown log topic count
			logr.Trace.Printf("Transaction %s\n has unrecognised log topic count %v", log.TxHash.Hex(), len(log.Topics))
			continue
		}
	} // end log loop
	return events
}
