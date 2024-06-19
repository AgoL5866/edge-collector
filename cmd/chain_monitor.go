package cmd

import (
	"sync"
	"time"

	"github.com/coolestowl/edge-collector/configs"
	"github.com/coolestowl/edge-collector/db"
	"github.com/coolestowl/edge-collector/env"
	"github.com/coolestowl/edge-collector/rpc"
	log "github.com/sirupsen/logrus"
)

type Monitor struct {
}

func (m *Monitor) Routine(cfg *configs.AppConf) {
	wg := sync.WaitGroup{}
	for _, c := range cfg.Chains {
		wg.Add(1)

		go func(chain configs.Chain) {
			defer wg.Done()
			m.collectChain(chain)
		}(c)
	}
	wg.Wait()
}

func (m *Monitor) collectChain(chain configs.Chain) {
	defer printPanicStack("[monitor] collectChain recover:")

	lock := sync.Mutex{}
	infoList := make([]*rpc.BlockHeaderInfo, 0, len(chain.Nodes))
	wg := sync.WaitGroup{}
	reqTime := time.Now().Truncate(time.Second).UTC()
	for name, url := range chain.Nodes {
		wg.Add(1)
		go func(name, url string) {
			defer wg.Done()
			info := &rpc.BlockHeaderInfo{
				Region: env.REGION, Chain: chain.Name, NodeAlias: name, ReqTime: reqTime, Timestamp: time.Unix(0, 0),
			}
			rpc := rpc.GetRpcProvider(chain.RpcType, url)

			start := time.Now()
			err := rpc.BlockHeader(info)
			if err != nil {
				log.Errorf("[monitor] get chain=%s node=%s block height: %s", info.Chain, info.NodeAlias, err)
			}
			info.Duration = time.Since(start).Milliseconds()

			//log the info with error to database
			lock.Lock()
			defer lock.Unlock()

			infoList = append(infoList, info)
			if configs.App.DebugMode {
				log.Debugf("[monitor] get block height: %+v", info)
			}
		}(name, url)
	}
	wg.Wait()

	if len(infoList) == 0 {
		log.Infof("[monitor] get block height all failed, %+v", infoList)
		return
	}

	orm := db.GetORM()
	if orm != nil {
		if err := orm.CreateInBatches(infoList, 100).Error; err != nil {
			log.Infof("[monitor] save block height info: %v", err)
		}
	}
	for k, v := range infoList {
		log.Infof("[monitor] collect chain block headers %d: %+v", k, v)
	}
}
