package services

import (
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
	"time"
)

const yearOfTransactionTimeEstimate = 2000

var testData = []*chain.TransferEvent{
	{BlockNumber: 0x1a00914,
		BlockTimestamp:                    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TransactionTimestampEstimate:      time.Date(yearOfTransactionTimeEstimate, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		TransactionTimestampEstimateIndex: 0x0, TxHash: common.HexToHash("0x244c8173556a0e15b26db7a7729a66cca0f8689a19a7b6e709ccfe47096074a0"),
		TxIndex:                      0x2,
		TransferType:                 "ERC20",
		LogIndex:                     0x8,
		LogAddressFrom:               common.HexToAddress("0x08f47FFbB40aAE4662eB5f4F284f2d056Deb0dc2"),
		LogAddressTo:                 common.Address{},
		LogAddressFromFirstSeen:      time.Date(yearOfTransactionTimeEstimate, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		LogAddressToFirstSeen:        time.Date(yearOfTransactionTimeEstimate, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		LogAddressFromFirstSeenIndex: 0x0,
		LogAddressToFirstSeenIndex:   0x0,
		LogTokenValue:                *big.NewInt(0),
		LogNftId:                     "",
		LogOperator:                  common.Address{},
		LogEmitterAddress:            common.HexToAddress("0x472361d3cA5F49c8E633FB50385BfaD1e018b445"),
	},
}

func TestGetTimeToIndexMap(t *testing.T) {

	timeToIndexMap := GetTimeToIndexMap(testData)
	if len(timeToIndexMap) != 1 {
		t.Errorf("time to index map entries expected 1 got %v", len(timeToIndexMap))
	}
	for timeStamp, timeIndex := range timeToIndexMap {
		if yearOfTransactionTimeEstimate != timeStamp.Year() {
			t.Errorf("expected year %v got %v", yearOfTransactionTimeEstimate, timeStamp.Year())
		}
		if timeIndex != 0 {
			t.Errorf("expected index %v got %v", 1, timeIndex)
		}
	}
}
