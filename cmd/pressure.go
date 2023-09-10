package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/internal/common"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	common3 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RespAddr struct {
	Status          string `json:"Status"`
	User_account    string `json:"User_account"`
	Key_id          string `json:"Key_id"`
	Public_key      string `json:"Public_key"`
	Mpc_address     string `json:"Mpc_address"`
	Initializer     string `json:"Initializer"`
	Reply_status    string `json:"Reply_status"`
	Reply_timestamp string `json:"Reply_timestamp"`
	Reply_enode     string `json:"Reply_enode"`
	Gid             string `json:"Gid"`
	Threshold       string `json:"Threshold"`
	Mode            string `json:"Mode"`
	Key_type        string `json:"Key_type"`
}

type TxDataReqAddr struct {
	TxType        string
	Account       string
	Nonce         string
	Keytype       string
	GroupID       string
	ThresHold     string
	Mode          string
	FixedApprover []string
	AcceptTimeOut string
	TimeStamp     string
	Sigs          string
	Comment       string
	Uuid          string
}

var (
	parallel   *int
	passwd     *string
	passwdfile *string
	url        *string
	cmd        *string
	ts         *string
	sendtxrpc  *string
	keystores  arrayFlags
	client     *ethrpc.EthRPC
)

var (
	t           int
	addresses   []string
	keyWrappers []*keystore.Key
	keyGenlock  sync.Mutex
	sign_addr   []string
	sign_pub    []string
	sign_gid    []string
	keygen_succ = true
	faucet_priv *ecdsa.PrivateKey
	faucet      string
)

func main() {
	initFlag()
	switch *cmd {
	case "KEY_GEN":
		for {
			testKeyGen()
			time.Sleep(30 * time.Second)
		}
	case "SEND_TX":
		testSendTx()
	}
}

type response struct {
	Status string      `json:"Status"`
	Tip    string      `json:"Tip"`
	Error  string      `json:"Error"`
	Data   interface{} `json:"Data"`
}

type RespGroupId struct {
	Gid  string
	Sigs string
	Uuid string
}

func testSendTx() {
	testKeyGen()
	if !keygen_succ {
		log.Info("keygen not success, can not process tx")
		return
	}
	if *sendtxrpc == "" {
		log.Info("sendtxrpc is blank")
		return
	}
	log.Info("please send some native coin to address " + faucet + ",which will be used to testing, the waiting time is 2 minutes")
	recharged := false
	//TODO: sendingAmt set as feeAmt from cmd line
	sendingAmt := new(big.Int).SetInt64(1000000000000000)
	for i := 0; i < 24; i++ {
		balance, err := GetETHBalance()
		if err != nil {
			log.Info("get balance error")
			return
		}
		if balance.String() != "0" {
			println("charging eth balance " + balance.String())
			if balance.Cmp(new(big.Int).Mul(sendingAmt, new(big.Int).SetInt64(int64(*parallel)))) > 0 {
				println("receiving enough test coin")
				recharged = true
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
	if !recharged {
		println("2 minutes , not receiving enough test coin, exit for now")
		return
	}
	// TODO should get nonce from blockchain.
	var nonce uint64 = 0
	for _, addr := range sign_addr {
		txid, err := TransferETH(sendingAmt.Bytes(), nonce, *sendtxrpc, common3.HexToAddress(addr))
		if err != nil {
			log.Info("send test coin error " + err.Error())
			return
		}
		log.Info("sending test coin to address" + addr + ",txid :" + txid)
		nonce++
	}

}

func testKeyGen() {
	var kids []string
	succHolder := make(map[string]bool)
	log.Info("pressure test keygen")
	for i := 0; i < *parallel; i++ {
		go func() {
			successResponse, err := client.Call("smw_getGroupId", *ts, addresses[0:t])
			if err != nil {
				log.Info(err.Error())
				return
			}
			var rep response
			if err := json.Unmarshal(successResponse, &rep); err != nil {
				log.Info("getJSONData Unmarshal json fail:", err)
				return
			}
			if !strings.EqualFold(rep.Status, "Success") {
				log.Info("getJSONData Unmarshal json status error")
				return
			}
			repData, err := json.Marshal(rep.Data)
			if err != nil {
				log.Info("getJSONData Marshal json fail:", err)
				return
			}
			var r RespGroupId
			err = json.Unmarshal(repData, &r)
			if err != nil {
				log.Info(err.Error())
				return
			}
			req := &TxDataReqAddr{
				TxType:        "REQSMPCADDR",
				Account:       addresses[0],
				Nonce:         "3",
				Keytype:       "EC256K1",
				GroupID:       r.Gid,
				ThresHold:     *ts,
				Mode:          "2",
				AcceptTimeOut: "604800",
				TimeStamp:     new(big.Int).SetInt64(time.Now().UnixMilli()).String(),
				Sigs:          r.Sigs,
				Uuid:          r.Uuid,
			}

			payload, err := json.Marshal(req)
			if err != nil {
				panic(err)
			}

			rsv, err := signMsg(keyWrappers[0].PrivateKey, payload)
			successResponse, err = client.Call("smw_keyGen", rsv, string(payload))
			if err != nil {
				log.Info(err.Error())
				return
			}
			if err = json.Unmarshal(successResponse, &rep); err != nil {
				log.Info("getJSONData Unmarshal json fail:", err)
				return
			}
			if !strings.EqualFold(rep.Status, "Success") {
				log.Info("getJSONData Unmarshal json status error")
				return
			}
			log.Info("key_id:" + rep.Data.(string))
			keyGenlock.Lock()
			defer keyGenlock.Unlock()
			kids = append(kids, rep.Data.(string))
		}()
	}

	// 5 minutes not success consider invalid
	for m := 0; m < 40; m++ {
		log.Info("wait 30s to check keygen status")
		time.Sleep(30 * time.Second)
		var wait sync.WaitGroup
		for _, kid := range kids {
			wait.Add(1)
			go func(kid string) {
				defer wait.Done()
				if succHolder[kid] == true {
					return
				}
				successResponse, err := client.Call("smw_getReqAddrStatus", kid)
				if err != nil {
					log.Info(err.Error())
					return
				}
				var rep response
				if err := json.Unmarshal(successResponse, &rep); err != nil {
					log.Info("getJSONData Unmarshal json fail:", err)
					return
				}
				if !strings.EqualFold(rep.Status, "Success") {
					log.Info("getJSONData Unmarshal json status error")
					return
				}
				repData, err := json.Marshal(rep.Data)
				if err != nil {
					log.Info("getJSONData Marshal json fail:", err)
					return
				}
				var addrStatus []RespAddr
				err = json.Unmarshal(repData, &addrStatus)
				if err != nil {
					log.Info("getJSONData Marshal json fail:", err)
					return
				}
				var ok = true
				for _, v := range addrStatus {
					if v.Status == "0" {
						ok = false
						break
					}
				}
				if ok {
					keyGenlock.Lock()
					succHolder[kid] = true
					sign_addr = append(sign_addr, addrStatus[0].Mpc_address)
					sign_pub = append(sign_pub, addrStatus[0].Public_key)
					sign_gid = append(sign_gid, addrStatus[0].Gid)
					keyGenlock.Unlock()
				}
			}(kid)
		}
		wait.Wait()
		if len(kids) == len(succHolder) {
			break
		}
	}
	for _, v := range kids {
		var out string
		out = "kid:" + v
		if succHolder[v] {
			out += ",status success"
		} else {
			keygen_succ = false
			out += ",status not success"
		}
		log.Info(out)
	}
}

func initFlag() {
	parallel = flag.Int("parallel", 1, "parallel count")
	passwd = flag.String("passwd", "", "Password")
	passwdfile = flag.String("passwdfile", "", "Password file")
	url = flag.String("url", "https://api.smpcwallet.com", "smpc api url")
	cmd = flag.String("cmd", "", "KEY_GEN|SIGN")
	ts = flag.String("ts", "2/2", "Threshold")
	sendtxrpc = flag.String("sendtxrpc", "", "send tx chain rpc")

	// array
	flag.Var(&keystores, "keystore", "Keystore file")

	flag.Parse()
	for _, keysto := range keystores {
		keyjson, err := ioutil.ReadFile(keysto)
		if err != nil {
			log.Info("Read keystore fail", err)
			panic(err)
		}
		if *passwd == "" && *passwdfile != "" {
			tmpPas, err := ioutil.ReadFile(*passwdfile)
			if err != nil {
				log.Info(err.Error())
				return
			}
			*passwd = string(tmpPas)
		}
		keyWrapper, _ := keystore.DecryptKey(keyjson, *passwd)
		addresses = append(addresses, keyWrapper.Address.String())
		keyWrappers = append(keyWrappers, keyWrapper)
		if faucet == "" {
			faucet_priv = keyWrapper.PrivateKey
			faucet = keyWrapper.Address.String()
		}
	}

	// init RPC client
	client = ethrpc.New(*url, ethrpc.WithHttpClient(&http.Client{Timeout: 5 * time.Second}))

	if len(*ts) != 3 {
		log.Info("invalid ts")
		return
	}
	var err error
	t, err = strconv.Atoi((*ts)[2:3])
	if err != nil {
		log.Info(err.Error())
		return
	}
	if t > len(addresses) {
		log.Info("max threshold bigger than keystores")
		return
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return fmt.Sprint(*i)
}
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// signMsg sign msg
func signMsg(privatekey *ecdsa.PrivateKey, playload []byte) (string, error) {
	// sign tx by privatekey
	hash := GetMsgSigHash(playload)
	//hash := crypto.Keccak256([]byte(header),playload)
	signature, signatureErr := crypto.Sign(hash, privatekey)
	if signatureErr != nil {
		log.Info("signature create error")
		panic(signatureErr)
	}
	rsv := common.ToHex(signature)
	return rsv, nil
}

func GetMsgSigHash(message []byte) []byte {
	msglen := []byte(strconv.Itoa(len(message)))

	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte{0x19})
	hash.Write([]byte("Ethereum Signed Message:"))
	hash.Write([]byte{0x0A})
	hash.Write(msglen)
	hash.Write(message)
	buf := hash.Sum([]byte{})
	return buf
}

func TransferETH(sendingAmt []byte, nonce uint64, rpc string, toAddress common3.Address) (string, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return "", err
	}
	privateKey := faucet_priv
	value := new(big.Int).SetBytes(sendingAmt) // in wei (0 eth)
	log.Info("sending amt " + value.String())
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	log.Info("gasPrice" + gasPrice.String())
	//gasLimit = gasLimit * 2
	gasLimit := uint64(1000000) // in units
	//log.Info(gasLimit) // 23256

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}
	txid := signedTx.Hash().Hex()
	return txid, nil
}

func GetETHBalance() (*big.Int, error) {
	client, err := ethclient.Dial(*sendtxrpc)
	if err != nil {
		return nil, err
	}
	b, err := client.BalanceAt(context.Background(), common3.HexToAddress(faucet), nil)
	if err != nil {
		return nil, err
	}

	return b, nil
}
