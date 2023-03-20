package chain

import (
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"
)

// Transfer types

const ERC20 string = "ERC20"
const ERC721 string = "ERC721"
const ERC1155_SINGLE string = "ERC1155_SINGLE"
const ERC1155_BATCH string = "ERC1155_BATCH"

type TransferEvent struct {
	// Block number
	BlockNumber uint64

	// Block timestamp from chain (enrichment)
	BlockTimestamp time.Time

	// Transaction timestamp estimated as a proportion through whole block time, based transaction index inside block (enrichment)
	TransactionTimestampEstimate time.Time

	// Index at which first seen in dataset, equivalent of TransactionTimestampEstimate
	TransactionTimestampEstimateIndex uint

	// Transaction hash
	TxHash common.Hash

	// Transaction index in the block
	TxIndex uint

	// Log transfer event is ERC20, ERC721, ERC1155
	TransferType string

	// Log index (can be many logs for one transaction)
	LogIndex uint

	// Transfer event log address from
	LogAddressFrom common.Address

	// Transfer event log address to
	LogAddressTo common.Address

	// Transfer event log address from, time first seen in dataset (enrichment)
	LogAddressFromFirstSeen time.Time

	// Transfer event log address to, time first seen in dataset (enrichment)
	LogAddressToFirstSeen time.Time

	// Transfer event log address from, index at which first seen in dataset (enrichment)
	LogAddressFromFirstSeenIndex uint

	// Transfer event log address to, index at which first seen in dataset (enrichment)
	LogAddressToFirstSeenIndex uint

	// Transfer event log, the value of tokens transferred, filled for ERC20, ERC1155, else 0
	LogTokenValue big.Int

	// Transfer event log, the NFT id transferred (natively uint256 on Ethereum), filled for ERC721, ERC1155, else empty
	LogNftId string

	// Transfer event log, for ERC1155 this holds the operator (the address of an account/contract that is approved to make the transfer), else 0x
	LogOperator common.Address

	// The address that emitted the transfer event log
	LogEmitterAddress common.Address
}

func (event *TransferEvent) Print(title string) {
	logr.Info.Println("------------", title, " -----------")
	logr.Info.Println("BlockNumber:", event.BlockNumber)
	logr.Info.Println("BlockTimestamp:", event.BlockTimestamp.Format("2006-01-02 15:04:05"))
	logr.Info.Println("TxHash:", event.TxHash.Hex())
	logr.Info.Println("TxIndex:", event.TxIndex)
	logr.Info.Println("TransferType:", event.TransferType)
	logr.Info.Println("LogIndex:", event.LogIndex)
	logr.Info.Println("LogAddressFrom:", event.LogAddressFrom.Hex())
	logr.Info.Println("LogAddressTo:", event.LogAddressTo.Hex())
	logr.Info.Println("LogTokenValue:", event.LogTokenValue.String())
	logr.Info.Println("LogNftId:", event.LogNftId)
	logr.Info.Println("LogEmitterAddress:", event.LogEmitterAddress.Hex())
	logr.Info.Println("---------------------------------------------------")
}
