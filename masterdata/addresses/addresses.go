package addresses

var addressMap = make(map[string]addressPopularData)

// Init is called manually to load just the master data for the desired chainId
// (standard init() didn't know the chainId yet)
func Init(chainId string) {
	if len(addressMap) == 0 {

		// Embedded address data
		loadAddressesEmbedded(chainId)

		// Local file of address data
		loadAddressesCached(chainId)
	}
}

func GetAddressesLoadedCount() uint {
	return uint(len(addressMap))
}

// GetAddressMasterData returns master data for address, or correct defaults if no
// master data found.
func GetAddressMasterData(addr string) (addressData AddressData, exists bool) {

	// Defaults
	addressData = AddressData{
		Description: "",
	}
	exists = false

	// lookup our addressMap keyed on address, to get its master data
	addressMapEntry, mapExists := addressMap[addr]
	if mapExists {
		addressData.Description = addressMapEntry.Description
		exists = true
	}
	return
}
