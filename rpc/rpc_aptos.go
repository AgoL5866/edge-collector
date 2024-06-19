package rpc

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type AptosHeight struct {
	ChainId             int    `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       string `json:"ledger_version"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	LedgerTimestamp     string `json:"ledger_timestamp"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   string `json:"oldest_block_height"`
	BlockHeight         string `json:"block_height"`
	GitHash             string `json:"git_hash"`
}

type AptosBlock struct {
	BlockHeight    string `json:"block_height"`
	BlockHash      string `json:"block_hash"`
	BlockTimestamp string `json:"block_timestamp"`
	FirstVersion   string `json:"first_version"`
	LastVersion    string `json:"last_version"`
	Transactions   any    `json:"-"`
}

type AptosApi struct {
	*BaseApi
}

func (api *AptosApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var height = &AptosHeight{}
	err = api.BaseApi.request(height, "/v1", nil)
	if err != nil {
		h.Reason = err.Error()
		return
	}

	if height.BlockHeight == "" {
		err = errors.New("blockHeight empty")
		h.Reason = err.Error()
		return
	}
	var b = &AptosBlock{}
	err = api.BaseApi.request(b, fmt.Sprintf("/v1/blocks/by_height/%s?with_transactions=false", height.BlockHeight), nil)
	if err != nil {
		h.Reason = err.Error()
		return
	}
	bh, err := strconv.Atoi(b.BlockHeight)
	if err != nil {
		h.Reason = err.Error()
		return
	}
	microSec, err := strconv.Atoi(b.BlockTimestamp)
	if err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = int64(bh)
	h.Timestamp = time.UnixMicro(int64(microSec))
	h.Hash = b.BlockHash
	return
}
