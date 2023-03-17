package tokens

import (
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum/common"
	"sort"
)

// tokenMap is package global holding all token master data from any source, could
// be from embedded data, or local cache or anywhere else.
var tokenMap = make(map[string]TokenDataFromSource)

// Init is called to load just the master data for the desired chainId
func Init(chainId string) {
	if len(tokenMap) == 0 {

		// Embedded tokens
		loadTokensEmbedded(chainId)

		// Locally cached tokens
		loadTokensCached(chainId)
	}
}

func GetGlobalTokensLoadedCount() uint {
	return uint(len(tokenMap))
}

// GetTokenMasterData returns master data for token, or correct defaults if no
// master data found. This function returns valid data regardless of the token
// type (ERC20, ERC721, ERC1155).
func GetTokenMasterData(tokenAddr string) (tokenData TokenData, exists bool) {

	// Defaults shine through if missing
	tokenData = TokenData{
		Name:     "Unknown",
		Symbol:   "UNKNOWN",
		Decimals: 0,
	}
	exists = false

	// lookup our tokenMap keyed on address, to see what the token Symbol is
	tokenMapEntry, mapExists := tokenMap[tokenAddr]
	if mapExists {
		tokenData.Name = tokenMapEntry.Name
		tokenData.Symbol = tokenMapEntry.Symbol
		tokenData.Decimals = tokenMapEntry.Decimals
		exists = true
	}
	return
}

// RemoveUnknownTokens takes a map of addresses and removes those that are not known in the global tokenMap
func RemoveUnknownTokens(addressesMap map[common.Address]*AddressMapValue) map[common.Address]*AddressMapValue {
	for address, _ := range addressesMap {
		_, mapExists := tokenMap[address.Hex()]
		if mapExists {
			delete(addressesMap, address)
		}
	}
	return addressesMap
}

func PrintContents(addressesMap map[common.Address]*AddressMapValue) {
	for address, addressMapValue := range addressesMap {
		logr.Info.Printf("    %s, %s, %v\n", address.Hex(), addressMapValue.TransferType, addressMapValue.Count)
	}
}

func MergeTokensIntoGlobalTokenMap(newTokens map[string]TokenDataFromSource) {
	for a, t := range newTokens {
		tokenMap[a] = t
	}
}

// uniqueAddressesMap := tokens.(uniqueAddressesMap)

func BuildUniqueTokenAddressesFromEvents(allEvents []*chain.TransferEvent) map[common.Address]*AddressMapValue {
	// Create a map to store unique addresses
	uniqueAddressesMap := make(map[common.Address]*AddressMapValue)

	// Iterate through the events and get unique emitter addresses (== token addresses)
	for _, event := range allEvents {
		_, exists := uniqueAddressesMap[event.LogEmitterAddress]
		if exists {
			// We've seen this address before
			uniqueAddressesMap[event.LogEmitterAddress].Count++
		} else {
			// New address not seen before
			uniqueAddressesMap[event.LogEmitterAddress] = &AddressMapValue{
				TransferType: event.TransferType,
				Count:        1,
			}
		}
	}
	return uniqueAddressesMap
}

func SortAddressesByCount(uniqueAddresses map[common.Address]*AddressMapValue) []AddressValue {
	// convert map to slice of AddressValue
	addressValues := make([]AddressValue, 0, len(uniqueAddresses))
	for address, addressMapValue := range uniqueAddresses {
		addressValues = append(addressValues, AddressValue{
			Address:      address,
			TransferType: addressMapValue.TransferType,
			Count:        addressMapValue.Count,
		})
	}
	// sort slice by count descending
	sort.Slice(addressValues, func(i, j int) bool {
		return addressValues[i].Count > addressValues[j].Count
	})
	return addressValues
}
