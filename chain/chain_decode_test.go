package chain

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func TestDecodeDataForErc1155Single(t *testing.T) {
	testCases := []struct {
		logDataHex  string
		expectedId  *big.Int
		expectedVal *big.Int
	}{
		{
			logDataHex:  "0x00000000000000000000000000000000000000000000000000000000000029810000000000000000000000000000000000000000000000000000000000000001",
			expectedId:  big.NewInt(10625),
			expectedVal: big.NewInt(1),
		},
		{
			logDataHex:  "0x00000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000004",
			expectedId:  big.NewInt(5),
			expectedVal: big.NewInt(4),
		},
		{
			logDataHex:  "0x01",
			expectedId:  big.NewInt(0),
			expectedVal: big.NewInt(0),
		},
	}

	for _, tc := range testCases {
		logData, err := hex.DecodeString(tc.logDataHex[2:]) // remove "0x" prefix
		if err != nil {
			t.Fatalf("failed to decode hex string: %v", err)
		}

		id, val := DecodeDataForErc1155Single(logData)
		if id.Cmp(tc.expectedId) != 0 {
			t.Errorf("expected ID: %s, got: %s", tc.expectedId.String(), id.String())
		}
		if val.Cmp(tc.expectedVal) != 0 {
			t.Errorf("expected value: %s, got: %s", tc.expectedVal.String(), val.String())
		}
	}
}

func TestDecodeDataForErc1155Batch(t *testing.T) {

	testCases := []struct {
		logDataHex   string
		expectedIds  []*big.Int
		expectedVals []*big.Int
	}{
		{
			// See https://etherscan.io/tx/0x6c0364729a9c06cb39a0d45850a7f0187fbc97b35b932841ffbed37a278041a9#eventlog
			logDataHex: "0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000051700000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
			expectedIds: []*big.Int{
				big.NewInt(1303),
			},
			expectedVals: []*big.Int{
				big.NewInt(1),
			},
		},
		{
			// See https://etherscan.io/tx/0xbb34bbf7a7d4f2186e549c9fc3afbafc597ef630959984c46c3827a447d29b04#eventlog
			logDataHex: "0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000029140000000000000000000000000000000000000000000000000000000000002933000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
			expectedIds: []*big.Int{
				big.NewInt(10516),
				big.NewInt(10547),
			},
			expectedVals: []*big.Int{
				big.NewInt(1),
				big.NewInt(1),
			},
		},
		{
			logDataHex:   "0x01",
			expectedIds:  make([]*big.Int, 0),
			expectedVals: make([]*big.Int, 0),
		},
	}

	for _, tc := range testCases {
		logData, err := hex.DecodeString(tc.logDataHex[2:]) // remove "0x" prefix
		if err != nil {
			t.Fatalf("failed to decode hex string: %v", err)
		}
		ids, values := DecodeDataForErc1155Batch(logData)

		// compare expectedIds with ids
		if len(ids) != len(tc.expectedIds) {
			t.Errorf("expected %d ids, but got %d", len(tc.expectedIds), len(ids))
			continue
		}

		for i, id := range ids {
			if id.Cmp(tc.expectedIds[i]) != 0 {
				t.Errorf("expected id[%d] to be %v, but got %v", i, tc.expectedIds[i], id)
			} else {
				t.Logf("expected id[%d] to be %v and got %v ok", i, tc.expectedIds[i], id)
			}
		}

		// compare expectedVals with values
		if len(values) != len(tc.expectedVals) {
			t.Errorf("expected %d values, but got %d", len(tc.expectedVals), len(values))
			continue
		}

		for i, value := range values {
			if value.Cmp(tc.expectedVals[i]) != 0 {
				t.Errorf("expected value[%d] to be %v, but got %v", i, tc.expectedVals[i], value)
			} else {
				t.Logf("expected value[%d] to be %v, and got %v ok", i, tc.expectedVals[i], value)
			}
		}
	}
}
