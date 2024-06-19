package rpc

import (
	"time"
)

type TonHeader struct {
	ConsensusBlock int64   `json:"consensus_block"`
	Timestamp      float64 `json:"timestamp"`
}

type TonApi struct {
	*BaseApi
}

func (api *TonApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var b = &TonHeader{}
	if err = api.BaseApi.request(b, "/getConsensusBlock", nil); err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = b.ConsensusBlock
	h.Timestamp = time.Unix(int64(b.Timestamp), 0)
	return nil
}
