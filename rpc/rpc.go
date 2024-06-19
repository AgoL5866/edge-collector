package rpc

import (
	"fmt"
	"sync"
	"time"
)

const (
	rpcTypeAda      = "ada-rpc"
	rpcTypeAptos    = "aptos-rpc"
	rpcTypeDot      = "dot-rpc"
	rpcTypeEth      = "eth-rpc"
	rpcTypeEgld     = "egld-rpc"
	rpcTypeSol      = "sol-rpc"
	rpcTypeStarknet = "starknet-rpc"
	rpcTypeTon      = "ton-rpc"
	rpcTypeSui      = "sui-rpc"
	rpcTypeCosmos   = "cosmos-rpc"
)

type ChainRpc interface {
	BlockHeader(h *BlockHeaderInfo) (e error)
}

var (
	gProviderLock sync.Mutex
	gProviders    = map[string]ChainRpc{} //url -> rpc
)

func GetRpcProvider(t, url string) ChainRpc {
	var rpc ChainRpc
	gProviderLock.Lock()
	rpc = gProviders[url]
	gProviderLock.Unlock()
	if rpc != nil {
		return rpc
	}

	base := NewBaseApi(url, nil, 8, 30*time.Second)
	switch t {
	case rpcTypeAda:
		base.noExtract = true
		rpc = &AdaApi{base}
	case rpcTypeAptos:
		base.noExtract = true
		rpc = &AptosApi{base}
	case rpcTypeEgld:
		base.noExtract = true
		rpc = &EgldApi{base}
	case rpcTypeStarknet:
		rpc = &StarknetApi{base}
	case rpcTypeDot:
		rpc = &DotApi{base}
	case rpcTypeSol:
		rpc = &SolApi{base}
	case rpcTypeTon:
		rpc = &TonApi{base}
	case rpcTypeSui:
		rpc = &SuiApi{BaseApi: base}
	case rpcTypeCosmos:
		base.noExtract = true
		rpc = &CosmosApi{BaseApi: base}
	case rpcTypeEth:
		fallthrough
	default:
		rpc = NewEthApi(url, "", 8, 30*time.Second)
	}
	if rpc == nil {
		panic(fmt.Sprintf("can not init rpc provider:%s, with url:%s", t, url))
	}
	gProviderLock.Lock()
	gProviders[url] = rpc
	gProviderLock.Unlock()
	return rpc
}
