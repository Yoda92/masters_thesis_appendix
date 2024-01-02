// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package tstemplates

var paramsTs = map[string]string{
	// *******************************
	"params.ts": `
$#emit importWasmTypes
$#emit importSc
$#each func paramsFunc
`,
	// *******************************
	"paramsFunc": `
$#if params paramsFuncParams
`,
	// *******************************
	"paramsFuncParams": `
$#set Kind Param
$#set mut Immutable
$#if param paramsProxyStruct
$#set mut Mutable
$#if param paramsProxyStruct
`,
	// *******************************
	"paramsProxyStruct": `
$#set TypeName $mut$FuncName$+Params
$#each param proxyContainers

export class $TypeName extends wasmtypes.ScProxy {
$#set separator $false
$#each param proxyMethods
}
`,
}
