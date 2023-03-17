package conv

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var zeroAddress = common.Address{}

func AddressToBigInt(addr common.Address) big.Int {
	var bi big.Int
	bi.SetBytes(addr[:])
	return bi
}

func SafeStringToInt(s string) int {
	u, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		u = int64(0)
	}
	return int(u)
}

// PrettyTrimAddress removes leading zeros in address
func PrettyTrimAddress(addr common.Address) (trimmed string) {
	if addr == zeroAddress {
		trimmed = ""
	} else {
		trimmed = "0x" + strings.TrimLeft(addr.Hex()[2:], "0")
	}
	return
}

// PrettyShortenAddress takes an address like 0xdA6c517C4ed4134E5Ae1f825d028559f267767df and
// returns 0xdA6c...df
func PrettyShortenAddress(addr string) string {
	if len(addr) <= 8 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-2:]
}

// PrettyBlockNumberWithUnderscores pretty prints a block number e.g. 16345454 as 16_345_454
func PrettyBlockNumberWithUnderscores(blockNumber uint64) string {
	blockNumberString := fmt.Sprint(blockNumber)
	blockNumberWithUnderscores := ""
	for i := len(blockNumberString) - 1; i >= 0; i-- {
		blockNumberWithUnderscores = string(blockNumberString[i]) + blockNumberWithUnderscores
		if (len(blockNumberString)-i)%3 == 0 && i != 0 {
			blockNumberWithUnderscores = "_" + blockNumberWithUnderscores
		}
	}
	return blockNumberWithUnderscores
}

// SafeBigIntToFloat converts a big.Int to a float or zero if any troubles
func SafeBigIntToFloat(n *big.Int) float64 {
	if !n.IsInt64() {
		return 0
	}
	x := n.Int64()
	if x > math.MaxInt64 || x < math.MinInt64 {
		return 0
	}
	return float64(x)
}

// SafeScaleTokenValue scales a token value by the given decimals and returns
// the value with two decimal places afterwards. So 123456789 and 6 decimals will
// return 123.46 (which is 123.456789 rounded up to 123.46). If anything goes wrong
// the expected will be 0.00.
func SafeScaleTokenValue(tokenValue *big.Int, decimals int) (resultFloat float64) {

	if decimals >= 0 && tokenValue.Cmp(big.NewInt(0)) != 0 {
		bigF := new(big.Float)
		bigF.SetInt(tokenValue)
		scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
		scaledValue := new(big.Float).Quo(bigF, new(big.Float).SetInt(scaleFactor))
		resultFloat, _ = scaledValue.Float64()
	}
	return math.Round(resultFloat*100) / 100
}
