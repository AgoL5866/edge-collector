package rpc

import (
	"bytes"
	"strconv"
	"time"
)

type CheckPointResp struct {
	Epoch                      string `json:"epoch"`
	SequenceNumber             string `json:"sequenceNumber"`
	Digest                     string `json:"digest"`
	NetworkTotalTransactions   string `json:"networkTotalTransactions"`
	PreviousDigest             string `json:"previousDigest"`
	EpochRollingGasCostSummary struct {
		ComputationCost         string `json:"computationCost"`
		StorageCost             string `json:"storageCost"`
		StorageRebate           string `json:"storageRebate"`
		NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
	} `json:"epochRollingGasCostSummary"`
	TimestampMs           string        `json:"timestampMs"`
	Transactions          []string      `json:"transactions"`
	CheckpointCommitments []interface{} `json:"checkpointCommitments"`
	ValidatorSignature    string        `json:"validatorSignature"`
}

type SuiApi struct {
	*BaseApi
}

func (api *SuiApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var number string
	var body = bytes.NewBufferString(`{"jsonrpc":"2.0", "id":1, "method":"sui_getLatestCheckpointSequenceNumber"}`)
	if err = api.BaseApi.request(&number, "", body); err != nil {
		h.Reason = err.Error()
		return
	}
	checkpoint, err := strconv.Atoi(number)
	if err != nil {
		h.Reason = err.Error()
		return err
	}
	var b = CheckPointResp{}
	body.Reset()
	body.WriteString(`{"jsonrpc":"2.0", "id":2, "method":"sui_getCheckpoint", "params":["` + number + `"]}`)
	if err = api.BaseApi.request(&b, "", body); err != nil {
		h.Reason = err.Error()
		return err
	}
	tsms, _ := strconv.Atoi(b.TimestampMs)
	h.Timestamp = time.Unix(int64(tsms)/1e3, int64(tsms)%1e3)
	h.Number = int64(checkpoint)
	h.Hash = b.Digest
	h.ParentHash = b.PreviousDigest
	return nil
}
