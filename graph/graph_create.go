// Package graph contains functions to create and write graphML
// It uses https://github.com/yaricom/goGraphML see readme there
package graph

import (
	"github.com/KevinSmall/ethgraph/chain"
	"github.com/KevinSmall/ethgraph/logr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/yaricom/goGraphML/graphml"
	"sort"
)

var zeroAddress = common.Address{}

const topNMissingTokens = 10
const topNExistingTokens = 10

// CreateGraph creates a GraphML graph from a slice of TransferEvents. CreateGraph does
// not care what selections were used to produce the slice of TransferEvents, for example
// a selection by address, or by block, or filtered by a single token, it does not know
// or care it just blindly converts events to a graph.
func CreateGraph(graphTitle string,
	events []*chain.TransferEvent) (
	graphMlRoot *graphml.GraphML,
	graphCreationResult CreationResult) {

	// for preparing test data
	//fmt.Print("var testData = []*chain.TransferEvent{")
	//for i, event := range events {
	//	if i > 0 {
	//		fmt.Print(", ")
	//	}
	//	fmt.Printf("{%#v}", *event)
	//}
	//fmt.Println("}")

	// Root
	graphMlRoot = graphml.NewGraphML(graphTitle)

	// Graph
	g, err := graphMlRoot.AddGraph(graphTitle, graphml.EdgeDirectionDirected, nil)
	if err != nil {
		logr.Error.Panicln(err)
	}

	// Movement from/to addresses become Address Graph Nodes, these nodes are always required, these have nodeType 1
	uniqueAddressesToNodeMap := addAddressNodesToGraph(events, g)
	var uniqueMovementsToNodeMap map[mvtNodeKey]*graphml.Node
	// Movement Events == More Nodes
	uniqueMovementsToNodeMap = addMovementNodesToGraph(events, g)

	// Movement Events == Edges as well, create edges AND add them to graph at the same time
	addEdgesToGraph(uniqueAddressesToNodeMap, uniqueMovementsToNodeMap, events, g)

	// Build results (not part of the graphML, this is for info)
	graphCreationResult = CreationResult{
		Nodes:  len(uniqueAddressesToNodeMap) + len(uniqueMovementsToNodeMap),
		Edges:  len(g.Edges),
		Events: len(events),
	}
	return
}

// getTopNTokenAddresses sorts tokensMap by its value descending, and retains the returnTopN
// rows, returns result as a slice.
func getTopNTokenAddresses(tokensMap map[string]int, returnTopN int) (topNTokens []TokenAddressCount) {

	// Convert map to slice of structs
	for address, count := range tokensMap {
		topNTokens = append(topNTokens, TokenAddressCount{Address: address, Count: count})
	}

	// Sort slice by count in descending order
	sort.SliceStable(topNTokens, func(i, j int) bool {
		return topNTokens[i].Count > topNTokens[j].Count
	})

	// Take top N tokens
	if len(topNTokens) > returnTopN {
		topNTokens = topNTokens[:returnTopN]
	}
	return topNTokens
}
