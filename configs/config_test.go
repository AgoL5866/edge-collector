package configs

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test123(t *testing.T) {
	var s int64 = 1683617872880205
	ti := time.UnixMicro(s)
	fmt.Println(ti)
}

func TestCheckType(t *testing.T) {
	d := AlertDiff{
		Height:    100,
		CheckType: map[CheckType]bool{CheckTypeBlockBehind: false},
	}
	enable := d.IsEnableCheckType(CheckTypeBlockBehind)
	fmt.Println(enable)
}

func TestDumpConfig(t *testing.T) {
	type args struct {
		c *AppConf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", wantErr: false, args: args{c: &AppConf{
			DebugMode: true,
			Chains: []Chain{
				{Name: "eth", RpcType: "evm", Nodes: map[string]string{
					"official": "https://rpc.eth.org",
					"rockx":    "https://eth.w3node.com/<key>/api",
				}},
				{Name: "bsc", RpcType: "evm", Nodes: map[string]string{
					"official": "https://rpc.bsc.org",
					"rockx":    "https://eth.w3node.com/<key>/api",
				}},
				{Name: "polygon", RpcType: "evm", Nodes: map[string]string{
					"official": "https://rpc.polygon.org",
					"rockx":    "https://eth.w3node.com/<key>/api",
				}},
			},
			Alerts: []Alert{
				{Name: "rockx portal", Slack: "https://hooks.slack.com/xxx", Diff: map[string]AlertDiff{"eth": {Height: 100}, "bsc": {Height: 200}, "polygon": {Height: 200}}},
				{Name: "amber group", Slack: "https://hooks.slack.com/xxx", Diff: map[string]AlertDiff{"eth": {Height: 100}, "bsc": {Height: 200}, "polygon": {Height: 200}}},
			},
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := DumpConfig(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("DumpConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			log.Println(out)
		})
	}
}

func TestParseConfig(t *testing.T) {
	type args struct {
		f string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", wantErr: false, args: args{f: "./block_height.yaml"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if c, err := ParseConfig([]byte(tt.args.f)); (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %+v, wantErr %+v", err, tt.wantErr)
			} else {
				log.Printf("%+v \n", c)
			}
		})
	}
}
