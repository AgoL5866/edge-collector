package configs

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type CheckType string

const (
	CheckTypeBlockMoving CheckType = "BlockMoving"
	CheckTypeNetwork     CheckType = "Network"
	CheckTypeBlockBehind CheckType = "BlockBehind"
	CheckTypeSoftFork    CheckType = "SoftFork"
)

var (
	App = &AppConf{}
)

type AppConf struct {
	DebugMode bool `yaml:"debugMode"`
	Chains    []Chain
	Alerts    []Alert
}

type Chain struct {
	Name    string
	RpcType string            `yaml:"rpcType"`
	Nodes   map[string]string // node alias -> endpoint url
}

type Alert struct {
	Name       string
	Slack      string //slack addr want to push
	SkipNotice bool   `yaml:"skipNotice"`
	Threshold  int
	Period     int
	Diff       map[string]AlertDiff // Chain.Name -> height diff
}

func (a *Alert) GetPeriodByChain(chain string) int {
	p := a.Period
	if ad := a.Diff[chain]; ad.Period > 0 {
		p = ad.Period
	}
	return p
}

func (a *Alert) GetThresholdByChain(chain string) int {
	p := a.Threshold
	if ad := a.Diff[chain]; ad.Threshold > 0 {
		p = ad.Threshold
	}
	return p
}

type AlertDiff struct {
	Height    int64
	Period    int
	Threshold int
	CheckType map[CheckType]bool `yaml:"checkType"` //default true
}

func (d *AlertDiff) IsEnableCheckType(t CheckType) bool {
	enable, ok := d.CheckType[t]
	return !ok || (enable == true)
}

func (c *AppConf) GetChainByName(n string) *Chain {
	for i := range c.Chains {
		ch := c.Chains[i]
		if ch.Name == n {
			return &ch
		}
	}
	return nil
}

func Init(data []byte) error {
	c, err := ParseConfig(data)
	if err != nil {
		return err
	}
	App = c
	return nil
}

func ParseConfig(data []byte) (c *AppConf, err error) {
	c = &AppConf{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	_, _ = DumpConfig(c)
	return c, nil
}

func DumpConfig(c *AppConf) (out string, err error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	out = string(data)
	log.Infof("dump app config: %s", string(data))
	return
}
