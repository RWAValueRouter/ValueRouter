package jobs

import (
	"encoding/hex"
	"encoding/json"
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/RWAValueRouter/ValueRouter/db"
	"github.com/onrik/ethrpc"
	"strings"
	"sync"
)

var (
	invalidKeygenKeyIdHolder     = make(map[string]bool)
	invalidKeygenKeyIdHolderLock sync.RWMutex
)

// listenKeygenKidStatus listen kengen keyid status and stored it into db
func listenKeygenKidStatus() {
	log.Info("listenKeygenKidStatus")
	list, err := db.Conn.GetStructValue("select ip_port, key_id, uuid from accounts_info where key_id is not null and status = 0 and DATE_ADD(local_system_time, INTERVAL 8 DAY) > NOW() group by key_id, ip_port, uuid", Data{})
	if err != nil {
		log.Error("listenKeygenKidStatus", "error", err.Error())
		return
	}
	for _, l := range list {
		d := l.(*Data)
		go runKeygenKidStatus(d)
	}
}

func runKeygenKidStatus(d *Data) {
	invalidKeygenKeyIdHolderLock.RLock()
	if invalidKeygenKeyIdHolder[d.Key_id] == true {
		invalidKeygenKeyIdHolderLock.RUnlock()
		return
	}
	invalidKeygenKeyIdHolderLock.RUnlock()
	maxStatus, err := db.Conn.GetIntValue("select max(status) from accounts_info where key_id = ?", d.Key_id)
	if err != nil {
		log.Error("internal db " + err.Error())
		return
	}
	if maxStatus >= 2 {
		invalidKeygenKeyIdHolderLock.Lock()
		defer invalidKeygenKeyIdHolderLock.Unlock()
		invalidKeygenKeyIdHolder[d.Key_id] = true
		log.Warn("max status over 2 skip", "key_id", d.Key_id)
		return
	}
	var statusJSON ReqAddrStatus
	uuid := d.Uuid
	client := ethrpc.New("http://" + d.Ip_port)
	log.Info("runKeygenKidStatus request url " + d.Ip_port)
	reqStatus, err := client.Call("smpc_getReqAddrStatus", d.Key_id)
	if err != nil {
		log.Error("smpc_getReqAddrStatus rpc error:" + err.Error())
		return
	}
	statusJSONStr, err := common.GetJSONResult(reqStatus)
	if err != nil {
		log.Error("smpc_getReqAddrStatus=NotStart", "keyID", d.Key_id, "error", err.Error())
		return
	}
	log.Debug("smpc_getReqAddrStatus", "keyId", d.Key_id, "result", statusJSONStr)
	if err := json.Unmarshal([]byte(statusJSONStr), &statusJSON); err != nil {
		log.Error(err.Error())
		return
	}
	if strings.ToLower(statusJSON.Status) != "pending" {
		log.Info("smpc_getReqAddrStatus", "smpc_getReqAddrStatus", statusJSON.Status, "keyID", d.Key_id)
		errMsg := statusJSON.Error
		tipMsg := statusJSON.Tip
		pub := statusJSON.PubKey
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
		for _, reply := range statusJSON.AllReply {
			addr := ""
			if pub != "" {
				pubBuf, err := hex.DecodeString(pub)
				if err != nil {
					log.Error("invalid statusJson public key", "error", err.Error())
					continue
				}
				addr = common.PublicKeyBytesToAddress(pubBuf).String()
			}
			_, err := db.BatchExecute("update accounts_info set error = ? , tip = ? , reply_timestamp = ?, reply_status = ? , initializer = ? , reply_enode = ? ,"+
				"mpc_address = ?, public_key = ? , status = ? where uuid = ? and substring(enode, 9, 128) = ?", tx, errMsg, tipMsg, reply.TimeStamp, reply.Status, reply.Initiator, reply.Enode,
				addr, pub, int(stat), uuid, reply.Enode)
			if err != nil {
				db.Conn.Rollback(tx)
				log.Error("internal db error", "error", err.Error())
				return
			}
			// if status == success only initializer return
			if stat == 1 {
				_, err = db.BatchExecute("update accounts_info set error = ? , tip = ? , reply_timestamp = ?, reply_status = ? , initializer = ? , reply_enode = substring(enode, 9, 128) ,"+
					"mpc_address = ?, public_key = ? , status = ? where uuid = ? and substring(enode, 9, 128) != ?", tx, errMsg, tipMsg, reply.TimeStamp, reply.Status, 0,
					addr, pub, int(stat), uuid, reply.Enode)
				if err != nil {
					db.Conn.Rollback(tx)
					log.Error("internal db error", "error", err.Error())
					return
				}
			}
		}
		db.Conn.Commit(tx)
	}
}

func init() {
	jobs.AddFunc("@every 10s", listenKeygenKidStatus)
}
