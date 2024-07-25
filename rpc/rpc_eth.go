package rpc

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthApi struct {
	url    string
	host   string
	client *ethclient.Client
}

func (api *EthApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	var b Block
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//如果是支持订阅的话，使用订阅模式去获取最新区块头
	if api.client.Client().SupportsSubscriptions() {
		err = api.OnceSubscribeNewHeads(h)
		if err != nil {
			h.Reason = err.Error()
		}
	} else {
		err = api.client.Client().CallContext(ctx, &b, "eth_getBlockByNumber", "latest", false)
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
	}
	return
}

// OnceSubscribeNewHeads 订阅一次新区块头
func (api *EthApi) OnceSubscribeNewHeads(h *BlockHeaderInfo) (err error) {
	ch := make(chan *types.Header)
	sub, err := api.client.SubscribeNewHead(context.Background(), ch)
	if err != nil {
		return err
	}

	//todo goroutine loop
	select {
	case <-sub.Err():
		return err
	case head := <-ch:
		h.Number = head.Number.Int64()
		h.Hash = head.Hash().String()
		h.ParentHash = head.ParentHash.String()
		h.Timestamp = time.Unix(int64(head.Time), 0).UTC()
	}
	return nil
}

func NewEthApi(url string, host string, maxConn int64, timeout time.Duration) *EthApi {
	httpCli := newHttpClient(int(maxConn), timeout)
	rpcCli, err := rpc.DialHTTPWithClient(url, httpCli)
	if err != nil {
		log.Println("new EthApi error: ", err)
		return nil
	}
	if host != "" {
		rpcCli.SetHeader("X-Forwarded-Host", host)
	}
	api := &EthApi{
		url: url, host: host, client: ethclient.NewClient(rpcCli),
	}
	return api
}
