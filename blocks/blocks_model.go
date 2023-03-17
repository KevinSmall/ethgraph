package blocks

import "time"

type BlockKey struct {
	BlockNumber uint64
}

type BlockMapValue struct {
	BlockTimestamp   time.Time
	TransactionCount uint
}

// BlockDataFromSource is the data received from source (the chain)
type BlockDataFromSource struct {
	BlockNumber      uint64
	BlockTimestamp   time.Time
	TransactionCount uint
}

type BlockMap map[BlockKey]BlockMapValue
