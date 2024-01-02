// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

//go:build wasm
// +build wasm

package wasmvmhost

import "github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"

// interface WasmLib to the VM host

// These 2 external functions are funneling the entire
// WasmLib functionality to their counterparts on the host.

//go:wasm-module WasmLib
//export hostStateGet
func hostStateGet(keyRef *byte, keyLen int32, valRef *byte, valLen int32) int32

//go:wasm-module WasmLib
//export hostStateSet
func hostStateSet(keyRef *byte, keyLen int32, valRef *byte, valLen int32)

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// ptr returns pointer to slice or nil when slice is empty
func ptr(buf []byte) *byte {
	// &buf[0] will panic on zero length slice, so use nil instead
	if len(buf) == 0 {
		return nil
	}
	return &buf[0]
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

type WasmVMHost struct {
}

// implements wasmlib.ScHost interface
var _ wasmlib.ScHost = &WasmVMHost{}

func ConnectWasmHost() {
	w := &WasmVMHost{}
	wasmlib.ConnectHost(w)
}

func (w *WasmVMHost) ExportName(index int32, name string) {
	// nil key indicates export name, with keyLen indicating index
	// this removes the need for a separate hostExportName function
	buf := []byte(name)
	hostStateSet(nil, index, ptr(buf), int32(len(buf)))
}

func (w *WasmVMHost) Sandbox(funcNr int32, params []byte) []byte {
	// call sandbox function, result value will be cached by host
	// always negative funcNr as keyLen indicates sandbox call
	// this removes the need for a separate hostSandbox function
	size := hostStateGet(nil, funcNr, ptr(params), int32(len(params)))

	// zero length, no need to retrieve cached value
	if size == 0 {
		return []byte{}
	}

	// retrieve cached result value from host
	result := make([]byte, size)
	_ = hostStateGet(nil, 0, &result[0], size)
	return result
}

func (w *WasmVMHost) StateDelete(key []byte) {
	// value size -1 means delete key
	// this removes the need for a separate hostStateDel function
	hostStateSet(&key[0], int32(len(key)), nil, -1)
}

func (w *WasmVMHost) StateExists(key []byte) bool {
	// value size -1 means only test for existence
	// returned size -1 indicates key not found (or error)
	// this removes the need for a separate hostStateExists function
	return hostStateGet(&key[0], int32(len(key)), nil, -1) >= 0
}

func (w *WasmVMHost) StateGet(key []byte) []byte {
	// variable sized result expected,
	// query size first by passing zero length buffer
	// value will be cached by host
	size := hostStateGet(&key[0], int32(len(key)), nil, 0)

	// -1 means non-existent
	if size < 0 {
		return []byte(nil)
	}

	// zero length, no need to retrieve cached value
	if size == 0 {
		return []byte{}
	}

	// retrieve cached value from host
	value := make([]byte, size)
	_ = hostStateGet(nil, 0, &value[0], size)
	return value
}

func (w *WasmVMHost) StateSet(key, value []byte) {
	hostStateSet(&key[0], int32(len(key)), ptr(value), int32(len(value)))
}
