package chain

import (
	"context"
	"embed"
	"encoding/json"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum/ethclient"
	"strings"
)

//go:embed chain_ids.json
var f embed.FS

type EvmClient struct {
	Name              string
	ChainId           string
	LatestBlockNumber uint64
	Url               string
	Client            *ethclient.Client
}

var chainIdMap = make(map[string]string)

const chainIdsFile = "chain_ids.json"

func init() {
	// Load chainIds and lookup chain name
	// read the JSON file into a byte slice
	fileContents, err := f.ReadFile(chainIdsFile)
	if err != nil {
		logr.Error.Panicf("Error when opening chain ids file %s: %s\n", chainIdsFile, err)
	}

	// unmarshal the JSON into the map
	err = json.Unmarshal(fileContents, &chainIdMap)
	if err != nil {
		logr.Error.Panicf("Error when reading chain ids file %s: %s\n", chainIdsFile, err)
	}
}

func getChainName(chainId string) string {
	chainName, exists := chainIdMap[chainId]
	if !exists {
		chainName = "Unknown"
	}
	chainName = strings.ReplaceAll(chainName, " ", "-")
	return chainName
}

// CreateEvmClient gets an EVM client and name. The name comes from reading embedded file copied
// from this JSON https://github.com/DefiLlama/chainlist/blob/main/constants/chainIds.json
func CreateEvmClient(url string) (EvmClient, error) {

	// connect to the EVM client
	client, err := ethclient.Dial(url)
	if err != nil {
		return EvmClient{}, err
	}

	// Get chainId
	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		logr.Error.Panicln("Error when reading chain to get chainId: ", err)
	}
	chainName := getChainName(chainId.String())

	// Get the latest block
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		logr.Error.Panicln("Error when reading chain to get latest block: ", err)
	}

	return EvmClient{
		Name:              chainName,
		ChainId:           chainId.String(),
		LatestBlockNumber: blockNumber,
		Url:               url,
		Client:            client,
	}, nil
}
