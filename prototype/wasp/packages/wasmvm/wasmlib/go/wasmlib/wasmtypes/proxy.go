// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmtypes

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type IKvStore interface {
	Bytes() []byte
	Delete(key []byte)
	Exists(key []byte) bool
	Get(key []byte) []byte
	Set(key, value []byte)
}

type Proxy struct {
	key     []byte
	kvStore IKvStore
}

func NewProxy(kvStore IKvStore) Proxy {
	return Proxy{kvStore: kvStore}
}

// Append returns a Proxy for a newly appended null element
// Note that this will essentially return the element at Length()
func (p Proxy) Append() Proxy {
	length := p.Length()
	p.expand(length + 1)
	return p.element(length)
}

func (p Proxy) Bytes() []byte {
	return p.kvStore.Bytes()
}

// ClearArray clears an array by deleting all elements
// TODO Note that this does not delete recursive container elements
func (p Proxy) ClearArray() {
	for length := p.Length(); length != 0; length-- {
		p.element(length - 1).Delete()
	}

	// clear the length counter
	p.Delete()
}

// ClearMap clears a map by deleting all elements
// TODO Note that this does not delete recursive container elements
func (p Proxy) ClearMap() {
	// TODO clearPrefix

	// clear the length counter
	p.Delete()
}

func (p Proxy) Delete() {
	p.kvStore.Delete(p.key)
}

func (p Proxy) element(index uint32) Proxy {
	enc := NewWasmEncoder()
	enc.VluEncode(uint64(index))
	return p.sub('#', enc.Buf())
}

func (p Proxy) Exists() bool {
	return p.kvStore.Exists(p.key)
}

// TODO have a Grow function that grows an array?
func (p Proxy) expand(length uint32) {
	// update the length counter
	enc := NewWasmEncoder()
	enc.VluEncode(uint64(length))
	p.Set(enc.Buf())
}

func (p Proxy) Get() []byte {
	return p.kvStore.Get(p.key)
}

// Index gets a Proxy for an element of an Array by its index
func (p Proxy) Index(index uint32) Proxy {
	size := p.Length()
	if index >= size {
		if index == size {
			panic("invalid index: use append")
		}
		panic("invalid index")
	}
	return p.element(index)
}

// Key gets a Proxy for an element of a Map by its key
func (p Proxy) Key(key []byte) Proxy {
	return p.sub('.', key)
}

// Length returns the number of elements in an Array
// Never try to access an index >= Length()
func (p Proxy) Length() uint32 {
	// get the length counter
	buf := p.Get()
	if len(buf) == 0 {
		return 0
	}
	dec := NewWasmDecoder(buf)
	return uint32(dec.VluDecode(32))
}

// Root returns a Proxy for an element of a root container (Params/Results/State).
// The key is always a string.
func (p Proxy) Root(key string) Proxy {
	return Proxy{kvStore: p.kvStore, key: []byte(key)}
}

func (p Proxy) Set(value []byte) {
	p.kvStore.Set(p.key, value)
}

// sub returns a proxy for an element of a container.
// The separator is significant, it prevents potential clashes with other elements.
// Different separators can be used to indicate different sub-containers
func (p Proxy) sub(sep byte, key []byte) Proxy {
	if len(p.key) == 0 {
		// this must be a root proxy
		return Proxy{kvStore: p.kvStore, key: key}
	}
	return Proxy{kvStore: p.kvStore, key: append(append(p.key, sep), key...)}
}
