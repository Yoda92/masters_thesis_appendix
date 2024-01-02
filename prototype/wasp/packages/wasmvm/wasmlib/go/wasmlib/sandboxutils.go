// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmlib

import (
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type ScSandboxUtils struct{}

// Bech32Decode decodes the specified bech32-encoded string value to its original address
func (u ScSandboxUtils) Bech32Decode(value string) wasmtypes.ScAddress {
	return wasmtypes.AddressFromBytes(Sandbox(FnUtilsBech32Decode, wasmtypes.StringToBytes(value)))
}

// Bech32Encode encodes the specified address to a bech32-encoded string
func (u ScSandboxUtils) Bech32Encode(addr wasmtypes.ScAddress) string {
	return wasmtypes.StringFromBytes(Sandbox(FnUtilsBech32Encode, wasmtypes.AddressToBytes(addr)))
}

func (u ScSandboxUtils) BlsAddressFromPubKey(pubKey []byte) wasmtypes.ScAddress {
	return wasmtypes.AddressFromBytes(Sandbox(FnUtilsBlsAddress, pubKey))
}

func (u ScSandboxUtils) BlsAggregateSignatures(pubKeys, sigs [][]byte) ([]byte, []byte) {
	enc := wasmtypes.NewWasmEncoder()
	wasmtypes.Uint32Encode(enc, uint32(len(pubKeys)))
	for _, pubKey := range pubKeys {
		enc.Bytes(pubKey)
	}
	wasmtypes.Uint32Encode(enc, uint32(len(sigs)))
	for _, sig := range sigs {
		enc.Bytes(sig)
	}
	result := Sandbox(FnUtilsBlsAggregate, enc.Buf())
	decode := wasmtypes.NewWasmDecoder(result)
	return decode.Bytes(), decode.Bytes()
}

func (u ScSandboxUtils) BlsValidSignature(data, pubKey, signature []byte) bool {
	enc := wasmtypes.NewWasmEncoder().Bytes(data).Bytes(pubKey).Bytes(signature)
	return wasmtypes.BoolFromBytes(Sandbox(FnUtilsBlsValid, enc.Buf()))
}

func (u ScSandboxUtils) Ed25519AddressFromPubKey(pubKey []byte) wasmtypes.ScAddress {
	return wasmtypes.AddressFromBytes(Sandbox(FnUtilsEd25519Address, pubKey))
}

func (u ScSandboxUtils) Ed25519ValidSignature(data, pubKey, signature []byte) bool {
	enc := wasmtypes.NewWasmEncoder().Bytes(data).Bytes(pubKey).Bytes(signature)
	return wasmtypes.BoolFromBytes(Sandbox(FnUtilsEd25519Valid, enc.Buf()))
}

// hashes the specified value bytes using Blake2b hashing and returns the resulting 32-byte hash
func (u ScSandboxUtils) HashBlake2b(value []byte) wasmtypes.ScHash {
	return wasmtypes.HashFromBytes(Sandbox(FnUtilsHashBlake2b, value))
}

// hashes the specified value bytes using Keccak hashing and returns the resulting 32-byte hash
func (u ScSandboxUtils) HashKeccak(value []byte) wasmtypes.ScHash {
	return wasmtypes.HashFromBytes(Sandbox(FnUtilsHashKeccak, value))
}

// hashes the specified value bytes using Blake2b hashing and returns the resulting 32-byte hash
func (u ScSandboxUtils) HashName(value string) wasmtypes.ScHname {
	return wasmtypes.HnameFromBytes(Sandbox(FnUtilsHashName, []byte(value)))
}

// hashes the specified value bytes using SHA3 hashing and returns the resulting 32-byte hash
func (u ScSandboxUtils) HashSha3(value []byte) wasmtypes.ScHash {
	return wasmtypes.HashFromBytes(Sandbox(FnUtilsHashSha3, value))
}

// converts an integer to its string representation
func (u ScSandboxUtils) String(value int64) string {
	return wasmtypes.Int64ToString(value)
}
