package tokens

import (
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/masterdata/tokens/contracts/erc20"
	"github.com/KevinSmall/ethgraph/masterdata/tokens/contracts/erc721"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"strings"
)

func GetTokenFromChain(chainId string, client *ethclient.Client, tokenAddr common.Address, transferType string) (tokenData TokenDataFromSource) {
	var tokenName, tokenSymbol string
	var tokenDecimals uint8
	if transferType == chain.ERC20 {
		tokenName, tokenSymbol, tokenDecimals = getTokenERC20FromChain(client, tokenAddr)
	} else if transferType == chain.ERC721 {
		tokenName, tokenSymbol = getTokenERC721FromChain(client, tokenAddr)
	} else if transferType == chain.ERC1155_SINGLE || transferType == chain.ERC1155_BATCH {
		// ERC1155 decided not to include name or symbol
		// See Metadata Choices section in https://eips.ethereum.org/EIPS/eip-1155
		// Try anyway, and if it fails defaults will show anyway.
		tokenName, tokenSymbol = getTokenERC721FromChain(client, tokenAddr)
	} else {
		logr.Warning.Printf("Unknown transferType for address %s transferType %s.", tokenAddr.Hex(), transferType)
		return
	}
	tokenData = TokenDataFromSource{
		ChainId:      chainId,
		Name:         tokenName,
		Symbol:       tokenSymbol,
		Decimals:     int(tokenDecimals),
		TokenAddress: tokenAddr.Hex(),
	}
	return tokenData
}

func getTokenERC20FromChain(client *ethclient.Client, tokenAddr common.Address) (name string, symbol string, decimals uint8) {

	// Defaults shine through if any errors
	name = "Unknown"
	symbol = "UNKNOWN"
	decimals = 0

	// Create a new ERC20 instance using the token address and client
	erc20Instance, err := erc20.NewErc20(tokenAddr, client)
	if err == nil {
		name, err = erc20Instance.Name(nil)
		if err == nil {
			name = strings.ReplaceAll(name, ",", " ")
		}
		symbol, err = erc20Instance.Symbol(nil)
		if err == nil {
			symbol = strings.ReplaceAll(symbol, ",", " ")
		}
		decimals, _ = erc20Instance.Decimals(nil)
	}
	return
}

func getTokenERC721FromChain(client *ethclient.Client, tokenAddr common.Address) (name string, symbol string) {

	// Defaults shine through if any errors
	name = "Unknown"
	symbol = "UNKNOWN"

	// Create a new ERC721 instance using the token address and client
	erc20Instance, err := erc721.NewErc721(tokenAddr, client)
	if err == nil {
		name, err = erc20Instance.Name(nil)
		if err == nil {
			name = strings.ReplaceAll(name, ",", " ")
		}
		symbol, err = erc20Instance.Symbol(nil)
		if err == nil {
			symbol = strings.ReplaceAll(symbol, ",", " ")
		}
	}
	return
}
