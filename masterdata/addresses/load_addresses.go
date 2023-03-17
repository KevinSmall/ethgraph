package addresses

import (
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/KevinSmall/ethgraph/logr"
	"io"
	"os"
	"strings"
)

//go:embed data/addresses_1.csv
//go:embed data/addresses_56.csv
//go:embed data/addresses_43114.csv
var f embed.FS

// loadAddressesEmbedded loads embedded CSV files of addresses and master data
// like name. It is intended for non-tokens, and the data is informational data not available
// on chain, so exchange names and the like that can be shown for a node.
func loadAddressesEmbedded(chainId string) {

	// Validate against embedded files
	if chainId != "1" && chainId != "56" && chainId != "43114" {
		logr.Trace.Printf("Not loading an embedded address file because no file available for chainId %s", chainId)
		return
	}

	addressFiles := []string{fmt.Sprintf("data/addresses_%s.csv", chainId)}

	for _, filename := range addressFiles {
		// Read embedded CSV file
		fileContents, err := f.ReadFile(filename)
		if err != nil {
			logr.Error.Panicln("Error when opening file: ", filename, " ", err)
		}

		// Parse CSV data
		reader := csv.NewReader(strings.NewReader(string(fileContents)))
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				logr.Error.Panicln("Error when processing file: ", filename, " ", err)
			}
			// Convert record to addressPopularData struct
			addressData := addressPopularData{
				ChainId:     record[0],
				Description: record[1],
				Address:     record[2],
			}
			// Only tokens for chosen chain, and only those with valid length for hex address
			if addressData.ChainId == chainId && len(addressData.Address) == 42 {
				addressMap[addressData.Address] = addressData
			}
		}
	}
}

// loadAddressesCached loads local CSV file of address master data
func loadAddressesCached(chainId string) {

	// Read cache CSV file
	filename := getAddressCacheFilename(chainId)

	file, err := os.Open(filename)
	if err != nil {
		logr.Trace.Printf("Failed to open file: %s. It probably doesn't exist yet: %s.", filename, err)
		return
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		logr.Warning.Printf("Local CSV cache not read for chainId %s. Error: %s. Try deleting it and rerun.", chainId, err)
		return
	}

	// Parse CSV data
	reader := csv.NewReader(strings.NewReader(string(fileContents)))

	// Skip first row of headers
	_, err = reader.Read()
	if err != nil {
		logr.Error.Panicln("Error when processing file: ", filename, ". Try deleting it and rerun.", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logr.Error.Panicln("Error when processing file: ", filename, ". Try deleting it and rerun.", err)
		}
		if len(record) != 3 {
			continue
		}

		// Convert record to addressPopularData struct
		addressData := addressPopularData{
			ChainId:     record[0],
			Description: record[1],
			Address:     record[2],
		}
		// Only tokens for chosen chain, and only those with valid length for hex address
		if addressData.ChainId == chainId && len(addressData.Address) == 42 {
			addressMap[addressData.Address] = addressData
		}
	}
}

func getAddressCacheFilename(chainId string) (filename string) {
	return fmt.Sprintf(".addresses_%s_cache.csv", chainId)
}
