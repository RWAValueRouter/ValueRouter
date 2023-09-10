package jobs

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"sync"
)

type ChainConfig struct {
	Rpc_list       []string
	Chain_name     string
	Chain_currency string
	Is_evm         string
	Chain_id       int
	Chain_type     int
}

var (
	CC = make(map[int]map[int]*ChainConfig)
	Cm sync.RWMutex
)

func listenConfig() {
	log.Info("listenConfig")
	l, err := db.Conn.GetStructValue("select * from chain_config", ChainConfig{})
	if err != nil {
		log.Error("internal db error " + err.Error())
		return
	}
	Cm.Lock()
	defer Cm.Unlock()
	for _, v := range l {
		c := v.(*ChainConfig)
		if exist, ok := CC[c.Chain_type]; ok {
			exist[c.Chain_id] = c
			CC[c.Chain_type] = exist
		} else {
			typ := make(map[int]*ChainConfig)
			typ[c.Chain_id] = c
			CC[c.Chain_type] = typ
		}
	}
}

func init() {
	jobs.AddFunc("@every 30s", listenConfig)
}
