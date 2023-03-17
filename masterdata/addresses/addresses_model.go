package addresses

type AddressData struct {
	// Description is free text description about this address, e.g. Binance, Gemini
	Description string
}

type addressPopularData struct {
	// ChainId uniquely identifies an EVM-compatible chain see https://chainlist.org/
	ChainId string

	// Description is free text description about this address, e.g. Binance, Gemini
	Description string

	// Address is the EoA or contract address
	Address string
}
