package rpc

import (
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type BlockHeaderInfo struct {
	ReqTime    time.Time
	Duration   int64  // ms
	Region     string // Singapore, Tokyo, Mumbai, Paris, Frankfurt, San Francisco
	Chain      string
	NodeAlias  string
	Number     int64
	Timestamp  time.Time
	Delay      int64 // ms, block delay time when use websocket subscribe
	Hash       string
	ParentHash string
	Reason     string //if have error, it maybe include error info
}

type Block struct {
	Number     rpc.BlockNumber `json:"number"`
	Hash       string          `json:"hash"`
	ParentHash string          `json:"parentHash"`
	Timestamp  string          `json:"timestamp"`
}

type HeaderNumberSorter []*BlockHeaderInfo

var _ sort.Interface = (*HeaderNumberSorter)(nil)

func (s HeaderNumberSorter) Len() int {
	return len(s)
}

func (s HeaderNumberSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s HeaderNumberSorter) Less(i, j int) bool {
	return s[i].Number < s[j].Number
}
