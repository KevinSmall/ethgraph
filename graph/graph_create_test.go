package graph

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
	"time"
)

// testData contains two transfer events.
//
//	|<- transfer event ->|   |<- transfer event ->|    <= the two transfer events from chain
//	an0    ->  te0   ->   an1   ->    te1  ->   an2    <= the graph generated as a result of above
//
// The graph consists of:
//
//	  5 nodes:
//		   3 address nodes an0, an1, an2 in above schematic (0x08f4..c2, zero address, 0xA010..01)
//		   2 nodes for transfer events, te0 and te1 in above schematic
//		 4 edges:
//		   The -> in the above schematic
var testData = []*chain.TransferEvent{
	{BlockNumber: 0x1a00914,
		BlockTimestamp:                    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TransactionTimestampEstimate:      time.Date(2022, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		TransactionTimestampEstimateIndex: 0,
		TxHash:                            common.HexToHash("0x244c8173556a0e15b26db7a7729a66cca0f8689a19a7b6e709ccfe47096074a0"),
		TxIndex:                           2,
		TransferType:                      "ERC20",
		LogIndex:                          8,
		LogAddressFrom:                    common.HexToAddress("0x08f47FFbB40aAE4662eB5f4F284f2d056Deb0dc2"),
		LogAddressTo:                      zeroAddress,
		LogAddressFromFirstSeen:           time.Date(2022, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		LogAddressToFirstSeen:             time.Date(2022, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		LogAddressFromFirstSeenIndex:      0,
		LogAddressToFirstSeenIndex:        0,
		LogTokenValue:                     *big.NewInt(0),
		LogNftId:                          "",
		LogOperator:                       zeroAddress,
		LogEmitterAddress:                 common.HexToAddress("0x472361d3cA5F49c8E633FB50385BfaD1e018b445")},
	{BlockNumber: 0x1a00914,
		BlockTimestamp:                    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
		TransactionTimestampEstimate:      time.Date(2022, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		TransactionTimestampEstimateIndex: 0,
		TxHash:                            common.HexToHash("0x244c8173556a0e15b26db7a7729a66cca0f8689a19a7b6e709ccfe47096074a0"),
		TxIndex:                           2,
		TransferType:                      "ERC20",
		LogIndex:                          9,
		LogAddressFrom:                    common.HexToAddress("0xA0107FFbB40aAE4662eB5f4F284f2d056Deb0d01"),
		LogAddressTo:                      zeroAddress,
		LogAddressFromFirstSeen:           time.Date(2022, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		LogAddressToFirstSeen:             time.Date(2022, time.September, 22, 0, 12, 43, 145000000, time.UTC),
		LogAddressFromFirstSeenIndex:      0x0,
		LogAddressToFirstSeenIndex:        0x0,
		LogTokenValue:                     *big.NewInt(0),
		LogNftId:                          "",
		LogOperator:                       zeroAddress,
		LogEmitterAddress:                 common.HexToAddress("0xFf1489227BbAAC61a9209A08929E4c2a526DdD17")},
}

func TestCreateGraph(t *testing.T) {
	wot, creationResult := CreateGraph("HelloWorld", testData)
	fmt.Println(wot)
	if creationResult.Nodes != 5 {
		t.Errorf("Expected 5 nodes, got %v", creationResult.Nodes)
	}
	if creationResult.Edges != 4 {
		t.Errorf("Expected 4 edges, got %v", creationResult.Edges)
	}
	if creationResult.Events != 2 {
		t.Errorf("Expected 2 events, got %v", creationResult.Events)
	}
	t.Logf("Created %v nodes, %v edges, %v events",
		creationResult.Nodes, creationResult.Edges, creationResult.Events)
}
