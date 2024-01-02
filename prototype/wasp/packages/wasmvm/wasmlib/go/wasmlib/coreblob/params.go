// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package coreblob

import (
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type MapStringToImmutableBytes struct {
	Proxy wasmtypes.Proxy
}

func (m MapStringToImmutableBytes) GetBytes(key string) wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(m.Proxy.Key(wasmtypes.StringToBytes(key)))
}

type ImmutableStoreBlobParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableStoreBlobParams() ImmutableStoreBlobParams {
	return ImmutableStoreBlobParams{Proxy: wasmlib.NewParamsProxy()}
}

// named chunks
func (s ImmutableStoreBlobParams) Blobs() MapStringToImmutableBytes {
	return MapStringToImmutableBytes(s)
}

// data schema for external tools
func (s ImmutableStoreBlobParams) DataSchema() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamDataSchema))
}

// smart contract program binary code
func (s ImmutableStoreBlobParams) ProgBinary() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamProgBinary))
}

// smart contract program source code
func (s ImmutableStoreBlobParams) Sources() wasmtypes.ScImmutableBytes {
	return wasmtypes.NewScImmutableBytes(s.Proxy.Root(ParamSources))
}

// VM type that must be used to run progBinary
func (s ImmutableStoreBlobParams) VMType() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(ParamVMType))
}

type MapStringToMutableBytes struct {
	Proxy wasmtypes.Proxy
}

func (m MapStringToMutableBytes) Clear() {
	m.Proxy.ClearMap()
}

func (m MapStringToMutableBytes) GetBytes(key string) wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(m.Proxy.Key(wasmtypes.StringToBytes(key)))
}

type MutableStoreBlobParams struct {
	Proxy wasmtypes.Proxy
}

// named chunks
func (s MutableStoreBlobParams) Blobs() MapStringToMutableBytes {
	return MapStringToMutableBytes(s)
}

// data schema for external tools
func (s MutableStoreBlobParams) DataSchema() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamDataSchema))
}

// smart contract program binary code
func (s MutableStoreBlobParams) ProgBinary() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamProgBinary))
}

// smart contract program source code
func (s MutableStoreBlobParams) Sources() wasmtypes.ScMutableBytes {
	return wasmtypes.NewScMutableBytes(s.Proxy.Root(ParamSources))
}

// VM type that must be used to run progBinary
func (s MutableStoreBlobParams) VMType() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(ParamVMType))
}

type ImmutableGetBlobFieldParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableGetBlobFieldParams() ImmutableGetBlobFieldParams {
	return ImmutableGetBlobFieldParams{Proxy: wasmlib.NewParamsProxy()}
}

// chunk name
func (s ImmutableGetBlobFieldParams) Field() wasmtypes.ScImmutableString {
	return wasmtypes.NewScImmutableString(s.Proxy.Root(ParamField))
}

// hash of the blob
func (s ImmutableGetBlobFieldParams) Hash() wasmtypes.ScImmutableHash {
	return wasmtypes.NewScImmutableHash(s.Proxy.Root(ParamHash))
}

type MutableGetBlobFieldParams struct {
	Proxy wasmtypes.Proxy
}

// chunk name
func (s MutableGetBlobFieldParams) Field() wasmtypes.ScMutableString {
	return wasmtypes.NewScMutableString(s.Proxy.Root(ParamField))
}

// hash of the blob
func (s MutableGetBlobFieldParams) Hash() wasmtypes.ScMutableHash {
	return wasmtypes.NewScMutableHash(s.Proxy.Root(ParamHash))
}

type ImmutableGetBlobInfoParams struct {
	Proxy wasmtypes.Proxy
}

func NewImmutableGetBlobInfoParams() ImmutableGetBlobInfoParams {
	return ImmutableGetBlobInfoParams{Proxy: wasmlib.NewParamsProxy()}
}

// hash of the blob
func (s ImmutableGetBlobInfoParams) Hash() wasmtypes.ScImmutableHash {
	return wasmtypes.NewScImmutableHash(s.Proxy.Root(ParamHash))
}

type MutableGetBlobInfoParams struct {
	Proxy wasmtypes.Proxy
}

// hash of the blob
func (s MutableGetBlobInfoParams) Hash() wasmtypes.ScMutableHash {
	return wasmtypes.NewScMutableHash(s.Proxy.Root(ParamHash))
}