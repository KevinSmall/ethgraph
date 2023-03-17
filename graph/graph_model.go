package graph

import (
	"github.com/KevinSmall/ethgraph/logr"
	"time"
)

type CreationResult struct {
	Nodes  int
	Edges  int
	Events int
}

type TokenAddressCount struct {
	Address string
	Count   int
}

type mvtNodeKey struct {
	edgeFrom string
	edgeTo   string
	txHash   string
	logIndex uint
	nftId    string
}

func formatTimestamp(time time.Time) string {
	return time.Format("2006-01-02 15:04:05.999")
}

func formatTimestampShort(time time.Time) string {
	return time.Format("15:04:05")
}

func (creationResult CreationResult) PrintSummary() {

	logr.Info.Printf("Chain Transfer Events: %v\n", creationResult.Events)
	logr.Info.Printf("GraphML Nodes: %v\n", creationResult.Nodes)
	logr.Info.Printf("GraphML Edges: %v\n", creationResult.Edges)
}
