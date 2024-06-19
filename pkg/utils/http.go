package utils

import (
	"io"
	"net/http"
)

func DefaultReqJSON(req *http.Request, ptr interface{}) error {
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	return RespParseJSON(rsp, ptr)
}

func RespParseJSON(resp *http.Response, ptr interface{}) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return JSONUnmarshal(data, ptr)
}
