package v1

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Echo(c *gin.Context) {
	var req struct {
		Msg string `json:"msg"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp := errResp(http.StatusBadRequest, err, "BodyParser()")
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, okResp(req.Msg))
}

func Req(c *gin.Context) {
	var req struct {
		Method  string            `json:"method"`
		Uri     string            `json:"uri"`
		Payload string            `json:"payload"`
		Headers map[string]string `json:"headers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		resp := errResp(http.StatusBadRequest, err, "BodyParser()")
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}

	r, err := http.NewRequest(req.Method, req.Uri, strings.NewReader(req.Payload))
	if err != nil {
		resp := errResp(http.StatusBadRequest, err, "NewRequest()")
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}

	start := time.Now()

	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		resp := errResp(http.StatusBadRequest, err, "DoRequest()")
		c.AbortWithStatusJSON(http.StatusBadRequest, resp)
		return
	}
	defer rsp.Body.Close()

	timeUsed := time.Since(start)

	data, _ := io.ReadAll(rsp.Body)

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status":   rsp.StatusCode,
		"body":     string(data),
		"duration": timeUsed / time.Millisecond,
	})
}

func errResp(code int, err error, msg string) any {
	return map[string]any{
		"code": code,
		"msg":  msg,
		"data": err.Error(),
	}
}

func okResp(data any) any {
	return map[string]any{
		"code": 0,
		"msg":  "ok",
		"data": data,
	}
}
