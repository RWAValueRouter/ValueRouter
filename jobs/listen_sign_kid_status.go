package jobs

import (
	"encoding/json"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"github.com/anyswap/fastmpc-service-middleware/db"
	"github.com/onrik/ethrpc"
	"math/big"
	"strings"
)

// listenSignKidStatus listen sign keyid status and stored it into db
func listenSignKidStatus() {
	log.Info("listenSignKidStatus")
	list, err := db.Conn.GetStructValue("select key_id,pubkey from signs_info where key_id is not null and status = 0 and DATE_ADD(local_system_time, INTERVAL 8 DAY) > NOW()", SignKids{})
	if err != nil {
		log.Error("listenKeygenKidStatus", "error", err.Error())
		return
	}
	for _, l := range list {
		d := l.(*SignKids)
		go runSignKidStatus(d)
	}
}

func runSignKidStatus(d *SignKids) {
	kid := d.Key_id
	ps, err := db.Conn.GetStructValue("select ip_port,user_account,enode from accounts_info where public_key = ? and status = 1", IpData{}, d.Pubkey)
	if err != nil {
		log.Error("internal db " + err.Error())
		return
	}
	for _, p := range ps {
		ip := p.(*IpData)
		client := ethrpc.New("http://" + ip.Ip_port)
		reqStatus, err := client.Call("smpc_getSignStatus", kid)
		if err != nil {
			log.Error("smpc_getSignStatus rpc error:" + err.Error())
			continue
		}
		statusJSONStr, err := common.GetJSONResult(reqStatus)
		if err != nil {
			log.Error("smpc_getSignStatus=NotStart", "keyID", d.Key_id, "error", err.Error())
			continue
		}
		var statusJSON SignStatus
		log.Debug("smpc_getSignStatus", "keyId", d.Key_id, "result", statusJSONStr)
		if err := json.Unmarshal([]byte(statusJSONStr), &statusJSON); err != nil {
			log.Error(err.Error())
			return
		}
		if strings.ToLower(statusJSON.Status) != "pending" {
			log.Info("smpc_getSignStatus", "smpc_getSignStatus", statusJSON.Status, "keyID", d.Key_id)
			stat, ok := ReqAddressStatusMap[strings.ToLower(statusJSON.Status)]
			if !ok {
				log.Error("can not find status in ReqAddressStatusMap")
				return
			}
			tx, err := db.Conn.Begin()
			if err != nil {
				log.Error("internal db " + err.Error())
				return
			}
			var maxTimestamp = new(big.Int).SetInt64(0)
			for _, reply := range statusJSON.AllReply {
				if reply.Approver == "" {
					db.Conn.Rollback(tx)
					log.Warn("reply approver can not be blank " + kid)
					return
				}
				exist, err := db.Conn.GetIntValue("select count(key_id) from signs_detail where key_id = ? and substring(enode, 9, 128) = ? and status != 0", kid, reply.Enode)
				if err != nil {
					db.Conn.Rollback(tx)
					log.Error("internal db error", "error", err.Error())
					return
				}
				if exist > 0 {
					log.Info("runSignKidStatus already exist record")
					continue
				}
				_, err = db.BatchExecute("update signs_detail set rsv = ? , reply_initializer = ? , reply_status = ? , reply_timestamp = ? , reply_enode = ? , initiator_public_key = ? , error = ? , tip = ? , status = ? where key_id = ? and substring(enode, 9, 128) = ?", tx,
					common.ConvertArrStrToStr(statusJSON.Rsv), reply.Initiator, reply.Status, reply.TimeStamp, reply.Enode, statusJSON.Initiator, statusJSON.Error, statusJSON.Tip, int(stat), kid, reply.Enode)
				if err != nil {
					db.Conn.Rollback(tx)
					log.Error("internal db error", "error", err.Error())
					return
				}

				if reply.TimeStamp != "" {
					tmpC, ok := new(big.Int).SetString(reply.TimeStamp, 10)
					if ok {
						if maxTimestamp.Cmp(tmpC) == -1 {
							maxTimestamp = tmpC
						}
					} else {
						db.Conn.Rollback(tx)
						log.Error("reply.TimeStamp value is fatal", "error", err.Error())
						return
					}
				}

				_, err = db.BatchExecute("update signing_list set status = 1 where user_account = ? and key_id = ?", tx, strings.ToLower(reply.Approver), kid)
				if err != nil {
					db.Conn.Rollback(tx)
					log.Error("internal db error", "error", err.Error())
					return
				}
			}

			_, err = db.BatchExecute("update signs_info set status = 1 where key_id = ?", tx, kid)
			if err != nil {
				db.Conn.Rollback(tx)
				log.Error("internal db error", "error", err.Error())
				return
			}

			if maxTimestamp.Cmp(new(big.Int).SetInt64(0)) != 0 {
				_, err = db.BatchExecute("update signs_detail set reply_timestamp = ? where key_id = ?", tx, maxTimestamp.String(), kid)
				if err != nil {
					db.Conn.Rollback(tx)
					log.Error("internal db error", "error", err.Error())
					return
				}
			}

			db.Conn.Commit(tx)
		}
	}
}

func init() {
	jobs.AddFunc("@every 10s", listenSignKidStatus)
}
