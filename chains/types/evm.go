package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/anyswap/FastMulThreshold-DSA/log"
	"github.com/anyswap/fastmpc-service-middleware/chains/client"
	common2 "github.com/anyswap/fastmpc-service-middleware/chains/common"
	"github.com/anyswap/fastmpc-service-middleware/chains/signer/core/apitypes"
	"github.com/anyswap/fastmpc-service-middleware/chains/tools/crypto"
	"github.com/anyswap/fastmpc-service-middleware/common"
	"math/big"
	"strconv"
	"strings"
)

var (
	errEmptyURLs              = errors.New("empty URLs")
	errTxInOrphanBlock        = errors.New("tx is in orphan block")
	errTxHashMismatch         = errors.New("tx hash mismatch with rpc result")
	errTxBlockHashMismatch    = errors.New("tx block hash mismatch with rpc result")
	errTxReceiptMissBlockInfo = errors.New("tx receipt missing block info")
)

type EVMChain struct {
	ChainType
}

var EvmChain = NewEVMChain()

func NewEVMChain() *EVMChain {
	return &EVMChain{
		EVM,
	}
}

type UnsignedEVMTx struct {
	From        string `json:"from"`
	To          string `json:"to"`
	ChainId     string `json:"chainId"`
	Value       string `json:"value"`
	Nonce       uint64 `json:"nonce"`
	Gas         uint64 `json:"gas"`
	GasPrice    int64  `json:"gasPrice"`
	Data        string `json:"data"`
	OriginValue string `json:"originValue"`
	Name        string `json:"name"`
}

func (*EVMChain) ValidateParam(txdata string) bool {
	unTx := &UnsignedEVMTx{}
	err := json.Unmarshal([]byte(txdata), unTx)
	if err != nil {
		return false
	}
	if !common.CheckEthereumAddress(unTx.From) {
		return false
	}
	if !common.CheckEthereumAddress(unTx.To) {
		return false
	}
	if common.IsBlank(unTx.ChainId) || common.IsBlank(unTx.Name) || common.IsBlank(unTx.Value) ||
		common.IsBlank(unTx.OriginValue) || unTx.Gas == 0 || unTx.GasPrice == 0 {
		return false
	}
	if len(unTx.Value) < 2 || len(unTx.ChainId) < 2 {
		return false
	}
	_, ok := new(big.Int).SetString(unTx.ChainId[2:], 16)
	if !ok {
		return false
	}
	_, ok = new(big.Int).SetString(unTx.Value[2:], 16)
	if !ok {
		return false
	}

	return true
}

func (e *EVMChain) GetPersonalSignHash(message string) (string, error) {
	str := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message
	return hex.EncodeToString(crypto.Keccak256([]byte(str))), nil
}

func (e *EVMChain) GetETHSignHash(message string) (string, error) {
	return hex.EncodeToString(crypto.Keccak256([]byte(message))), nil
}

func (e *EVMChain) GetTypedHash(message string) (string, error) {
	hash, err := apitypes.GetTypedHash(message)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func (e *EVMChain) GetUnsignedTransactionHash(data string) (string, error) {
	if !e.ValidateParam(data) {
		return "", errors.New("invalid unsigned tx data")
	}
	unTx := &UnsignedEVMTx{}
	_ = json.Unmarshal([]byte(data), unTx)
	chainId, _ := new(big.Int).SetString(unTx.ChainId[2:], 16)
	signer := NewEIP155Signer(chainId)
	amount, _ := new(big.Int).SetString(unTx.Value[2:], 16)
	var payload []byte
	var err error
	if unTx.Data == "" || strings.EqualFold(unTx.Data, "0x") {
		payload = []byte{}
	} else {
		if strings.HasPrefix(unTx.Data, "0x") || strings.HasPrefix(unTx.Data, "0X") {
			payload, err = hex.DecodeString(unTx.Data[2:])
			if err != nil {
				return "", err
			}
		} else {
			payload, err = hex.DecodeString(unTx.Data)
			if err != nil {
				return "", err
			}
		}
	}
	tx := NewTransaction(unTx.Nonce, common2.HexToAddress(unTx.To), amount, unTx.Gas, new(big.Int).SetInt64(unTx.GasPrice), payload)
	return signer.Hash(tx).String(), nil
}

func (e *EVMChain) ValidatePersonalSignHash(src string, cmp string) bool {
	type Msg struct {
		Data    string
		ChainId string
	}
	msg := &Msg{}
	err := json.Unmarshal([]byte(src), msg)
	if err != nil {
		log.Warn("ValidatePersonalSignHash", "msg", "not valid msg"+err.Error())
		return false
	}
	hash, err := e.GetPersonalSignHash(msg.Data)
	if err != nil {
		log.Error("ValidateUnsignedTransactionHash", "hash", err)
		return false
	}
	if !strings.EqualFold(hash, cmp) {
		return false
	}
	return true
}

func (e *EVMChain) ValidateETHSignHash(src string, cmp string) bool {
	type Msg struct {
		Data    string
		ChainId string
	}
	msg := &Msg{}
	err := json.Unmarshal([]byte(src), msg)
	if err != nil {
		log.Warn("ValidateETHSignHash", "msg", "not valid msg"+err.Error())
		return false
	}
	hash, err := e.GetETHSignHash(msg.Data)
	if err != nil {
		log.Error("ValidateUnsignedTransactionHash", "hash", err)
		return false
	}
	if !strings.EqualFold(hash, cmp) {
		return false
	}
	return true
}

func (e *EVMChain) ValidateTypedHash(src string, cmp string) bool {
	type Msg struct {
		Data    string
		ChainId string
	}
	msg := &Msg{}
	err := json.Unmarshal([]byte(src), msg)
	if err != nil {
		log.Warn("ValidateTypedHash", "msg", "not valid msg"+err.Error())
		return false
	}
	hash, err := e.GetTypedHash(msg.Data)
	if err != nil {
		log.Error("ValidateTypedHash", "hash", err)
		return false
	}
	if !strings.EqualFold(hash, cmp) {
		return false
	}
	return true
}

func (e *EVMChain) ValidateUnsignedTransactionHash(src string, cmp string) bool {
	if !e.ValidateParam(src) {
		log.Warn("ValidateUnsignedTransactionHash", "src", "not valid src")
		return false
	}
	hash, err := e.GetUnsignedTransactionHash(src)
	if err != nil {
		log.Error("ValidateUnsignedTransactionHash", "hash", err)
		return false
	}
	if !strings.EqualFold(hash, cmp) {
		return false
	}
	return true
}

func (e *EVMChain) GetUnsignedTransaction(data string) (*Transaction, error) {
	if !e.ValidateParam(data) {
		return nil, errors.New("invalid unsigned tx data")
	}
	unTx := &UnsignedEVMTx{}
	_ = json.Unmarshal([]byte(data), unTx)
	amount, _ := new(big.Int).SetString(unTx.Value[2:], 16)
	var payload []byte
	var err error
	if unTx.Data == "" || strings.EqualFold(unTx.Data, "0x") {
		payload = []byte{}
	} else {
		if strings.HasPrefix(unTx.Data, "0x") || strings.HasPrefix(unTx.Data, "0X") {
			payload, err = hex.DecodeString(unTx.Data[2:])
			if err != nil {
				return nil, err
			}
		} else {
			payload, err = hex.DecodeString(unTx.Data)
			if err != nil {
				return nil, err
			}
		}
	}
	tx := NewTransaction(unTx.Nonce, common2.HexToAddress(unTx.To), amount, unTx.Gas, new(big.Int).SetInt64(unTx.GasPrice), payload)
	return tx, nil
}

func (e *EVMChain) SignTxWithSignature(signer Signer, tx *Transaction, signature []byte, signerAddr common2.Address) (*Transaction, error) {
	vPos := crypto.SignatureLength - 1
	for i := 0; i < 2; i++ {
		signedTx, err := tx.WithSignature(signer, signature)
		if err != nil {
			return nil, err
		}

		sender, err := Sender(signer, signedTx)
		if err != nil {
			return nil, err
		}

		if sender == signerAddr {
			return signedTx, nil
		}

		signature[vPos] ^= 0x1 // v can only be 0 or 1
	}

	return nil, errors.New("wrong sender address")
}

func (e *EVMChain) GetSigner(data string) (Signer, error) {
	if !e.ValidateParam(data) {
		return nil, errors.New("invalid unsigned tx data")
	}
	unTx := &UnsignedEVMTx{}
	_ = json.Unmarshal([]byte(data), unTx)
	chainId, _ := new(big.Int).SetString(unTx.ChainId[2:], 16)
	s := NewEIP155Signer(chainId)
	return &s, nil
}

func (e *EVMChain) GetChainId(data string) (int64, error) {
	if !e.ValidateParam(data) {
		return -1, errors.New("invalid unsigned tx data")
	}
	unTx := &UnsignedEVMTx{}
	_ = json.Unmarshal([]byte(data), unTx)
	chainId, _ := new(big.Int).SetString(unTx.ChainId[2:], 16)
	return chainId.Int64(), nil
}

func (e *EVMChain) SendRawTransaction(hexData string, url string) (string, error) {
	var result string
	err := client.RPCPost(&result, url, "eth_sendRawTransaction", hexData)
	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func (e *EVMChain) GetTransactionReceipt(txHash string, url string) (result *RPCTxReceipt, err error) {
	err = client.RPCPost(&result, url, "eth_getTransactionReceipt", txHash)
	if err == nil && result != nil {
		if result.BlockNumber == nil || result.BlockHash == nil || result.TxIndex == nil {
			return nil, errTxReceiptMissBlockInfo
		}
		if !strings.EqualFold(result.TxHash.Hex(), txHash) {
			return nil, errTxHashMismatch
		}
		return result, nil
	}
	return nil, err
}

func (e *EVMChain) GetTransactionByHash(txHash string, url string) (result *RPCTransaction, err error) {
	err = client.RPCPost(&result, url, "eth_getTransactionByHash", txHash)
	if err == nil && result != nil {
		if !strings.EqualFold(result.Hash.Hex(), txHash) {
			return nil, errTxHashMismatch
		}
		return result, nil
	}
	return nil, err
}
