package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEVMChain_GetUnsignedTransactionHash(t *testing.T) {
	str := "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"chainId\":\"0x5\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}"
	c := NewEVMChain()
	hash, err := c.GetUnsignedTransactionHash(str)
	if assert.NoError(t, err) {
		assert.Equal(t, "0x8dd93ccbbc3ed6b2d4ae7976c73e8b35e1a695aae91ca41211a8746892703677", hash)
	}
}

func TestEVMChain_ValidateUnsignedTransactionHash(t *testing.T) {
	str := "{\"from\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"to\":\"0xd6a643ECE3E3d21Af05D7DA3875860B2EFcA164c\",\"chainId\":\"0x5\",\"value\":\"0xb5e620f48000\",\"nonce\":0,\"gas\":47984,\"gasPrice\":5338199,\"data\":\"\",\"originValue\":\"0.0002\",\"name\":\"ETH Goerli\"}"
	c := NewEVMChain()
	hash := "0x8dd93ccbbc3ed6b2d4ae7976c73e8b35e1a695aae91ca41211a8746892703677"
	assert.True(t, c.ValidateUnsignedTransactionHash(str, hash))
}
