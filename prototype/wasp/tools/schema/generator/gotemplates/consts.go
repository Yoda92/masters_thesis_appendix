// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gotemplates

var constsGo = map[string]string{
	// *******************************
	"consts.go": `
package $package

$#emit importWasmTypes

const (
	ScName        = "$scName"
	ScDescription = "$scDesc"
	HScName       = wasmtypes.ScHname(0x$hscName)
)
$#if params constParams
$#if results constResults
$#if state constState
$#if funcs constFuncs
`,
	// *******************************
	"constParams": `

const (
$#set constPrefix Param
$#each params constField
)
`,
	// *******************************
	"constResults": `

const (
$#set constPrefix Result
$#each results constField
)
`,
	// *******************************
	"constState": `

const (
$#set constPrefix State
$#each state constField
)
`,
	// *******************************
	"constFuncs": `

const (
$#each func constFunc
)

const (
$#each func constHFunc
)
`,
	// *******************************
	"constField": `
	$constPrefix$FldName$fldPad = "$fldAlias"
`,
	// *******************************
	"constFunc": `
	$Kind$FuncName$funcPad = "$funcAlias"
`,
	// *******************************
	"constHFunc": `
	H$Kind$FuncName$funcPad = wasmtypes.ScHname(0x$hFuncName)
`,
}
