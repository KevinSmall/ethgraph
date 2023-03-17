package tokens

import (
	"embed"
	"encoding/csv"
	"fmt"
	"github.com/KevinSmall/ethgraph/conv"
	"github.com/KevinSmall/ethgraph/logr"
	"io"
	"os"
	"strconv"
	"strings"
)

//go:embed data/tokens_1.csv
//go:embed data/tokens_56.csv
//go:embed data/tokens_43114.csv
var f embed.FS

// loadTokensEmbedded loads embedded CSV files of token addresses and master data
// like name, symbol. This function returns a map, keyed on token address ( == the address that
// emitted the transfer event) with a value containing a struct of token master data.
func loadTokensEmbedded(chainId string) {

	// Validate against embedded files
	if chainId != "1" && chainId != "56" && chainId != "43114" {
		logr.Trace.Printf("Not loading an embedded token file because no file available for chainId %s", chainId)
		return
	}

	tokenFiles := []string{fmt.Sprintf("data/tokens_%s.csv", chainId)}

	for _, filename := range tokenFiles {
		// Read embedded CSV file
		fileContents, err := f.ReadFile(filename)
		if err != nil {
			logr.Error.Panicln("Error when opening file: ", filename, " ", err)
		}
		addFileContentsToGlobalTokenMap(chainId, filename, fileContents)
	}
}

func addFileContentsToGlobalTokenMap(chainId string, filename string, fileContents []byte) {
	// Parse CSV data
	reader := csv.NewReader(strings.NewReader(string(fileContents)))

	// Skip first row of headers
	_, err := reader.Read()
	if err != nil {
		logr.Error.Panicf("Error when processing token file %s %s. Try deleting it and rerunning.", filename, err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logr.Warning.Printf("Error when processing file %s contents: %s", filename, err)
			break
		}
		if len(record) != 5 {
			continue
		}

		// Convert record to TokenData struct
		tokenData := TokenDataFromSource{
			ChainId:      record[0],
			Name:         record[1],
			Symbol:       record[2],
			Decimals:     conv.SafeStringToInt(record[3]),
			TokenAddress: record[4],
		}
		// Only tokens for chosen chain, and only those with valid length for hex address
		if tokenData.ChainId == chainId && len(tokenData.TokenAddress) == 42 {
			tokenMap[tokenData.TokenAddress] = tokenData
		}
	}
}

// loadTokensCached loads the cached CSV file of token addresses and master data
// like name, symbol. This function updates the global map.
// The data is filtered to only be for the passed chainId.
func loadTokensCached(chainId string) {

	// Read cache CSV file
	filename := getTokenCacheFilename(chainId)

	file, err := os.Open(filename)
	if err != nil {
		logr.Trace.Printf("Failed to open file: %s. It probably doesn't exist yet: %s.", filename, err)
		return
	}
	defer file.Close()

	fileContents, err := io.ReadAll(file)
	if err != nil {
		logr.Warning.Printf("Local CSV cache not read for chainId %s. Error: %s.", chainId, err)
		return
	}
	addFileContentsToGlobalTokenMap(chainId, filename, fileContents)
}

func getTokenCacheFilename(chainId string) (filename string) {
	return fmt.Sprintf(".tokens_%s_cache.csv", chainId)
}

// WriteGlobalTokenMapToCache takes the whole of the tokenMap master data and dumps it all out
// to the local cache .csv file. If any troubles, it just does nothing.
func WriteGlobalTokenMapToCache(chainId string) {
	var tokenInfo [][]string
	for _, t := range tokenMap {
		tokenInfo = append(tokenInfo, []string{
			t.ChainId,
			t.Name,
			t.Symbol,
			strconv.Itoa(t.Decimals),
			t.TokenAddress})
	}
	filename := getTokenCacheFilename(chainId)
	WriteTokenInfoToFile(filename, tokenInfo)
}

func DeleteTokenCache(chainId string) {
	filename := getTokenCacheFilename(chainId)
	err := os.Remove(filename)
	if err != nil {
		// file doesn't exist, which is fine
		return
	}
	// Check if the file has been deleted
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		logr.Info.Printf("Token cache file %s has been deleted\n", filename)
	} else if err != nil {
		logr.Info.Printf("Token cache file %s not deleted: %s\n", filename, err)
	} else {
		logr.Info.Printf("Token cache file %s still exists\n", filename)
	}
}

func WriteTokenInfoToFile(filename string, tokenInfo [][]string) {
	// Create a new file for writing
	file, err := os.Create(filename)
	if err != nil {
		logr.Error.Panicln(err)
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	err = writer.Write([]string{"ChainId", "Name", "Symbol", "Decimals", "TokenAddress"})
	if err != nil {
		logr.Error.Panicln(err)
	}

	// Write the token information rows
	for _, row := range tokenInfo {
		err = writer.Write(row)
		if err != nil {
			logr.Error.Panicln(err)
		}
	}
}
