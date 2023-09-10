package types

import (
	"github.com/anyswap/fastmpc-service-middleware/chains/common"
	"math/big"
	"testing"
)

func Test_getHash(t *testing.T) {
	to := common.HexToAddress("0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c")
	println(GetEVMMsgHash(0x5, 0, to, big.NewInt(0xb5e620f48000), 47984, big.NewInt(5338199), nil))
}

func GetEVMMsgHash(chainId int64, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) string {
	signer := NewEIP155Signer(big.NewInt(chainId))
	tx := NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	return signer.Hash(tx).String()
}
