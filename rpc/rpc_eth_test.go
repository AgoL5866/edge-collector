package rpc

import (
	"testing"
	"time"
)

func TestNewEthApi(t *testing.T) {
	type args struct {
		url     string
		host    string
		maxConn int64
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "", args: args{
			url:     "wss://eth.w3node.com/53e984408d8aec59503f1e8c9a93db128131d1f650d947fda53c7ccb79a6f90b/ws",
			host:    "",
			maxConn: 10,
			timeout: 10,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEthApi(tt.args.url, tt.args.host, tt.args.maxConn, tt.args.timeout)
			if got == nil {
				t.Fatal("NewEthApi() = nil")
			}
		})
	}
}
