package rpc

import (
	"encoding/json"
	"strconv"
	"time"
)

type CosmosHeader struct {
	Version struct {
		Block string `json:"block"`
		App   string `json:"app"`
	} `json:"version"`
	ChainId            string          `json:"chain_id"`
	Height             string          `json:"height"`
	Time               time.Time       `json:"time"`
	LastBlockId        json.RawMessage `json:"last_block_id"`
	LastCommitHash     string          `json:"last_commit_hash"`
	DataHash           string          `json:"data_hash"`
	ValidatorsHash     string          `json:"validators_hash"`
	NextValidatorsHash string          `json:"next_validators_hash"`
	ConsensusHash      string          `json:"consensus_hash"`
	AppHash            string          `json:"app_hash"`
	LastResultsHash    string          `json:"last_results_hash"`
	EvidenceHash       string          `json:"evidence_hash"`
	ProposerAddress    string          `json:"proposer_address"`
}

type CosmosBlock struct {
	BlockId json.RawMessage `json:"block_id"`
	Block   struct {
		Header     CosmosHeader    `json:"header"`
		Data       json.RawMessage `json:"-"`
		Evidence   json.RawMessage `json:"-"`
		LastCommit json.RawMessage `json:"-"`
	} `json:"block"`
	SdkBlock json.RawMessage `json:"-"`
}

type CosmosApi struct {
	*BaseApi
}

func (api *CosmosApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var b = &CosmosBlock{}
	err = api.BaseApi.request(b, "/cosmos/base/tendermint/v1beta1/blocks/latest", nil)
	if err != nil {
		h.Reason = err.Error()
		return
	}

	header := b.Block.Header
	bh, err := strconv.Atoi(header.Height)
	if err != nil {
		h.Reason = err.Error()
		return
	}
	microSec := header.Time.UnixMicro()
	h.Number = int64(bh)
	h.Timestamp = time.UnixMicro(microSec)
	h.Hash = header.DataHash
	return
}
