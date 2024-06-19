package rpc

import "time"

type EgldStats struct {
	Shards         int64 `json:"shards"`
	Blocks         int64 `json:"blocks"`
	Accounts       int64 `json:"accounts"`
	Transactions   int64 `json:"transactions"`
	ScResults      int64 `json:"scResults"`
	RefreshRate    int64 `json:"refreshRate"`
	Epoch          int64 `json:"epoch"`
	RoundsPassed   int64 `json:"roundsPassed"`
	RoundsPerEpoch int64 `json:"roundsPerEpoch"`
}

// EgldApi == MultiverseX
// only for egld-api
type EgldApi struct {
	*BaseApi
}

func (api *EgldApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var b = &EgldStats{}
	if err = api.BaseApi.request(b, "/stats", nil); err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = b.Blocks
	return nil
}
