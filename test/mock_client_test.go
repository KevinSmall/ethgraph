package test

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestGetMockClient(t *testing.T) {
	client, address, err := GetMockClient()

	if err != nil {
		t.Fatalf("GetMockClient returned an error: %v", err)
	}

	if client == nil {
		t.Fatalf("GetMockClient returned an invalid client: %T", client)
	}

	if address == (common.Address{}) {
		t.Fatalf("GetMockClient returned an empty address: %v", address.Hex())
	}
}
