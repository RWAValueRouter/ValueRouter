package jobs

import (
	"encoding/hex"
	"encoding/json"
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/RWAValueRouter/ValueRouter/db"
	"github.com/onrik/ethrpc"
	"strings"
)

func listenSigningList() {
	log.Info("listenSigningList")
	signingAccts, err := db.Conn.GetStructValue("select user_account,ip_port from signs_detail where status = 0 and DATE_ADD(local_system_time, INTERVAL 8 DAY) > NOW() group by user_account, ip_port", UserAccount{})
	if err != nil {
		log.Error("listenSigningList", "internal db error", err.Error())
		return
	}
	for _, acct := range signingAccts {
		a := acct.(*UserAccount)
		go runListenSigningList(a)
	}

	go func() {
		_, err = db.Conn.CommitOneRow("UPDATE signing_list JOIN signs_detail ON signing_list.key_id = signs_detail.key_id SET signing_list.status = 1 WHERE signs_detail.status >= 1 and signing_list.status = 0")
		if err != nil {
			log.Error("db error " + err.Error())
			return
		}
	}()
}

func runListenSigningList(acct *UserAccount) {
	// get approve list of condominium account
	client := ethrpc.New("http://" + acct.Ip_port)
	log.Info("runListenSigningList " + acct.Ip_port)
	reqListRep, err := client.Call("smpc_getCurNodeSignInfo", acct.User_account)
	if err != nil {
		log.Error("runListenSigningList", "rpc call", err.Error())
		return
	}
	singingKids, err := db.Conn.GetStructValue("select key_id from signing_list where user_account = ? and ip_port = ? and status = 0", SigningKids{}, acct.User_account, acct.Ip_port)
	if err != nil {
		log.Error("internal db error " + err.Error())
		return
	}
	existed := extractSingingKids(singingKids)
	reqListJSON, _ := common.GetJSONData(reqListRep)

	tx, err := db.Conn.Begin()
	if err != nil {
		log.Error("internal db error " + err.Error())
		return
	}
	if !strings.EqualFold(string(reqListJSON), "null") {
		log.Debug("smpc_getCurNodeSignInfo", "msg", string(reqListJSON))
		var signing []SignCurNodeInfo
		if err = json.Unmarshal(reqListJSON, &signing); err != nil {
			log.Error("Unmarshal SignCurNodeInfo fail:", "msg", err.Error())
			return
		}

		for _, ing := range signing {
			delete(existed, ing.Key)
			c, err := db.Conn.GetIntValue("select count(key_id) from signing_list where key_id = ? and lower(user_account) = ?", ing.Key, strings.ToLower(acct.User_account))
			if err != nil {
				log.Error("internal db error " + err.Error())
				db.Conn.Rollback(tx)
				return
			}
			if c > 0 {
				continue
			}
			pubBuf, err := hex.DecodeString(ing.PubKey)
			if err != nil {
				log.Error("invalid public key", "error", err.Error())
				return
			}
			addr := common.PublicKeyBytesToAddress(pubBuf).String()
			status := 0
			isInitial, err := db.Conn.GetIntValue("select count(key_id) from signs_info where key_id = ? and lower(account) = ?", ing.Key, strings.ToLower(acct.User_account))
			if err != nil {
				log.Error("internal db error " + err.Error())
				db.Conn.Rollback(tx)
				return
			}
			if isInitial > 0 {
				status = 1
			}
			cs, err := db.Conn.GetStructValue("select chain_id,chain_type from signs_detail where key_id = ?", ChainSpec{}, ing.Key)
			if err != nil {
				log.Error("internal db error " + err.Error())
				db.Conn.Rollback(tx)
				return
			}
			var chainId, chainType int
			if len(cs) > 0 {
				chainId = cs[0].(*ChainSpec).Chain_id
				chainType = cs[0].(*ChainSpec).Chain_type
			}
			_, err = db.BatchExecute("insert into signing_list(user_account, group_id, key_id, key_type, `mode`, msg_context, msg_hash, nonce, public_key,mpc_address, threshold, `timestamp`, ip_port, status, chain_id, chain_type) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
				tx, acct.User_account, ing.GroupID, ing.Key, ing.KeyType, ing.Mode, common.ConvertArrStrToStr(ing.MsgContext), common.ConvertArrStrToStr(ing.MsgHash), ing.Nonce, ing.PubKey, addr, ing.ThresHold, ing.TimeStamp, acct.Ip_port, status, chainId, chainType)
			if err != nil {
				log.Error("internal db error " + err.Error())
				db.Conn.Rollback(tx)
				return
			}
		}
	}

	// update status if not exist
	if len(existed) > 0 {
		for k, _ := range existed {
			_, err = db.BatchExecute("update signing_list set status = 1 where key_id = ?", tx, k)
			if err != nil {
				log.Error("internal db error " + err.Error())
				db.Conn.Rollback(tx)
				return
			}
		}
	}
	db.Conn.Commit(tx)
}

func extractSingingKids(d []interface{}) map[string]bool {
	m := make(map[string]bool)
	for _, v := range d {
		kids := v.(*SigningKids)
		m[kids.Key_id] = true
	}
	return m
}

func init() {
	jobs.AddFunc("@every 10s", listenSigningList)
}
