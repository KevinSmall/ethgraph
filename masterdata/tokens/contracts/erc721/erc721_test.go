package erc721

import (
	"github.com/KevinSmall/ethgraph/test"
	"testing"
)

func TestErc721(t *testing.T) {

	client, address, err := test.GetMockClient()
	if err != nil {
		t.Fatalf("Unable to create mock client %s", err)
	}
	// there isn't a valid token contract here, but we can create this anyway
	_, err = NewErc721(address, client)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Called NewErc721 without error")
	}
}
