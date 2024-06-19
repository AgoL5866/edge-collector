package rpc

import (
	"bytes"
	"time"
)

type StarknetHeader struct {
	BlockHash   string `json:"block_hash"`
	BlockNumber int    `json:"block_number"`
}

type StarknetApi struct {
	*BaseApi
}

func (api *StarknetApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var b = &StarknetHeader{}
	var body = bytes.NewBufferString(`{"jsonrpc":"2.0", "id":1, "method":"starknet_blockHashAndNumber"}`)
	if err = api.BaseApi.request(b, "", body); err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = int64(b.BlockNumber)
	h.Hash = b.BlockHash
	return nil
}
