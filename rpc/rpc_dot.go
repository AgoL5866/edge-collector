package rpc

import (
	"bytes"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

// curl -X POST -H "Content-Type: application/json" --ssl-no-revoke https://full-dot.w3node.com/7ed7b7b1882003f7f294d3ef9ca97d54473f64ccc8e6865527061115109983a1/api
// -d '{"jsonrpc":"2.0", "id":1, "method":"chain_getHeader"}' |jq
//
//	{
//	 "jsonrpc": "2.0",
//	 "result": {
//	   "parentHash": "0x12f414bf86fd4c0d2fb37a82d915b5ce4f8ff14084d5a3f6e3987dd9ac5b73da",
//	   "number": "0xeba7a6",
//	   "stateRoot": "0x07390138aea035d82df0d265b16f2203a2d5a3336720c34d922a03e19b857413",
//	   "extrinsicsRoot": "0x858a719fc4ebdd69d5ccbd67ffe649abe54c94e4c7429370980870a3b69dd243",
//	   "digest": {
//	     "logs": [
//	       "0x0642414245b501032201000074abb91000000000ba3e334db5efd09d566c3b01df6d6af91a373b093e4a5fae1cfa0c1ec0d9f746b62ba287faef4e7c54cbdf5d1b9994764f574bae0d4eb8296e8c1a1e337c58081c7855fdcf7f4b389868e3dc3592607232a6f7b7762583d1d9796d771797d10f",
//	       "0x05424142450101040e0a735cb12ab474f6531ce7a46d2c04c20bb7cf1a493cb0046f5126acf810b2186e4b4c5b64bfa55b72e72a6228c4331e5d2a4988523387424b729b8cfe86"
//	     ]
//	   }
//	 },
//	 "id": 1
//	}

type DotHeader struct {
	ParentHash     string          `json:"parentHash"`
	Number         rpc.BlockNumber `json:"number"`
	StateRoot      string          `json:"stateRoot"`
	ExtrinsicsRoot string          `json:"extrinsicsRoot"`
	Digest         any             `json:"-"`
}

type DotApi struct {
	*BaseApi
}

// BlockHeader get header
func (api *DotApi) BlockHeader(h *BlockHeaderInfo) (err error) {
	if h.ReqTime.IsZero() {
		h.ReqTime = time.Now().Truncate(time.Second).UTC()
	}
	var b = &DotHeader{}
	var body = bytes.NewBufferString(`{"jsonrpc":"2.0", "id":1, "method":"chain_getHeader"}`)
	if err = api.BaseApi.request(b, "", body); err != nil {
		h.Reason = err.Error()
		return
	}
	h.Number = b.Number.Int64()
	h.Hash = b.StateRoot
	h.ParentHash = b.ParentHash
	return nil
}
