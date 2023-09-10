package jobs

import (
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/RWAValueRouter/ValueRouter/chains/common"
	"github.com/RWAValueRouter/ValueRouter/chains/tools/crypto"
	"github.com/RWAValueRouter/ValueRouter/chains/types"
	"github.com/RWAValueRouter/ValueRouter/db"
	"strings"
)

var (
	failedTimeCounter = make(map[string]int)
	MAX_FAIL          = 10
)

type ToBeSend struct {
	Rsv         string
	Chain_type  int
	Key_id      string
	Msg_context string
	Mpc_address string
}

func listenTransactions() {
	log.Info("listenTransactions")
	l, err := db.Conn.GetStructValue("select rsv, chain_type , key_id, msg_context, mpc_address from signs_detail where status = 1 and rsv is not null and sign_type = 0 group by rsv, chain_type , key_id, msg_context, mpc_address", ToBeSend{})
	if err != nil {
		log.Error("internal db error " + err.Error())
		return
	}
	for _, v := range l {
		tbs := v.(*ToBeSend)
		go fireTx(tbs)
	}
}

func fireTx(send *ToBeSend) {
	switch types.ChainType(send.Chain_type) {
	case types.EVM:
		go doEVMTx(send)
	default:
		log.Error("unrecognized chain type")
		return
	}
}

func doEVMTx(send *ToBeSend) {
	if len(CC) == 0 {
		log.Info("initial not finished")
		return
	}
	log.Info("doEVMTx", "Msg_context", send.Msg_context, " key_id ", send.Key_id)
	tx, err := types.EvmChain.GetUnsignedTransaction(send.Msg_context)
	if err != nil {
		log.Error("fireTx", "msg", err.Error())
		return
	}
	if len(strings.Split(send.Rsv, "|")) != 1 {
		log.Error("rsv len must be 1")
		return
	}
	signature := common.FromHex(send.Rsv)
	if len(signature) != crypto.SignatureLength {
		log.Error("DcrmSignTransaction wrong length of signature")
		return
	}
	signer, err := types.EvmChain.GetSigner(send.Msg_context)
	if err != nil {
		log.Error("get wrong signer " + err.Error())
		return
	}
	signedTx, err := types.EvmChain.SignTxWithSignature(signer, tx, signature, common.HexToAddress(send.Mpc_address))
	if err != nil {
		log.Error("SignTxWithSignature", "msg", err.Error())
		return
	}
	hash := tx.Hash()
	if hash == common.EmptyHash {
		log.Error("empty hash")
		return
	}
	data, err := signedTx.MarshalBinary()
	if err != nil {
		log.Error(err.Error())
		return
	}
	log.Info("call eth_sendRawTransaction start", "txHash", signedTx.Hash().String())
	Cm.RLock()
	defer Cm.RUnlock()
	hexData := common.ToHex(data)
	config := CC[send.Chain_type]
	var cf *ChainConfig
	var txid string
	if config != nil {
		if send.Chain_type == 0 {
			chainID, err := types.EvmChain.GetChainId(send.Msg_context)
			if err != nil {
				log.Error(err.Error())
				return
			}
			cf = config[int(chainID)]
		} else {
			cf = config[0]
		}
		if cf != nil {
			switch send.Chain_type {
			case int(types.EVM):
				serverDown := false
				for _, rpc := range cf.Rpc_list {
					txid, err = types.EvmChain.SendRawTransaction(hexData, rpc)
					if err != nil {
						log.Error("send tx transaction " + err.Error())
						continue
					}
					break
				}
				if err != nil {
					if strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Client.Timeout") {
						serverDown = true
					}
				}
				if txid == "" && serverDown == false {
					if failedTime, ok := failedTimeCounter[txid]; ok {
						failedTimeCounter[txid] = failedTime + 1
					} else {
						failedTimeCounter[txid] = 1
					}
					if failedTimeCounter[txid] > MAX_FAIL {
						_, err = db.Conn.CommitOneRow("update signs_detail set status = 7 where key_id = ?", send.Key_id)
						if err != nil {
							log.Error("internal db error" + err.Error())
							return
						}
						return
					}
					log.Error("error send tx will try latter")
					return
				}
			default:
				log.Error("not support chain type")
				return
			}
		} else {
			log.Error("can not find config")
			return
		}
	} else {
		log.Error("unrecognized chain type")
		return
	}
	if txid != "" {
		_, err = db.Conn.CommitOneRow("update signs_detail set txid = ?, status = 4 where key_id = ?", txid, send.Key_id)
		if err != nil {
			log.Error("internal db error" + err.Error())
			return
		}
	}
}

func init() {
	jobs.AddFunc("@every 10s", listenTransactions)
}
