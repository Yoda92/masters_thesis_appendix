// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmtypes

import (
	"encoding/binary"
	"strconv"
)

const ScInt32Length = 4

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

func Int32Decode(dec *WasmDecoder) int32 {
	return Int32FromBytes(dec.FixedBytes(ScInt32Length))
}

func Int32Encode(enc *WasmEncoder, value int32) {
	enc.FixedBytes(Int32ToBytes(value), ScInt32Length)
}

func Int32FromBytes(buf []byte) int32 {
	if len(buf) == 0 {
		return 0
	}
	if len(buf) != ScInt32Length {
		panic("invalid Int32 length")
	}
	return int32(binary.LittleEndian.Uint32(buf))
}

func Int32ToBytes(value int32) []byte {
	tmp := make([]byte, ScInt32Length)
	binary.LittleEndian.PutUint32(tmp, uint32(value))
	return tmp
}

func Int32FromString(value string) int32 {
	return int32(IntFromString(value, 32))
}

func Int32ToString(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScImmutableInt32 struct {
	proxy Proxy
}

func NewScImmutableInt32(proxy Proxy) ScImmutableInt32 {
	return ScImmutableInt32{proxy: proxy}
}

func (o ScImmutableInt32) Exists() bool {
	return o.proxy.Exists()
}

func (o ScImmutableInt32) String() string {
	return Int32ToString(o.Value())
}

func (o ScImmutableInt32) Value() int32 {
	return Int32FromBytes(o.proxy.Get())
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScMutableInt32 struct {
	ScImmutableInt32
}

func NewScMutableInt32(proxy Proxy) ScMutableInt32 {
	return ScMutableInt32{ScImmutableInt32{proxy: proxy}}
}

func (o ScMutableInt32) Delete() {
	o.proxy.Delete()
}

func (o ScMutableInt32) SetValue(value int32) {
	o.proxy.Set(Int32ToBytes(value))
}
