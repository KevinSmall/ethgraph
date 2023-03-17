package tokens

import "testing"

const embeddedTokenPresent = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
const embeddedTokenGarbage = "0xThisIsNotAValidToken"

func TestLoadingEmbeddedTokens(t *testing.T) {

	// Load Eth embedded tokens
	loadTokensEmbedded("1")
	if len(tokenMap) > 0 {
		t.Logf("Loaded some embedded tokens")
	} else {
		t.Errorf("Failed to load any embedded tokens")
	}

	// Inspect an actual token
	token, exists := GetTokenMasterData(embeddedTokenPresent)
	if !exists {
		t.Errorf("Failed to read an expected embedded token %s", embeddedTokenPresent)
	} else {
		if token.Symbol == "WETH" {
			t.Logf("Successfully read an embedded token")
		} else {
			t.Errorf("Read an embedded token %s but its Symbol was wrong", embeddedTokenPresent)
		}
	}

	// Inspect a garbage token
	token, exists = GetTokenMasterData(embeddedTokenGarbage)
	if exists {
		t.Errorf("Should not be able to read a non-existing token %s", embeddedTokenGarbage)
	} else {
		if token.Symbol == "UNKNOWN" {
			t.Logf("Read correct default master data for non-existing token")
		} else {
			t.Errorf("Read a non-existing token %s but its default Symbol was wrong", embeddedTokenGarbage)
		}
	}
}
