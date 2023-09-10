package service

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/RWAValueRouter/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/chains/types"
	"github.com/RWAValueRouter/ValueRouter/common"
	"github.com/RWAValueRouter/ValueRouter/db"
	common2 "github.com/RWAValueRouter/ValueRouter/internal/common"
	"github.com/RWAValueRouter/ValueRouter/jobs"
	"github.com/google/uuid"
	"github.com/onrik/ethrpc"
	"strconv"
	"strings"
)

func doGetLatestMpcAddressStatus(mpc_address string, chain_id, chain_type int) (interface{}, error) {
	if !common.CheckEthereumAddress(mpc_address) {
		return nil, errors.New("invalid mpc address")
	}
	status, err := db.Conn.GetIntValue("select max(status) from signs_detail where key_id = (select key_id from signs_detail where lower(mpc_address) =  ? and chain_id = ? and chain_type = ? and sign_type = 0 order by id desc limit 1)", strings.ToLower(mpc_address), chain_id, chain_type)
	if err != nil {
		return nil, errors.New("db error " + err.Error())
	}
	return status, nil
}

func doGetNodesNumber() (interface{}, error) {
	c, err := db.Conn.GetIntValue("select count(status) from nodes_info where status = 0")
	if err != nil {
		return nil, errors.New("db error " + err.Error())
	}
	return c, nil
}

func doGetApprovalListByKeyId(key_id string) (interface{}, error) {
	if !common.ValidateKeyId(key_id) {
		return nil, errors.New("invalid key id")
	}
	l, err := db.Conn.GetStructValue("select a.*, b.reply_status, b.reply_initializer from signing_list a left join signs_detail b on a.user_account = b.user_account and a.key_id = b.key_id where  a.key_id = ?", SignCurNodeInfoDetail{}, strings.ToLower(key_id))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	c, err := db.Conn.GetIntValue("select count(status) from signing_list where key_id = ? and status = 1", strings.ToLower(key_id))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	initial_account, err := db.Conn.GetStringValue("select account from signs_info where key_id = ?", strings.ToLower(key_id))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	for _, v := range l {
		s := v.(*SignCurNodeInfoDetail)
		if strings.EqualFold(s.User_account, initial_account) {
			s.Reply_initializer = 1
		}
		s.Signed = c
		leastApproval, err := strconv.Atoi(s.Mode[0:1])
		if err != nil {
			return nil, errors.New("internal logic error Mode value invalid")
		}
		if c >= leastApproval {
			if s.Status == 1 && s.Reply_initializer == 0 && common.IsBlank(s.Reply_status) {
				s.Reply_status = "TIMEOUT"
			}
		}

		if s.Reply_initializer == 1 {
			s.Reply_status = "AGREE"
		}

		cs, err := db.Conn.GetStructValue("select chain_id,chain_type,sign_type,rsv from signs_detail where key_id = ?", jobs.ChainSpec{}, s.Key_id)
		if err != nil {
			return nil, err
		}
		if len(cs) > 0 {
			s.Chain_id = cs[0].(*jobs.ChainSpec).Chain_id
			s.Chain_type = cs[0].(*jobs.ChainSpec).Chain_type
			s.Sign_type = cs[0].(*jobs.ChainSpec).Sign_type
		}

		for _, sig := range cs {
			rsv := sig.(*jobs.ChainSpec).Rsv
			if rsv != "" {
				s.Rsv = rsv
			}
		}
	}
	return l, nil
}

func doGetMpcAddressDetail(mpc_account string) (interface{}, error) {
	if !common.CheckEthereumAddress(mpc_account) {
		return nil, errors.New("invalid account")
	}

	l, err := db.Conn.GetStructValue("select a.status, a.user_account, a.key_id, a.public_key, a.mpc_address, a.initializer, a.reply_status ,a.reply_timestamp ,a.reply_enode, a.gid , a.threshold, b.key_type , b.mode from accounts_info a, groups_info b where a.uuid = b.uuid and lower(mpc_address) = ?", RespAddr{}, strings.ToLower(mpc_account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return l, nil
}

func doGetAsset(account string, chain_id, chain_type int) (interface{}, error) {
	if !common.CheckEthereumAddress(account) {
		return nil, errors.New("invalid account")
	}

	asset, err := db.Conn.GetStringValue("select asset from asset_info where lower(mpc_address) = ? and chain_id = ? and chain_type = ?", strings.ToLower(account), chain_id, chain_type)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	type AddAsset struct {
		Symbol   string `json:"Symbol"`
		Contract string `json:"Contract"`
		Name     string `json:"Name"`
		Decimal  int    `json:"Decimal"`
	}
	var assets []AddAsset
	all := strings.Split(asset, "|")
	for _, a := range all {
		inner := strings.Split(a, "&")
		if len(inner) == 4 {
			decimal, err := strconv.Atoi(inner[3])
			if err != nil {
				log.Error("invalid db data decimal convert error" + err.Error())
				continue
			}
			assets = append(assets, AddAsset{
				Symbol:   inner[0],
				Contract: inner[1],
				Name:     inner[2],
				Decimal:  decimal,
			})
		}
	}
	return assets, nil
}

func doAddAssetForMpcAddress(rsv string, msg string) (interface{}, error) {
	log.Info("doAddAssetForMpcAddress request param", "rsv", rsv, "msg", msg)
	if err := common.VerifyAccount(rsv, msg); err != nil {
		return nil, err
	}
	req := AddAssetForMpcAddress{}
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}

	if req.TxType != "ADDASSETFORMPCADDRESS" {
		return nil, errors.New("request tx type invalid")
	}
	if common.IsSomeOneBlank(req.Nonce, req.Symbol, req.Contract, req.TimeStamp, req.Name, req.MpcAddress) {
		return nil, errors.New("request param can not be blank")
	}

	if !common.CheckEthereumAddress(req.Contract) {
		return nil, errors.New("contract address is invalid")
	}

	if !common.CheckEthereumAddress(req.MpcAddress) {
		return nil, errors.New("mpc_address address is invalid")
	}

	if req.Decimal == 0 {
		return nil, errors.New("decimal is invalid")
	}
	isOwnMpc, err := db.Conn.GetIntValue("select count(status) from accounts_info where lower(mpc_address) = ? and lower(user_account) = ? and status = 1", strings.ToLower(req.MpcAddress), strings.ToLower(req.Account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	if isOwnMpc == 0 {
		return nil, errors.New("account " + req.Account + " does not own " + req.MpcAddress)
	}

	req.Account = req.MpcAddress
	var asset string
	var exist bool
	asset, err = db.Conn.GetStringValue("select asset from asset_info where lower(mpc_address) = ? and chain_type = ? and chain_id = ?", strings.ToLower(req.Account), req.ChainType, req.ChainId)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	if _, ok := common.StripAsset(asset)[strings.ToLower(req.Contract)]; ok {
		return nil, errors.New("already exist asset " + req.Contract)
	}
	if asset != "" {
		exist = true
	}
	asset += strings.ToLower(req.Symbol) + "&" + strings.ToLower(req.Contract) + "&" + strings.ToLower(req.Name) + "&" + strings.ToLower(strconv.Itoa(req.Decimal)) + "|"
	if exist == true {
		_, err = db.Conn.CommitOneRow("update asset_info set asset = ? where mpc_address = ? and chain_id = ? and chain_type = ?", strings.ToLower(asset), strings.ToLower(req.Account), req.ChainId, req.ChainType)
	} else {
		_, err = db.Conn.CommitOneRow("insert into asset_info(mpc_address,asset,chain_id,chain_type) values(?,?,?,?)", strings.ToLower(req.Account), strings.ToLower(asset), req.ChainId, req.ChainType)
	}
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return "success", nil
}

func doAddAsset(rsv string, msg string) (interface{}, error) {
	log.Info("doAddAsset request param", "rsv", rsv, "msg", msg)
	if err := common.VerifyAccount(rsv, msg); err != nil {
		return nil, err
	}
	req := AddAsset{}
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}

	if req.TxType != "ADDASSET" {
		return nil, errors.New("request tx type invalid")
	}
	if common.IsSomeOneBlank(req.Nonce, req.Symbol, req.Contract, req.TimeStamp, req.Name) {
		return nil, errors.New("request param can not be blank")
	}

	if !common.CheckEthereumAddress(req.Contract) {
		return nil, errors.New("contract address is invalid")
	}

	if req.Decimal == 0 {
		return nil, errors.New("decimal is invalid")
	}

	var asset string
	var exist bool
	asset, err = db.Conn.GetStringValue("select asset from asset_info where lower(mpc_address) = ? and chain_type = ? and chain_id = ?", strings.ToLower(req.Account), req.ChainType, req.ChainId)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	if _, ok := common.StripAsset(asset)[strings.ToLower(req.Contract)]; ok {
		return nil, errors.New("already exist asset " + req.Contract)
	}
	if asset != "" {
		exist = true
	}
	asset += strings.ToLower(req.Symbol) + "&" + strings.ToLower(req.Contract) + "&" + strings.ToLower(req.Name) + "&" + strings.ToLower(strconv.Itoa(req.Decimal)) + "|"
	if exist == true {
		_, err = db.Conn.CommitOneRow("update asset_info set asset = ? where mpc_address = ? and chain_id = ? and chain_type = ?", strings.ToLower(asset), strings.ToLower(req.Account), req.ChainId, req.ChainType)
	} else {
		_, err = db.Conn.CommitOneRow("insert into asset_info(mpc_address,asset,chain_id,chain_type) values(?,?,?,?)", strings.ToLower(req.Account), strings.ToLower(asset), req.ChainId, req.ChainType)
	}
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return "success", nil
}

func getTxStatusByKeyId(keyId string) (interface{}, error) {
	if !common.ValidateKeyId(keyId) {
		return nil, errors.New("invalid keyId")
	}
	status, err := db.Conn.GetIntValue("select status from signs_detail where key_id = ? limit 1", keyId)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return status, nil
}

func getTxHashByKeyId(keyId string) (interface{}, error) {
	if !common.ValidateKeyId(keyId) {
		return nil, errors.New("invalid keyId")
	}
	txHash, err := db.Conn.GetStringValue("select txid from signs_detail where key_id = ?", keyId)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	return txHash, nil
}

func acceptSign(rsv string, msg string) (interface{}, error) {
	log.Info("acceptSign request param", "rsv", rsv, "msg", msg)
	err := common.VerifyAccount(rsv, msg)
	if err != nil {
		return nil, err
	}
	req := AcceptSignData{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}
	if len(req.MsgHash) != len(req.MsgContext) {
		return nil, errors.New("message hash and message context length not match")
	}

	if len(req.MsgHash) == 0 {
		return nil, errors.New("message hash and message context can not be blank")
	}

	if common.IsSomeOneBlank(req.Accept, req.Nonce, req.Key, req.TxType, req.TimeStamp) {
		return nil, errors.New("request param not valid")
	}

	if req.TxType != "ACCEPTSIGN" {
		return nil, errors.New("invalid tx type")
	}

	if req.Accept != "AGREE" && req.Accept != "DISAGREE" {
		return nil, errors.New("invalid accept value")
	}

	if c, err := db.Conn.GetIntValue("select sign_type from signs_detail where key_id = ?", req.Key); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	} else {
		if c != -1 {
			switch types.ChainType(req.ChainType) {
			case types.EVM:
				for i, hash := range req.MsgHash {
					if !types.EvmChain.ValidateUnsignedTransactionHash(req.MsgContext[i], hash) &&
						!types.EvmChain.ValidatePersonalSignHash(req.MsgContext[i], hash) &&
						!types.EvmChain.ValidateETHSignHash(req.MsgContext[i], hash) &&
						!types.EvmChain.ValidateTypedHash(req.MsgContext[i], hash) {
						return nil, errors.New("message hash and msg context value not match")
					}
				}
			default:
				return nil, errors.New("unrecognized chain type")
			}
		}
	}

	if c, err := db.Conn.GetIntValue("select count(key_id) from signs_detail where key_id = ? and user_account = ? and msg_hash = ? and status = 0", req.Key, strings.ToLower(req.Account), common.ConvertArrStrToStr(req.MsgHash)); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	} else if c == 0 {
		return nil, errors.New("request param invalid")
	}

	ipPort, err := db.Conn.GetStringValue("select ip_port from signs_detail where key_id = ? and user_account = ?", req.Key, strings.ToLower(req.Account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	if common.IsBlank(ipPort) {
		return nil, errors.New("invalid request param can not find ip port")
	}

	client := ethrpc.New("http://" + ipPort)
	// send rawTx
	acceptSignRep, err := client.Call("smpc_acceptSigning", rsv, msg)
	if err != nil {
		return nil, err
	}
	// get result
	acceptRet, err := common.GetJSONResult(acceptSignRep)
	if err != nil {
		return nil, err
	}
	log.Info("smpc_acceptSign: ", "result", acceptRet)

	_, err = db.Conn.CommitOneRow("insert into signs_result(key_type,account,nonce,key_id,msg_hash,msg_context,timestamp,accept) values(?,?,?,?,?,?,?,?)",
		req.TxType, strings.ToLower(req.Account), req.Nonce, req.Key, common.ConvertArrStrToStr(req.MsgHash), common.ConvertArrStrToStr(req.MsgContext), req.TimeStamp, req.Accept)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	_, err = db.Conn.CommitOneRow("update signs_detail set reply_status = ? where key_id = ? and lower(user_account) = ?", strings.ToUpper(req.Accept), req.Key, strings.ToLower(req.Account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return acceptRet, nil
}

func getSignHistoryByPagination(userAccount string, page, pageSize int) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("user account is not valid")
	}
	if pageSize > 100 {
		return nil, errors.New("exceed max pageSize 100")
	}
	if page < 1 {
		return nil, errors.New("minimum page is 1")
	}
	l, err := db.Conn.GetStructValue("select *,UNIX_TIMESTAMP(a.local_system_time) * 1000 as local_timestamp from signs_detail a where user_account = ? order by id desc limit ?,?", SignHistory{}, strings.ToLower(userAccount), (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	for _, v := range l {
		s := v.(*SignHistory)
		c, err := db.Conn.GetIntValue("select count(status) from signing_list where key_id = ? and status = 1", s.Key_id)
		if err != nil {
			return nil, errors.New("internal db error " + err.Error())
		}
		s.Signed = c

		cs, err := db.Conn.GetStructValue("select rsv from signs_detail where key_id = ?", jobs.ChainSpec{}, s.Key_id)
		if err != nil {
			return nil, err
		}
		for _, sig := range cs {
			rsv := sig.(*jobs.ChainSpec).Rsv
			if rsv != "" {
				s.Rsv = rsv
			}
		}
	}
	return l, nil
}

func getSignHistory(userAccount string) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("user account is not valid")
	}
	l, err := db.Conn.GetStructValue("select *,UNIX_TIMESTAMP(a.local_system_time) * 1000 as local_timestamp from signs_detail a where user_account = ? order by a.id desc limit 100", SignHistory{}, strings.ToLower(userAccount))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	for _, v := range l {
		s := v.(*SignHistory)
		c, err := db.Conn.GetIntValue("select count(status) from signing_list where key_id = ? and status = 1", s.Key_id)
		if err != nil {
			return nil, errors.New("internal db error " + err.Error())
		}
		s.Signed = c

		cs, err := db.Conn.GetStructValue("select rsv from signs_detail where key_id = ?", jobs.ChainSpec{}, s.Key_id)
		if err != nil {
			return nil, err
		}
		for _, sig := range cs {
			rsv := sig.(*jobs.ChainSpec).Rsv
			if rsv != "" {
				s.Rsv = rsv
			}
		}
	}
	return l, nil
}

func getApprovalListByPagination(userAccount string, status, page, pageSize int) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("user account is not valid")
	}
	if status != 0 && status != 1 {
		return nil, errors.New("status can only be 0 or 1")
	}
	if pageSize > 100 {
		return nil, errors.New("exceed max pageSize 100")
	}
	if page < 1 {
		return nil, errors.New("minimum page is 1")
	}
	l, err := db.Conn.GetStructValue("select * from signing_list where user_account = ? and status = ? order by id desc limit ?,?", SignCurNodeInfo{}, strings.ToLower(userAccount), status, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	for _, v := range l {
		s := v.(*SignCurNodeInfo)
		c, err := db.Conn.GetIntValue("select count(status) from signing_list where key_id = ? and status = 1", s.Key_id)
		if err != nil {
			return nil, errors.New("internal db error " + err.Error())
		}
		s.Signed = c
		cs, err := db.Conn.GetStructValue("select chain_id,chain_type,sign_type,rsv from signs_detail where key_id = ?", jobs.ChainSpec{}, s.Key_id)
		if err != nil {
			return nil, err
		}
		if len(cs) > 0 {
			s.Chain_id = cs[0].(*jobs.ChainSpec).Chain_id
			s.Chain_type = cs[0].(*jobs.ChainSpec).Chain_type
			s.Sign_type = cs[0].(*jobs.ChainSpec).Sign_type
		}
		for _, sig := range cs {
			rsv := sig.(*jobs.ChainSpec).Rsv
			if rsv != "" {
				s.Rsv = rsv
			}
		}
	}
	return l, nil
}

func getApprovalList(userAccount string) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("user account is not valid")
	}
	l, err := db.Conn.GetStructValue("select * from signing_list where user_account = ? order by id desc limit 100", SignCurNodeInfo{}, strings.ToLower(userAccount))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	for _, v := range l {
		s := v.(*SignCurNodeInfo)
		c, err := db.Conn.GetIntValue("select count(status) from signing_list where key_id = ? and status = 1", s.Key_id)
		if err != nil {
			return nil, errors.New("internal db error " + err.Error())
		}
		s.Signed = c
		cs, err := db.Conn.GetStructValue("select chain_id,chain_type,sign_type,rsv from signs_detail where key_id = ?", jobs.ChainSpec{}, s.Key_id)
		if err != nil {
			return nil, err
		}
		if len(cs) > 0 {
			s.Chain_id = cs[0].(*jobs.ChainSpec).Chain_id
			s.Chain_type = cs[0].(*jobs.ChainSpec).Chain_type
			s.Sign_type = cs[0].(*jobs.ChainSpec).Sign_type
		}

		for _, sig := range cs {
			rsv := sig.(*jobs.ChainSpec).Rsv
			if rsv != "" {
				s.Rsv = rsv
			}
		}
	}
	return l, nil
}

func getUnsigedTransactionHash(unsignedTx string, chain int) (interface{}, error) {
	var c types.Chain
	switch types.ChainType(chain) {
	case types.EVM:
		c = types.EvmChain
	default:
		return nil, errors.New("unrecognized chain")
	}
	hash, err := c.GetUnsignedTransactionHash(unsignedTx)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func doSign(rsv string, msg string) (interface{}, error) {
	log.Info("doSign request param", "rsv", rsv, "msg", msg)
	err := common.VerifyAccount(rsv, msg)
	if err != nil {
		return nil, err
	}
	req := TxDataSign{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}
	if len(req.MsgHash) != len(req.MsgContext) {
		return nil, errors.New("message hash and message context length not match")
	}

	if len(req.MsgHash) != 1 {
		return nil, errors.New("message hash and message context len must be one")
	}
	chainId := 0
	var sign_type int
	switch types.ChainType(req.ChainType) {
	case types.EVM:
		if types.EvmChain.ValidateETHSignHash(req.MsgContext[0], req.MsgHash[0]) {
			sign_type = 2
		} else if types.EvmChain.ValidatePersonalSignHash(req.MsgContext[0], req.MsgHash[0]) {
			sign_type = 1
		} else if types.EvmChain.ValidateUnsignedTransactionHash(req.MsgContext[0], req.MsgHash[0]) {
			sign_type = 0
		} else if types.EvmChain.ValidateTypedHash(req.MsgContext[0], req.MsgHash[0]) {
			sign_type = 99
		} else {
			return nil, errors.New("unrecognized hash")
		}

		msgContext := MsgContext{}
		err = json.Unmarshal([]byte(req.MsgContext[0]), &msgContext)
		if err != nil {
			return nil, err
		}
		if len(msgContext.ChainId) < 2 {
			return nil, errors.New("invalid chainId")
		}
		ret, err := strconv.ParseInt(msgContext.ChainId[2:], 16, 64)
		if err != nil {
			return nil, err
		}
		chainId = int(ret)
	default:
		return nil, errors.New("unrecognized chain type")
	}

	if req.TxType != "SIGN" {
		return nil, errors.New("tx type must be SIGN")
	}

	if chainId == 0 {
		return nil, errors.New("chainId can not be 0")
	}

	if common.IsSomeOneBlank(req.Nonce, req.PubKey, req.Keytype, req.GroupID, req.ThresHold, req.Mode, req.AcceptTimeOut, req.TimeStamp) {
		return nil, errors.New("invalid request param")
	}

	ipPort, err := db.Conn.GetStringValue("select ip_port from accounts_info where public_key = ? and user_account = ? and status = 1", req.PubKey, strings.ToLower(req.Account))
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	if ipPort == "" {
		return nil, errors.New("can not find pub key and account responded ip port")
	}
	client := ethrpc.New("http://" + ipPort)
	reqKeyID, err := client.Call("smpc_signing", rsv, msg)
	if err != nil {
		return nil, err
	}
	keyID, err := common.GetJSONResult(reqKeyID)
	if err != nil {
		return nil, err
	}
	log.Info("smpc_sign keyID = %s", keyID)

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error" + err.Error())
	}
	_, err = db.BatchExecute("insert into signs_info(account,nonce,pubkey,msg_hash,msg_context,key_type,group_id,threshold,`mod`,accept_timeout,`timestamp`,key_id) values(?,?,?,?,?,?,?,?,?,?,?,?)",
		tx, strings.ToLower(req.Account), req.Nonce, req.PubKey, common.ConvertArrStrToStr(req.MsgHash), common.ConvertArrStrToStr(req.MsgContext), req.Keytype, req.GroupID, req.ThresHold, req.Mode, req.AcceptTimeOut, req.TimeStamp, keyID)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error" + err.Error())
	}
	pubBuf, err := hex.DecodeString(req.PubKey)
	if err != nil {
		return nil, errors.New("invalid req pub key")
	}
	addr := common.PublicKeyBytesToAddress(pubBuf).String()
	accts, err := db.Conn.GetStructValue("select user_account, enode, ip_port from accounts_info where public_key = ?", Account{}, req.PubKey)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error" + err.Error())
	}
	if len(accts) == 0 {
		db.Conn.Rollback(tx)
		return nil, errors.New("invalid public key")
	}
	for _, acct := range accts {
		a := acct.(*Account)
		_, err = db.BatchExecute("insert into signs_detail(key_id, user_account, group_id, threshold, msg_hash, msg_context, public_key, mpc_address, key_type, mode, status, enode, ip_port, chain_type, chain_id,sign_type) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
			tx, keyID, a.User_account, req.GroupID, req.ThresHold, common.ConvertArrStrToStr(req.MsgHash), common.ConvertArrStrToStr(req.MsgContext), req.PubKey, addr, req.Keytype, req.Mode, 0, a.Enode, a.Ip_port, req.ChainType, chainId, sign_type)
		if err != nil {
			db.Conn.Rollback(tx)
			return nil, errors.New("internal db error" + err.Error())
		}
	}

	db.Conn.Commit(tx)
	return keyID, nil
}

func getAccountList(userAccount string) (interface{}, error) {
	if !common.CheckEthereumAddress(userAccount) {
		return nil, errors.New("invalid userAccount")
	}

	l, err := db.Conn.GetStructValue("select * from accounts_info a, groups_info b where a.uuid = b.uuid and user_account = ? and status = 1", RespAddr{}, strings.ToLower(userAccount))
	if err != nil {
		return nil, err
	}

	return l, nil
}

func getReqAddrStatus(keyId string) (interface{}, error) {
	if !common.ValidateKeyId(keyId) {
		return nil, errors.New("keyId is not valid")
	}

	v, err := db.Conn.GetStructValue("select a.status, a.user_account, a.key_id, a.public_key, a.mpc_address, a.initializer, a.reply_status ,a.reply_timestamp ,a.reply_enode, a.gid , a.threshold, b.key_type , b.mode from accounts_info a, groups_info b where a.uuid = b.uuid and a.key_id = ? order by a.id asc", RespAddr{}, keyId)
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	if len(v) == 0 {
		return nil, errors.New("no such keyId")
	}

	return v, nil
}

func doKeyGenByRawData(raw string) (interface{}, error) {
	type Msg struct {
		Rsv string
		Msg string
	}
	m := Msg{}
	err := json.Unmarshal(common2.FromHex(raw), &m)
	if err != nil {
		return nil, err
	}
	return doKeyGen(m.Rsv, m.Msg)
}

func doKeyGen(rsv string, msg string) (interface{}, error) {
	log.Info("doKeyGen request param", "rsv", rsv, "msg", msg)
	err := common.VerifyAccount(rsv, msg)
	if err != nil {
		return nil, err
	}
	req := TxDataReqAddr{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return nil, err
	}
	if req.TxType != "REQSMPCADDR" {
		return nil, errors.New("tx type must be REQSMPCADDR")
	}
	if common.IsSomeOneBlank(req.Nonce, req.Keytype, req.AcceptTimeOut, req.TimeStamp, req.Sigs) {
		return nil, errors.New("request param invalid")
	}
	if req.Mode != "2" {
		return nil, errors.New("service keygen mod must be 2")
	}
	req.Account = strings.ToLower(req.Account)
	ipStr, err := db.Conn.GetStringValue("select ip_port from accounts_info where gid = ? and user_account = ? and threshold = ? and key_id is null and uuid = ?",
		req.GroupID, req.Account, req.ThresHold, req.Uuid)
	if err != nil {
		return nil, errors.New("GetStringValue error " + err.Error())
	}
	if ipStr == "" {
		return nil, errors.New("Can not find ip port through uuid " + req.Uuid)
	}
	client := ethrpc.New("http://" + ipStr)
	reqKeyID, err := client.Call("smpc_reqKeyGen", rsv, msg)
	if err != nil {
		return nil, errors.New("smpc_reqKeyGen error " + err.Error())
	}
	keyID, err := common.GetJSONResult(reqKeyID)
	if err != nil {
		return nil, errors.New("getJSONResult error" + err.Error())
	}
	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	_, err = db.BatchExecute("update accounts_info set key_id = ? where uuid = ?",
		tx, keyID, req.Uuid)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error " + err.Error())
	}
	_, err = db.BatchExecute("insert into groups_info(tx_type, account, nonce, key_type, group_id, thres_hold, mode, accept_timeout, sigs , key_id, uuid, timestamp) "+
		"values(?,?,?,?,?,?,?,?,?,?,?,?)", tx, req.TxType, req.Account, req.Nonce, req.Keytype, req.GroupID, req.ThresHold, req.Mode, req.AcceptTimeOut, req.Sigs, keyID, req.Uuid, req.TimeStamp)
	if err != nil {
		db.Conn.Rollback(tx)
		return nil, errors.New("internal db error " + err.Error())
	}
	if err = db.Conn.Commit(tx); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return keyID, nil
}

func getGroupIdByRawData(raw string) (interface{}, error) {
	type Msg struct {
		Threshold                 string
		UserAccountsAndIpPortAddr []string
	}
	m := Msg{}
	err := json.Unmarshal(common2.FromHex(raw), &m)
	if err != nil {
		return nil, err
	}
	return getGroupId(m.Threshold, m.UserAccountsAndIpPortAddr)
}

// getGroupId threshold 2/3, userAccountsAndIpPortAddr user1|ip:port user2 user3|ip:port
func getGroupId(threshold string, userAccountsAndIpPortAddr []string) (interface{}, error) {
	_, p2, err := common.CheckThreshold(threshold)
	if err != nil {
		return nil, err
	}
	accounts, ipPorts, err := common.CheckUserAccountsAndIpPortAddr(userAccountsAndIpPortAddr)
	if err != nil {
		return nil, err
	}
	if p2 != len(accounts) {
		return nil, errors.New("threshold value and accounts number can not match")
	}
	var selectCount int
	var filledIpPort string
	for _, v := range ipPorts {
		if v == "" {
			selectCount++
		} else {
			filledIpPort = filledIpPort + v + ","
		}
	}
	if filledIpPort != "" {
		filledIpPort = filledIpPort[0 : len(filledIpPort)-1]
	}
	filledIpPortCount, err := db.Conn.GetIntValue("select count(*) from nodes_info where ip_addr in (?)", filledIpPort)
	if err != nil {
		return nil, err
	}
	if filledIpPortCount != len(ipPorts)-selectCount {
		return nil, errors.New("filled ip port not valid")
	}
	type ip struct {
		Ip_addr string
	}
	totalNodes, err := db.Conn.GetIntValue("select count(ip_addr) from nodes_info")
	if err != nil {
		return nil, err
	}
	if totalNodes < selectCount {
		return nil, errors.New("total nodes less than needed")
	}
	allIp, err := db.Conn.GetStructValue("select ip_addr from nodes_info where ip_addr not in (?) ORDER BY RAND() limit ?", ip{}, filledIpPort, selectCount)
	if err != nil {
		return nil, err
	}
	index := 0
	var acct_enode []account_enode
	var enodeList []string
	for i, acc := range accounts {
		var enode string
		var ipStr string
		if ipPorts[i] == "" {
			ipStr = allIp[index].(*ip).Ip_addr
			index++
		} else {
			ipStr = ipPorts[i]
		}
		client := ethrpc.New("http://" + ipStr)
		enodeRep, err := client.Call("smpc_getEnode")
		if err != nil {
			return nil, errors.New("IP_addr" + allIp[index].(ip).Ip_addr + " can not reach")
		}
		log.Info(fmt.Sprintf("getEnode = %s\n\n", enodeRep))
		type dataEnode struct {
			Enode string `json:"Enode"`
		}
		var enodeJSON dataEnode
		enodeData, _ := common.GetJSONData(enodeRep)
		if err := json.Unmarshal(enodeData, &enodeJSON); err != nil {
			return nil, err
		}
		log.Info(fmt.Sprintf("enode = %s\n", enodeJSON.Enode))
		enode = enodeJSON.Enode

		acct_enode = append(acct_enode, account_enode{
			Enode:   enode,
			Account: acc,
			Ip_port: ipStr,
		})

		enodeList = append(enodeList, enode)
	}

	client := ethrpc.New("http://" + acct_enode[0].Ip_port)
	// get gid by send createGroup
	groupRep, err := client.Call("smpc_createGroup", threshold, enodeList)
	if err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("smpc_createGroup = %s\n", groupRep))
	var groupJSON groupInfo
	groupData, _ := common.GetJSONData(groupRep)
	if err := json.Unmarshal(groupData, &groupJSON); err != nil {
		return nil, err
	}
	log.Info(fmt.Sprintf("\nGid = %s\n\n", groupJSON.Gid))

	tx, err := db.Conn.Begin()
	if err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}
	uid := uuid.New().String()
	sigs := ""
	for _, v := range acct_enode {
		_, err = db.BatchExecute("insert into accounts_info(threshold, gid , user_account, ip_port, enode, uuid) values(?,?,?,?,?,?)", tx,
			threshold, groupJSON.Gid, strings.ToLower(v.Account), v.Ip_port, v.Enode, uid)
		if err != nil {
			db.Conn.Rollback(tx)
			return nil, errors.New("internal db error " + err.Error())
		}
		sigs += common.StripEnode(v.Enode) + ":" + strings.ToLower(v.Account) + ":"
	}
	if err = db.Conn.Commit(tx); err != nil {
		return nil, errors.New("internal db error " + err.Error())
	}

	return Group{Gid: groupJSON.Gid, Sigs: strconv.Itoa(len(acct_enode)) + ":" + sigs[:len(sigs)-1], Uuid: uid}, nil
}
