package rpc

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

type EthApi struct {
	url  string
	host string
	rpc  *rpc.Client
}

func (api *EthApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	var b Block
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = api.rpc.CallContext(ctx, &b, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = b.Number.Int64()
	ts, ok := big.NewInt(0).SetString(b.Timestamp, 0)
	if ok {
		h.Timestamp = time.Unix(ts.Int64(), 0).UTC()
	}
	h.Hash = b.Hash
	h.ParentHash = b.ParentHash
	return
}

func NewEthApi(url string, host string, maxConn int64, timeout time.Duration) *EthApi {
	httpCli := newHttpClient(int(maxConn), timeout)
	cli, err := rpc.DialHTTPWithClient(url, httpCli)
	if err != nil {
		log.Errorf("new EthApi error: %s", err)
		return nil
	}
	if host != "" {
		cli.SetHeader("X-Forwarded-Host", host)
	}
	api := &EthApi{
		url: url, host: host, rpc: cli,
	}
	return api
}
