package common

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anyswap/FastMulThreshold-DSA/crypto"
	"github.com/anyswap/FastMulThreshold-DSA/crypto/secp256k1"
	"github.com/anyswap/fastmpc-service-middleware/internal/common"
	"github.com/fsn-dev/cryptoCoins/coins"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StripAsset(asset string) map[string]interface{} {
	ret := make(map[string]interface{})
	if asset == "" {
		return ret
	}
	all := strings.Split(asset, "|")
	for _, a := range all {
		inner := strings.Split(a, "&")
		if len(inner) > 2 {
			ret[inner[1]] = inner[0]
		}
	}
	return ret
}

func IsSomeOneBlank(s ...string) bool {
	for _, v := range s {
		if IsBlank(v) {
			return true
		}
	}
	return false
}

func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	if len(strings.ReplaceAll(s, " ", "")) == 0 {
		return true
	}

	return false
}

// ConvertArrStrToStr convert array string into a string with | as separator
func ConvertArrStrToStr(strArr []string) string {
	ret := ""
	for i, v := range strArr {
		ret += v
		if i != len(strArr)-1 {
			ret += "|"
		}
	}
	return ret
}

// ValidateKeyId validate key id
func ValidateKeyId(keyId string) bool {
	if strings.HasPrefix(keyId, "0x") && len(keyId) == 66 {
		return true
	}
	if len(keyId) == 64 {
		return true
	}
	return false
}

// StripEnode get enode from reformed format
func StripEnode(enode string) string {
	s1 := strings.Split(enode, "//")
	if len(s1) != 2 {
		return ""
	}
	s2 := strings.Split(s1[1], "@")
	if len(s2) != 2 {
		return ""
	}
	if len(s2[0]) != 128 {
		return ""
	}
	return s2[0]
}

// PublicKeyBytesToAddress get address from public key
func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}

func SplitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}

// Map2Struct convert map into struct
// Field name must match
func Map2Struct(src map[string]interface{}, destStrct interface{}) {
	value := reflect.ValueOf(destStrct)
	e := value.Elem()
	for k, v := range src {
		f := e.FieldByName(strings.ToUpper(k[:1]) + k[1:])
		if !f.IsValid() {
			continue
		}
		if !f.CanSet() {
			continue
		}
		mv := reflect.ValueOf(v)
		// map value type
		mvt := mv.Kind().String()
		// struct field type
		sft := f.Kind().String()
		if sft != mvt {
			if mvt == "string" && (strings.Index(sft, "int") != -1) {
				if sft == "int64" {
					i, err := strconv.ParseInt(v.(string), 10, 64)
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "int32" {
					i, err := strconv.ParseInt(v.(string), 10, 32)
					r := int32(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				} else if sft == "int" {
					i, err := strconv.Atoi(v.(string))
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "uint64" {
					i, err := strconv.ParseUint(v.(string), 10, 64)
					if err == nil {
						f.Set(reflect.ValueOf(i))
					}
				} else if sft == "uint32" {
					i, err := strconv.ParseUint(v.(string), 10, 32)
					r := uint32(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				} else if sft == "uint" {
					i, err := strconv.ParseUint(v.(string), 10, 0)
					r := uint(i)
					if err == nil {
						f.Set(reflect.ValueOf(r))
					}
				}
			}

			// make string and string[] more friendly
			if mvt == "string" && sft == "slice" {
				_, ok := f.Interface().([]string)
				if ok {
					f.Set(reflect.ValueOf(strings.Split(v.(string), "|")))
				}
			}

			// make string and float more friendly
			if mvt == "string" && (strings.Index(sft, "float") != -1) {
				i, err := strconv.ParseFloat(v.(string), 64)
				if err == nil {
					f.Set(reflect.ValueOf(i))
				}
			}

			// make int to bool more friendly
			if mvt == "string" && sft == "bool" {
				i, err := strconv.Atoi(v.(string))
				if err == nil {
					if i == 1 {
						f.Set(reflect.ValueOf(true))
					} else if i == 0 {
						f.Set(reflect.ValueOf(false))
					}
				}
			}
			continue
		}
		f.Set(mv)
	}
}

// RecoverAddress recover address from data and signature
func RecoverAddress(data, sig string) (string, error) {
	hash := GetMsgSigHash([]byte(data))
	public, err := crypto.SigToPub(hash, common.FromHex(sig))
	if err != nil {
		return "", err
	}
	pub := secp256k1.S256("EC256K1").Marshal(public.X, public.Y)
	pubStr := hex.EncodeToString(pub)
	h := coins.NewCryptocoinHandler("ETH")
	if h == nil {
		return "", errors.New("h is zero")
	}
	addr, err := h.PublicKeyToAddress(pubStr)
	if err != nil {
		return "", err
	}
	return addr, nil
}

func CheckThreshold(threshold string) (int, int, error) {
	if !strings.Contains(threshold, "/") {
		return -1, -1, errors.New("invalid threshold")
	}
	parts := strings.Split(threshold, "/")
	if len(parts) != 2 {
		return -1, -1, errors.New("invalid threshold")
	}
	p1, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return -1, -1, errors.New("invalid threshold")
	}
	p2, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return -1, -1, errors.New("invalid threshold")
	}
	if p1 > p2 {
		return -1, -1, errors.New("invalid threshold")
	}
	if p1 < 1 || p2 < 1 {
		return -1, -1, errors.New("invalid threshold")
	}
	return int(p1), int(p2), nil
}

func CheckUserAccountsAndIpPortAddr(userAccountsAndIpPortAddr []string) ([]string, []string, error) {
	if len(userAccountsAndIpPortAddr) < 2 {
		return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
	}
	var accounts []string
	var ipPort []string
	accountsHolder := make(map[string]bool)
	ipPortHolder := make(map[string]bool)
	for _, v := range userAccountsAndIpPortAddr {
		// only contains ethereum address
		if !strings.Contains(v, "|") {
			if !CheckEthereumAddress(v) {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			} else {
				accounts = append(accounts, v)
				if accountsHolder[v] == true {
					return nil, nil, errors.New("duplicated account")
				} else {
					accountsHolder[v] = true
				}
				ipPort = append(ipPort, "")
			}
		} else {
			// contains ethereum address and ip:port
			parts := strings.Split(v, "|")
			if len(parts) != 2 {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			}
			if !CheckEthereumAddress(parts[0]) {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			} else {
				accounts = append(accounts, parts[0])
				if accountsHolder[parts[0]] == true {
					return nil, nil, errors.New("duplicated account")
				} else {
					accountsHolder[parts[0]] = true
				}
			}
			left := strings.TrimPrefix(parts[1], "https://")
			left = strings.TrimPrefix(left, "https://")
			if !strings.Contains(left, ":") {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			}
			p := strings.Split(left, ":")
			if len(p) != 2 {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			}
			if net.ParseIP(p[0]) == nil {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			}
			port, err := strconv.ParseInt(p[1], 10, 64)
			if err != nil {
				if net.ParseIP(p[0]) == nil {
					return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
				}
			}
			if !(0 < port && port < 65535) {
				return nil, nil, errors.New("invalid UserAccountsAndIpPortAddr")
			}
			ipPort = append(ipPort, left)
			if ipPortHolder[left] == true {
				return nil, nil, errors.New("duplicated ip port")
			} else {
				ipPortHolder[left] = true
			}
		}
	}
	return accounts, ipPort, nil
}

func CheckEthereumAddress(addr string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(addr)
}

func VerifyAccount(rsv string, msg string) error {
	sig := common.FromHex(rsv)
	if sig == nil {
		return errors.New("rsv from hex failed")
	}

	hash := GetMsgSigHash([]byte(msg))
	public, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return errors.New("SigToPub error " + err.Error())
	}
	type ReqData struct {
		Keytype string
		Account string
	}
	req := ReqData{}
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return err
	}
	pub := secp256k1.S256(req.Keytype).Marshal(public.X, public.Y)
	pubStr := hex.EncodeToString(pub)
	h := coins.NewCryptocoinHandler("ETH")
	if h == nil {
		return errors.New("NewCryptocoinHandler error")
	}
	addr, err := h.PublicKeyToAddress(pubStr)
	if err != nil {
		return errors.New("PublicKeyToAddress error " + err.Error())
	}
	if !strings.EqualFold(addr, req.Account) {
		return errors.New("verify sig fail")
	}
	return nil
}

func GetJSONData(successResponse json.RawMessage) ([]byte, error) {
	var rep response
	if err := json.Unmarshal(successResponse, &rep); err != nil {
		fmt.Println("getJSONData Unmarshal json fail:", err)
		return nil, err
	}
	if rep.Status != "Success" {
		return nil, errors.New(rep.Error)
	}
	repData, err := json.Marshal(rep.Data)
	if err != nil {
		fmt.Println("getJSONData Marshal json fail:", err)
		return nil, err
	}
	return repData, nil
}

// GetJSONResult parse result from rpc return data
func GetJSONResult(successResponse json.RawMessage) (string, error) {
	var data dataResult
	repData, err := GetJSONData(successResponse)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(repData, &data); err != nil {
		fmt.Println("getJSONResult Unmarshal json fail:", err)
		return "", err
	}
	return data.Result, nil
}

func GetRandomIndex(max int) int {
	return r.Intn(max)
}

var r *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
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
