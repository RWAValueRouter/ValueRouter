// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package apitypes

import (
	"testing"
)

func TestIsPrimitive(t *testing.T) {
	// Expected positives
	for i, tc := range []string{
		"int24", "int24[]", "uint88", "uint88[]", "uint", "uint[]", "int256", "int256[]",
		"uint96", "uint96[]", "int96", "int96[]", "bytes17[]", "bytes17",
	} {
		if !isPrimitiveTypeValid(tc) {
			t.Errorf("test %d: expected '%v' to be a valid primitive", i, tc)
		}
	}
	// Expected negatives
	for i, tc := range []string{
		"int257", "int257[]", "uint88 ", "uint88 []", "uint257", "uint-1[]",
		"uint0", "uint0[]", "int95", "int95[]", "uint1", "uint1[]", "bytes33[]", "bytess",
	} {
		if isPrimitiveTypeValid(tc) {
			t.Errorf("test %d: expected '%v' to not be a valid primitive", i, tc)
		}
	}

	typedData := "{\"types\":{\"EIP712Domain\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"version\",\"type\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\"}],\"RelayRequest\":[{\"name\":\"target\",\"type\":\"address\"},{\"name\":\"encodedFunction\",\"type\":\"bytes\"},{\"name\":\"gasData\",\"type\":\"GasData\"},{\"name\":\"relayData\",\"type\":\"RelayData\"}],\"GasData\":[{\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"name\":\"pctRelayFee\",\"type\":\"uint256\"},{\"name\":\"baseRelayFee\",\"type\":\"uint256\"}],\"RelayData\":[{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"senderNonce\",\"type\":\"uint256\"},{\"name\":\"relayWorker\",\"type\":\"address\"},{\"name\":\"paymaster\",\"type\":\"address\"}]},\"domain\":{\"name\":\"GSN Relayed Transaction\",\"version\":\"1\",\"chainId\":42,\"verifyingContract\":\"0x6453D37248Ab2C16eBd1A8f782a2CBC65860E60B\"},\"primaryType\":\"RelayRequest\",\"message\":{\"target\":\"0x9cf40ef3d1622efe270fe6fe720585b4be4eeeff\",\"encodedFunction\":\"0xa9059cbb0000000000000000000000002e0d94754b348d208d64d52d78bcd443afa9fa520000000000000000000000000000000000000000000000000000000000000007\",\"gasData\":{\"gasLimit\":\"39507\",\"gasPrice\":\"1700000000\",\"pctRelayFee\":\"70\",\"baseRelayFee\":\"0\"},\"relayData\":{\"senderAddress\":\"0x22d491bde2303f2f43325b2108d26f1eaba1e32b\",\"senderNonce\":\"3\",\"relayWorker\":\"0x3baee457ad824c94bd3953183d725847d023a2cf\",\"paymaster\":\"0x957F270d45e9Ceca5c5af2b49f1b5dC1Abb0421c\"}}}"
	v, err := GetTypedHash(typedData)
	if err != nil {
		t.Errorf("expected nil , but got error %s", err.Error())
	}

	if v != "b21808615920f4a43f5da837cdba41d2859694b4d197e6d33ab93e7eb1b9f10e" {
		t.Errorf("expected get b21808615920f4a43f5da837cdba41d2859694b4d197e6d33ab93e7eb1b9f10e , but got %s", v)
	}

}
