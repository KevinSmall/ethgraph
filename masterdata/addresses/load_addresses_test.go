package addresses

import "testing"

const embeddedAddressPresent = "0x71660c4005BA85c37ccec55d0C4493E66Fe775d3"
const embeddedAddressGarbage = "0xThisIsNotAValidAddress"

func TestLoadingEmbeddedAddresses(t *testing.T) {

	// Load Eth embedded addresses
	loadAddressesEmbedded("1")
	if len(addressMap) > 0 {
		t.Logf("Loaded some embedded addresses ok")
	} else {
		t.Errorf("Failed to load any embedded addresses")
	}

	// Inspect an actual address
	address, exists := GetAddressMasterData(embeddedAddressPresent)
	if !exists {
		t.Errorf("Failed to read expected embedded address %s", embeddedAddressPresent)
	} else {
		if address.Description == "Coinbase" {
			t.Logf("Successfully read an embedded address")
		} else {
			t.Errorf("Read an embedded address %s but its Description was wrong", embeddedAddressPresent)
		}
	}

	// Inspect a garbage address
	address, exists = GetAddressMasterData(embeddedAddressGarbage)
	if exists {
		t.Errorf("Should not be able to read a non-existing address %s", embeddedAddressGarbage)
	} else {
		if address.Description == "" {
			t.Logf("Read correct default master data for non-existing address")
		} else {
			t.Errorf("Read a non-existing address %s but its default Description was wrong", embeddedAddressGarbage)
		}
	}
}
