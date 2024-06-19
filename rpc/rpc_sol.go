package rpc

import (
	"bytes"
	"fmt"
	"time"
)

type SolHeader struct {
	BlockHeight       int64  `json:"blockHeight"`
	BlockTime         int64  `json:"blockTime"`
	BlockHash         string `json:"blockhash"`
	ParentSlot        int64  `json:"parentSlot"`
	PreviousBlockhash string `json:"previousBlockhash"`
}

type SolApi struct {
	*BaseApi
}

func (api *SolApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var number int64
	var body = bytes.NewBufferString(`{"jsonrpc":"2.0", "id":1, "method":"getSlot"}`)
	if err = api.BaseApi.request(&number, "", body); err != nil {
		h.Reason = err.Error()
		return
	}
	var b SolHeader
	body = bytes.NewBufferString(fmt.Sprintf(`{"jsonrpc":"2.0", "id":1, "method":"getBlock","params":[%d, {"transactionDetails":"none","rewards":false }]}`, number))
	if err = api.BaseApi.request(&b, "", body); err != nil {
		h.Reason = err.Error()
		return
	}
	h.Timestamp = time.Unix(b.BlockTime, 0)
	h.Number = b.BlockHeight
	h.Hash = b.BlockHash
	h.ParentHash = b.PreviousBlockhash
	return nil
}
