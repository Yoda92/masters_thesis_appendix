// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscmagic

import (
	"bytes"
	_ "embed"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv/codec"
)

type addressKind uint8

const (
	addressKindISCMagic = addressKind(iota)
	addressKindERC20BaseTokens
	addressKindERC20NativeTokens
	addressKindERC721NFTs
	addressKindERC721NFTCollection
	addressKindERC20ExternalNativeTokens
	addressKindInvalid
)

var (
	AddressPrefix = []byte{0x10, 0x74}
	Address       = packMagicAddress(addressKindISCMagic, nil)

	kindByteIndex    = len(AddressPrefix)
	headerLength     = len(AddressPrefix) + 1 // AddressPrefix + kind (byte)
	maxPayloadLength = common.AddressLength - headerLength
)

// ERC20NativeTokensAddress returns the Ethereum address of the ERC20 contract for
// native tokens with an on-chain foundry.
func ERC20NativeTokensAddress(foundrySN uint32) common.Address {
	return packMagicAddress(addressKindERC20NativeTokens, codec.EncodeUint32(foundrySN))
}

// ERC20ExternalNativeTokensAddress creates an Ethereum address for an ERC20 contract for
// native tokens with an off-chain foundry.
// NOTE: since len(NativeTokenID) == 38 and len(Address) == 20, it is not possible to assign
// a unique address to every native token. This function tries to form the ERC20 address
// from the first 17 bytes of hash(nativeTokenID). In case of a collision, it reapplies the
// hash and checks for a collision again, repeating until it gives up.
func ERC20ExternalNativeTokensAddress(
	nativeTokenID iotago.NativeTokenID,
	isTaken func(common.Address) bool,
) (common.Address, error) {
	const maxAttempts = 10
	hash := hashing.HashData(nativeTokenID[:])
	for i := 0; i < maxAttempts; i++ {
		addr := packMagicAddress(addressKindERC20ExternalNativeTokens, hash[:maxPayloadLength])
		if !isTaken(addr) {
			return addr, nil
		}
		hash = hashing.HashData(hash[:])
	}
	return common.Address{}, errors.New("all suitable addresses are taken")
}

func ERC20NativeTokensFoundrySN(addr common.Address) (uint32, error) {
	kind, payload, err := unpackMagicAddress(addr)
	if err != nil {
		return 0, err
	}
	if kind != addressKindERC20NativeTokens {
		return 0, errors.New("ERC20NativeTokensFoundrySN: invalid address kind")
	}
	if !allZero(payload[4:]) {
		return 0, errors.New("ERC20NativeTokensFoundrySN: invalid address format")
	}
	return codec.MustDecodeUint32(payload[0:4]), nil
}

func ERC721NFTCollectionAddress(collectionID iotago.NFTID) common.Address {
	return packMagicAddress(addressKindERC721NFTCollection, collectionID[:maxPayloadLength])
}

func packMagicAddress(kind addressKind, payload []byte) common.Address {
	var ret common.Address
	copy(ret[:], AddressPrefix)
	ret[kindByteIndex] = byte(kind)
	if len(payload) > maxPayloadLength {
		panic("packMagicAddress: invalid payload length")
	}
	copy(ret[headerLength:], payload)
	return ret
}

func unpackMagicAddress(addr common.Address) (addressKind, []byte, error) {
	if !bytes.Equal(addr[0:len(AddressPrefix)], AddressPrefix) {
		return 0, nil, errors.New("unpackMagicAddress: expected magic address prefix")
	}
	kind := addressKind(addr[kindByteIndex])
	if kind >= addressKindInvalid {
		return 0, nil, errors.New("unpackMagicAddress: unknown address kind")
	}
	payload := addr[headerLength:]
	return kind, payload, nil
}

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}
