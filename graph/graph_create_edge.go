package graph

import (
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/yaricom/goGraphML/graphml"
)

func addEdgesToGraph(
	uniqueAddressesToNodeMap map[string]*graphml.Node,
	uniqueMovementsToNodeMap map[mvtNodeKey]*graphml.Node,
	events []*chain.TransferEvent,
	g *graphml.Graph) {

	for _, event := range events {

		// Locate required nodes, per transfer there are 3: address-from, address-to with
		// the transfer event node in the middle
		nodeFrom, exists := uniqueAddressesToNodeMap[event.LogAddressFrom.Hex()]
		if !exists {
			logr.Error.Panicln("NodeFrom missing in internal map")
		}
		nodeTo, exists := uniqueAddressesToNodeMap[event.LogAddressTo.Hex()]
		if !exists {
			logr.Error.Panicln("NodeTo missing in internal map")
		}
		nodeForEventKey := mvtNodeKey{
			edgeFrom: event.LogAddressFrom.Hex(),
			edgeTo:   event.LogAddressTo.Hex(),
			txHash:   event.TxHash.Hex(),
			logIndex: event.LogIndex,
			nftId:    event.LogNftId,
		}
		nodeForEvent, exists := uniqueMovementsToNodeMap[nodeForEventKey]
		if !exists {
			logr.Error.Panicln("nodeForTransferEventMovement missing in internal map ", nodeForEventKey)
		}

		// Edge creation
		attributes := make(map[string]interface{})
		attributes["transferType"] = event.TransferType
		nodeAttrs, err := nodeForEvent.GetAttributes()
		if err != nil {
			logr.Error.Panicln(err)
		}
		attributes["symbol"] = nodeAttrs["symbol"]
		attributes["timestampEstimate"] = formatTimestamp(event.TransactionTimestampEstimate)
		attributes["appearanceIndex"] = int(event.TransactionTimestampEstimateIndex)

		label := ""
		// Create edge from "transfer log from-address" to "transfer event middle node"
		_, err = g.AddEdge(nodeFrom, nodeForEvent, attributes, graphml.EdgeDirectionDirected, label)
		if err != nil {
			if event.TransferType == chain.ERC1155_BATCH {
				// Occasional bad data in ERC1155 batches, just ignore
				continue
			} else {
				logr.Error.Panicln(err)
			}
		}
		// Create edge "transfer event middle node" to "transfer log to-address"
		_, err = g.AddEdge(nodeForEvent, nodeTo, attributes, graphml.EdgeDirectionDirected, label)
		if err != nil {
			if event.TransferType == chain.ERC1155_BATCH {
				// Occasional bad data in ERC1155 batches, just ignore
				continue
			} else {
				logr.Error.Panicln(err)
			}
		}
	}
}
