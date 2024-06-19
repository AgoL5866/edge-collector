package rpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var defaultTransport = &http.Transport{
	Proxy: nil,
	DialContext: func() func(ctx context.Context, network string, address string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 120 * time.Second,
			Resolver:  net.DefaultResolver,
		}
		return dialer.DialContext
	}(),
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          8,
	IdleConnTimeout:       120 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ResponseHeaderTimeout: 15 * time.Second,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
}

func newHttpClient(maxConn int, timeout time.Duration) *http.Client {
	tr := defaultTransport.Clone()
	cli := &http.Client{
		Transport: tr,
	}
	if timeout > 0 {
		cli.Timeout = timeout
	}
	if maxConn > 0 {
		tr.MaxIdleConns = maxConn
	}
	return cli
}

type BaseApi struct {
	url        string
	header     http.Header
	cli        *http.Client
	noExtract  bool //提取 {result:<data>} 部份
	headerLock sync.Mutex
}

func (api *BaseApi) setHeader(key, val string) {
	api.headerLock.Lock()
	api.header.Set(key, val)
	api.headerLock.Unlock()
}

func NewBaseApi(url string, header map[string]string, maxConn int64, timeout time.Duration) *BaseApi {
	httpCli := newHttpClient(int(maxConn), timeout)
	api := &BaseApi{
		url: url, header: http.Header{}, cli: httpCli,
	}
	api.headerLock.Lock()
	api.header.Add("accept", "application/json")
	api.header.Add("content-type", "application/json")
	for k, v := range header {
		api.header.Set(k, v)
	}
	api.headerLock.Unlock()
	return api
}

type response struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *jsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

// request
// rst : the result of request
// path: /api_key/api/{path}
func (api *BaseApi) request(rst any, path string, body io.Reader) (err error) {
	if rst == nil {
		return errors.New("rst can not be nil")
	}
	method := http.MethodGet
	if body != nil {
		method = http.MethodPost
	}
	urlAddr := api.url
	if path != "" {
		urlAddr = api.url + path
	}
	req, err := http.NewRequest(method, urlAddr, body)
	if err != nil {
		return
	}
	api.headerLock.Lock()
	req.Header = api.header.Clone()
	api.headerLock.Unlock()
	resp, err := api.cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status + " Body: " + string(data))
	}

	if !api.noExtract {
		var r = response{}
		if err = json.Unmarshal(data, &r); err != nil {
			return
		}
		if r.Error != nil {
			err = fmt.Errorf("unmarshal to response %w, %s", r.Error, string(data))
			return
		}
		data = r.Result
	}
	if err = json.Unmarshal(data, rst); err != nil {
		return
	}
	return nil
}
