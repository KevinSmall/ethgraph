package conv

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func TestFormatBlockNumberWithUnderscores(t *testing.T) {

	expected := "12_123_456"
	actual := PrettyBlockNumberWithUnderscores(uint64(12_123_456))
	if actual == expected {
		t.Logf("Correctly formatted to %s", actual)
	} else {
		t.Errorf("Failed to format to %s got %s", expected, actual)
	}
}

func TestSafeScaleTokenValue(t *testing.T) {

	testCases := []struct {
		tokenValue *big.Int
		decimals   int
		expected   float64
	}{
		{big.NewInt(123456789), 6, 123.46},
		{big.NewInt(123456789), 0, 123456789},
		{big.NewInt(0), 0, 0.00},
		{big.NewInt(0), -1, 0.00},
	}

	for _, tc := range testCases {
		result := SafeScaleTokenValue(tc.tokenValue, tc.decimals)
		if result != tc.expected {
			t.Errorf("SafeScaleTokenValue(%v, %v) = %v, expected %v. Fail",
				tc.tokenValue, tc.decimals, result, tc.expected)
		} else {
			t.Logf("SafeScaleTokenValue(%v, %v) = %v, expected %v. Pass",
				tc.tokenValue, tc.decimals, result, tc.expected)
		}
	}
}

func TestPrettyTrimAddress(t *testing.T) {
	testCases := []struct {
		addr     string
		expected string
	}{
		{"0x0000007C4ed4134E5Ae1f825d028559f267767df", "0x7C4ED4134E5Ae1f825d028559F267767dF"},
		{"0x0000000000000000000000000000000000000000", ""},
	}

	for _, tc := range testCases {
		addr := common.HexToAddress(tc.addr)
		actual := PrettyTrimAddress(addr)
		if actual == tc.expected {
			t.Logf("Pass: expected %s, got %s", tc.expected, actual)
		} else {
			t.Errorf("Fail: expected %s, got %s", tc.expected, actual)
		}
	}
}
