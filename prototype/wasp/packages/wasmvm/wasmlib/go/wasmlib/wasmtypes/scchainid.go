// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmtypes

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

const ScChainIDLength = 32

type ScChainID struct {
	id [ScChainIDLength]byte
}

// Address returns the alias address that the chain ID actually represents
func (o ScChainID) Address() ScAddress {
	buf := []byte{ScAddressAlias}
	return AddressFromBytes(append(buf, o.id[:]...))
}

func (o ScChainID) Bytes() []byte {
	return ChainIDToBytes(o)
}

func (o ScChainID) String() string {
	return ChainIDToString(o)
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

func ChainIDDecode(dec *WasmDecoder) ScChainID {
	o := ScChainID{}
	copy(o.id[:], dec.FixedBytes(ScChainIDLength))
	return o
}

func ChainIDEncode(enc *WasmEncoder, value ScChainID) {
	enc.FixedBytes(value.id[:], ScChainIDLength)
}

func ChainIDFromBytes(buf []byte) ScChainID {
	o := ScChainID{}
	if len(buf) == 0 {
		return o
	}
	if len(buf) != ScChainIDLength {
		panic("invalid ChainID length")
	}
	copy(o.id[:], buf)
	return o
}

func ChainIDToBytes(value ScChainID) []byte {
	return value.id[:]
}

func ChainIDFromString(value string) ScChainID {
	addr := AddressFromString(value)
	if addr.id[0] != ScAddressAlias {
		panic("invalid ChainID address type")
	}
	return ChainIDFromBytes(addr.id[1:])
}

func ChainIDToString(value ScChainID) string {
	return AddressToString(value.Address())
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScImmutableChainID struct {
	proxy Proxy
}

func NewScImmutableChainID(proxy Proxy) ScImmutableChainID {
	return ScImmutableChainID{proxy: proxy}
}

func (o ScImmutableChainID) Exists() bool {
	return o.proxy.Exists()
}

func (o ScImmutableChainID) String() string {
	return ChainIDToString(o.Value())
}

func (o ScImmutableChainID) Value() ScChainID {
	return ChainIDFromBytes(o.proxy.Get())
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type ScMutableChainID struct {
	ScImmutableChainID
}

func NewScMutableChainID(proxy Proxy) ScMutableChainID {
	return ScMutableChainID{ScImmutableChainID{proxy: proxy}}
}

func (o ScMutableChainID) Delete() {
	o.proxy.Delete()
}

func (o ScMutableChainID) SetValue(value ScChainID) {
	o.proxy.Set(ChainIDToBytes(value))
}
