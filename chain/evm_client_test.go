package chain

import (
	"testing"
)

func TestGetChainName(t *testing.T) {

	chainId := "1"
	expectedName := "ethereum"
	actualName := getChainName(chainId)

	if actualName == expectedName {
		t.Logf("chainId %s expected name %s actual name %s", chainId, expectedName, actualName)
	} else {
		t.Fatalf("chainId %s expected name %s actual name %s", chainId, expectedName, actualName)
	}
}
