// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package gotemplates

var contractGo = map[string]string{
	// *******************************
	"contract.go": `
package $package
$#if funcs emitContract
`,
	// *******************************
	"emitContract": `

$#emit importWasmLib
$#each func FuncNameCall

type Funcs struct{}

var ScFuncs Funcs
$#each func FuncNameForCall
$#if core coreOnDispatch
`,
	// *******************************
	"FuncNameCall": `
$#emit alignCalculate
$#emit setupInitFunc

type $FuncName$+Call struct {
	Func$falign *wasmlib.Sc$initFunc$Kind
$#if param MutableFuncNameParams
$#if result ImmutableFuncNameResults
}
`,
	// *******************************
	"MutableFuncNameParams": `
	Params$align Mutable$FuncName$+Params
`,
	// *******************************
	"ImmutableFuncNameResults": `
	Results Immutable$FuncName$+Results
`,
	// *******************************
	"FuncNameForCall": `
$#emit setupInitFunc

$#each funcComment _funcComment
func (sc Funcs) $FuncName(ctx wasmlib.Sc$Kind$+ClientContext) *$FuncName$+Call {
$#set thisView f.Func
$#if func setThisView
$#set complex $false
$#if param setComplex
$#if result setComplex
$#if complex initComplex initSimple
}
`,
	// *******************************
	"setThisView": `
$#set thisView &f.Func.ScView
`,
	// *******************************
	"setComplex": `
$#set complex $true
`,
	// *******************************
	"coreOnDispatch": `

var exportMap = wasmlib.ScExportMap{
	Names: []string{
$#each func coreExportName
	},
	Funcs: []wasmlib.ScFuncContextFunction{
$#each func coreExportFunc
	},
	Views: []wasmlib.ScViewContextFunction{
$#each func coreExportView
	},
}

func OnDispatch(index int32) *wasmlib.ScExportMap {
	if index < 0 {
		return exportMap.Dispatch(index)
	}

	panic("Calling core contract?")
}
`,
	// *******************************
	"initComplex": `
	f := &$FuncName$+Call{Func: wasmlib.NewSc$initFunc$Kind(ctx, HScName, H$Kind$FuncName)}
$#if param initParams
$#if result initResults
	return f
`,
	// *******************************
	"initParams": `
	f.Params.Proxy = wasmlib.NewCallParamsProxy($thisView)
`,
	// *******************************
	"initResults": `
	wasmlib.NewCallResultsProxy($thisView, &f.Results.Proxy)
`,
	// *******************************
	"initSimple": `
	return &$FuncName$+Call{Func: wasmlib.NewSc$initFunc$Kind(ctx, HScName, H$Kind$FuncName)}
`,
	// *******************************
	"coreExportName": `
		$Kind$FuncName,
`,
	// *******************************
	"coreExportFunc": `
$#if func coreExportFuncThunk
`,
	// *******************************
	"coreExportFuncThunk": `
		wasmlib.$Kind$+Error,
`,
	// *******************************
	"coreExportView": `
$#if view coreExportViewThunk
`,
	// *******************************
	"coreExportViewThunk": `
		wasmlib.$Kind$+Error,
`,
}
