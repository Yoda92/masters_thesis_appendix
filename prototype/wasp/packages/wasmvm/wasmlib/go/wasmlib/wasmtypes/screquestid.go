// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmtypes

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

const (
	ScRequestIDLength = 34
)

type ScRequestID struct {
	id [ScRequestIDLength]byte
}

func (o ScRequestID) Bytes() []byte {
	return RequestIDToBytes(o)
}

func (o ScRequestID) String() string {
	return RequestIDToString(o)
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

func RequestIDDecode(dec *WasmDecoder) ScRequestID {
	return requestIDFromBytesUnchecked(dec.FixedBytes(ScRequestIDLength))
}

func RequestIDEncode(enc *WasmEncoder, value ScRequestID) {
	enc.FixedBytes(value.id[:], ScRequestIDLength)
}

func RequestIDFromBytes(buf []byte) ScRequestID {
	if len(buf) == 0 {
		return ScRequestID{}
	}
	if len(buf) != ScRequestIDLength {
		panic("invalid RequestID length")
	}
	// final uint16 output index must be > ledgerstate.MaxOutputCount
	if buf[ScRequestIDLength-2] > 127 || buf[ScRequestIDLength-1] != 0 {
		panic("invalid RequestID: output index > 127")
	}
	return requestIDFromBytesUnchecked(buf)
}

func RequestIDToBytes(value ScRequestID) []byte {
	return value.id[:]
}

func RequestIDFromString(value string) ScRequestID {
	return RequestIDFromBytes(HexDecode(value))
}

func RequestIDToString(value ScRequestID) string {
	return HexEncode(value.id[:])
}

func requestIDFromBytesUnchecked(buf []byte) ScRequestID {
	o := ScRequestID{}
	copy(o.id[:], buf)
	return o
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScImmutableRequestID struct {
	proxy Proxy
}

func NewScImmutableRequestID(proxy Proxy) ScImmutableRequestID {
	return ScImmutableRequestID{proxy: proxy}
}

func (o ScImmutableRequestID) Exists() bool {
	return o.proxy.Exists()
}

func (o ScImmutableRequestID) String() string {
	return RequestIDToString(o.Value())
}

func (o ScImmutableRequestID) Value() ScRequestID {
	return RequestIDFromBytes(o.proxy.Get())
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScMutableRequestID struct {
	ScImmutableRequestID
}

func NewScMutableRequestID(proxy Proxy) ScMutableRequestID {
	return ScMutableRequestID{ScImmutableRequestID{proxy: proxy}}
}

func (o ScMutableRequestID) Delete() {
	o.proxy.Delete()
}

func (o ScMutableRequestID) SetValue(value ScRequestID) {
	o.proxy.Set(RequestIDToBytes(value))
}
