package apitypes

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/RWAValueRouter/ValueRouter/chains/tools/crypto"
)

type AuthToken struct {
	TypedData string `json:"typedData"`
	Signature string `json:"signature"`
	Address   string `json:"address"`
}

func GetTypedHash(src string) (string, error) {
	data := "{\"typedData\":\"" + base64.StdEncoding.EncodeToString([]byte(src)) + "\"}"
	var authToken AuthToken
	if err := json.Unmarshal([]byte(data), &authToken); err != nil {
		return "", fmt.Errorf("unmarshal auth token: %w", err)
	}

	typedDataBytes, err := base64.StdEncoding.DecodeString(authToken.TypedData)
	if err != nil {
		return "", fmt.Errorf("decode typed data: %w", err)
	}

	typedData := TypedData{}
	if err := json.Unmarshal(typedDataBytes, &typedData); err != nil {
		return "", fmt.Errorf("unmarshal typed data: %w", err)
	}

	// EIP-712 typed data marshalling
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return "", fmt.Errorf("eip712domain hash struct: %w", err)
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", fmt.Errorf("primary type hash struct: %w", err)
	}

	// add magic string prefix
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return hex.EncodeToString(crypto.Keccak256(rawData)), nil
}
