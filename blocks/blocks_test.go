package blocks

import (
	"github.com/KevinSmall/ethgraph/test"
	"testing"
)

func TestGetBlockFromChain(t *testing.T) {

	client, _, err := test.GetMockClient()
	if err != nil {
		t.Fatalf("Unable to create mock client %s", err)
	}
	blockData, err := GetBlockFromChain(client, BlockKey{0})
	if err != nil {
		t.Fatal(err)
	} else if blockData.BlockNumber == 0 {
		t.Logf("Got blockData without error")
	} else {
		t.Fatalf("Got blockData with unexpected blocknumber %v", blockData.BlockNumber)
	}
}
