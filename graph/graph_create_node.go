// Package graph contains functions to create and write graphML
// It uses https://github.com/yaricom/goGraphML see readme there
package graph

import (
	"fmt"
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/conv"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/KevinSmall/ethgraph/masterdata/addresses"
	"github.com/KevinSmall/ethgraph/masterdata/tokens"
	"github.com/yaricom/goGraphML/graphml"
)

func addAddressNodesToGraph(events []*chain.TransferEvent, gr *graphml.Graph) (
	uniqueAddressesAsNodesMap map[string]*graphml.Node) {

	// Populate map to store unique addresses == nodes
	uniqueAddressesAsNodesMap = make(map[string]*graphml.Node)
	for _, event := range events {
		// Node using event log FROM address
		_, exists := uniqueAddressesAsNodesMap[event.LogAddressFrom.Hex()]
		if !exists {
			// we've not seen node before
			n := createAddressNodeAndAddToGraph(true, event, gr)
			uniqueAddressesAsNodesMap[event.LogAddressFrom.Hex()] = n
		}

		// Node using event log TO address
		_, exists = uniqueAddressesAsNodesMap[event.LogAddressTo.Hex()]
		if !exists {
			// we've not seen node before
			n := createAddressNodeAndAddToGraph(false, event, gr)
			uniqueAddressesAsNodesMap[event.LogAddressTo.Hex()] = n
		}
	}
	return uniqueAddressesAsNodesMap
}

func createAddressNodeAndAddToGraph(isFromAddress bool, event *chain.TransferEvent, gr *graphml.Graph) *graphml.Node {
	attributes := make(map[string]interface{})
	address := event.LogAddressTo.Hex()
	if isFromAddress {
		address = event.LogAddressFrom.Hex()
	}
	attributes["address"] = address
	label := conv.PrettyShortenAddress(address)
	addressData, addressMasterDataExists := addresses.GetAddressMasterData(address)
	if addressMasterDataExists {
		label = addressData.Description
	}
	attributes["description"] = addressData.Description
	attributes["nodeType"] = 1
	timestamp := formatTimestamp(event.LogAddressToFirstSeen)
	timeIndex := event.LogAddressToFirstSeenIndex
	if isFromAddress {
		timestamp = formatTimestamp(event.LogAddressFromFirstSeen)
		timeIndex = event.LogAddressFromFirstSeenIndex
	}
	attributes["timestampEstimate"] = timestamp
	attributes["appearanceIndex"] = int(timeIndex)

	n, err := gr.AddNode(attributes, label)
	if err != nil {
		logr.Error.Panicln(err)
	}
	return n
}

func addMovementNodesToGraph(events []*chain.TransferEvent, gr *graphml.Graph) (
	uniqueMovementsAsNodesMap map[mvtNodeKey]*graphml.Node) {

	// Populate map to store unique movements == nodes
	uniqueMovementsAsNodesMap = make(map[mvtNodeKey]*graphml.Node)
	for _, event := range events {

		// Create new node(s) for each transfer event
		attributes := make(map[string]interface{})
		tokenData, tokenMasterDataExists := tokens.GetTokenMasterData(event.LogEmitterAddress.Hex())
		attributes["symbol"] = tokenData.Symbol
		tokenValue := float64(0)
		if tokenMasterDataExists {
			tokenValue = conv.SafeScaleTokenValue(&event.LogTokenValue, tokenData.Decimals)
		}
		attributes["value"] = tokenValue
		attributes["nodeType"] = 0
		attributes["nftId"] = event.LogNftId
		attributes["transferType"] = event.TransferType
		attributes["txHash"] = event.TxHash.Hex()
		attributes["txIndex"] = int(event.TxIndex)
		attributes["timestampEstimate"] = formatTimestamp(event.TransactionTimestampEstimate)
		attributes["appearanceIndex"] = int(event.TransactionTimestampEstimateIndex)

		// Add node to graph
		label := ""
		timeStamp := formatTimestampShort(event.TransactionTimestampEstimate)
		switch event.TransferType {
		case chain.ERC20:
			label = fmt.Sprintf("%v %s (%s)", tokenValue, tokenData.Symbol, timeStamp)
		case chain.ERC721:
			label = fmt.Sprintf("NFT %s %s (%s)", event.LogNftId, tokenData.Symbol, timeStamp)
		case chain.ERC1155_SINGLE, chain.ERC1155_BATCH:
			label = fmt.Sprintf("%v of NFT %s %s (%s)", tokenValue, event.LogNftId, tokenData.Symbol, timeStamp)
		default:
			logr.Warning.Printf("Unknown transfer type %s.", event.TransferType)
			continue
		}

		n, err := gr.AddNode(attributes, label)
		if err != nil {
			logr.Error.Panicln(err)
		}

		// Store graph node for later reference during edge processing
		movementNodeKey := mvtNodeKey{
			edgeFrom: event.LogAddressFrom.Hex(),
			edgeTo:   event.LogAddressTo.Hex(),
			txHash:   event.TxHash.Hex(),
			logIndex: event.LogIndex,
			nftId:    event.LogNftId,
		}
		uniqueMovementsAsNodesMap[movementNodeKey] = n
	}
	return uniqueMovementsAsNodesMap
}
