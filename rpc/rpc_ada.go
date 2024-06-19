package rpc

import (
	"bytes"
	"time"
)

type AdaHeader struct {
	Data struct {
		Cardano struct {
			Tip struct {
				EpochNo int    `json:"epochNo"`
				Number  int    `json:"number"`
				Hash    string `json:"hash"`
			} `json:"tip"`
		} `json:"cardano"`
	} `json:"data"`
}
type AdaApi struct {
	*BaseApi
}

func (api *AdaApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var b = &AdaHeader{}
	var body = bytes.NewBufferString(`{"query":"{cardano{tip{epochNo number hash}}}"}`)
	if err = api.BaseApi.request(b, "/graphql", body); err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = int64(b.Data.Cardano.Tip.Number)
	h.Hash = b.Data.Cardano.Tip.Hash
	return nil
}
