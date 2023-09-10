package jobs

import (
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/chains/types"
	"github.com/RWAValueRouter/ValueRouter/db"
	"time"
)

type Receipt struct {
	Key_id      string
	Txid        string
	Chain_type  int
	Msg_context string
	Mpc_address string
}

func listenTxsReceipts() {
	log.Info("listenTxsReceipts")
	l, err := db.Conn.GetStructValue("select key_id ,txid, chain_type, msg_context, mpc_address from signs_detail where status = 4 and rsv is not null group by key_id ,txid, chain_type, msg_context, mpc_address", Receipt{})
	if err != nil {
		log.Error("listenTxsReceipts internal db error " + err.Error())
		return
	}
	for _, v := range l {
		r := v.(*Receipt)
		go runCheckTxReceipt(r)
	}
}

func runCheckTxReceipt(r *Receipt) {
	Cm.RLock()
	defer Cm.RUnlock()
	if len(CC) == 0 {
		log.Info("initial not finished")
		return
	}
	config, ok := CC[r.Chain_type]
	var cf *ChainConfig
	var chainID int64
	var err error
	if ok && r.Chain_type == 0 {
		chainID, err = types.EvmChain.GetChainId(r.Msg_context)
		if err != nil {
			log.Error(err.Error())
			return
		}
		cf = config[int(chainID)]
	} else {
		cf = config[0]
	}
	switch r.Chain_type {
	case int(types.EVM):
		isReplaced := false
		tx, err := db.Conn.Begin()
		if err != nil {
			log.Error("internal db error " + err.Error())
			return
		}
		for _, rpc := range cf.Rpc_list {
			ret, err := types.EvmChain.GetTransactionReceipt(r.Txid, rpc)
			if err != nil {
				log.Info("rpc error + " + err.Error())
				continue
			}
			if ret != nil && ret.Status != nil && *ret.Status == 1 { // succ
				_, err = db.BatchExecute("update signs_detail set status = 5 where key_id = ?", tx, r.Key_id)
				if err != nil {
					log.Error("internal db error " + err.Error())
					db.Conn.Rollback(tx)
					return
				}
				break
			} else if ret != nil && ret.Status != nil && *ret.Status == 0 { // fail
				_, err = db.BatchExecute("update signs_detail set status = 6 where key_id = ?", tx, r.Key_id)
				if err != nil {
					log.Error("internal db error " + err.Error())
					db.Conn.Rollback(tx)
					return
				}
				break
			}
		}
		db.Conn.Commit(tx)

		go func() {
		out:
			for _, rpc := range cf.Rpc_list {
				for try := 0; try < 10; try++ {
					time.Sleep(2 * time.Second)
					retByHash, err := types.EvmChain.GetTransactionByHash(r.Txid, rpc)
					if err != nil {
						log.Info("rpc error + " + err.Error())
						continue
					}
					if retByHash != nil {
						if isReplaced == true {
							log.Info("retry success")
						}
						isReplaced = false
						break out
					}
					log.Info("ret is null retry " + r.Key_id)
					isReplaced = true
				}
				if isReplaced {
					log.Warnf("ret is null " + r.Txid + " rpc " + rpc)
				}
			}
			if isReplaced {
				_, err = db.Conn.CommitOneRow("update signs_detail set status = 8 where key_id = ?", r.Key_id)
				if err != nil {
					log.Error("internal db error " + err.Error())
					return
				}
			}
		}()
	default:
		log.Error("not supported chain type")
	}
}

func init() {
	jobs.AddFunc("@every 10s", listenTxsReceipts)
}
