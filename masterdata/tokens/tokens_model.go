package tokens

import "github.com/ethereum/go-ethereum/common"

type AddressMapValue struct {
	TransferType string
	Count        uint
}

type AddressValue struct {
	Address      common.Address
	TransferType string
	Count        uint
}

// TokenData is what is returned by GetTokenData calls to the tokens package
type TokenData struct {
	// Name is the token free text name
	Name string

	// Symbol is the short token symbol or ticker
	Symbol string

	// Decimals are how many decimal places the token has (ERC20 only, 0 for ERC721)
	Decimals int
}

// TokenDataFromSource is token master data from source: CSV file, cache, or chain
type TokenDataFromSource struct {
	// ChainId uniquely identifies an EVM-compatible chain see https://chainlist.org/
	ChainId string

	// Name is the token free text name
	Name string

	// Symbol is the short token symbol or ticker
	Symbol string

	// Decimals are how many decimal places the token has (ERC20 only, 0 for ERC721)
	Decimals int

	// TokenAddress is the contract address of the token
	TokenAddress string
}
