// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

//nolint:dupl,staticcheck
package testcoreimpl

import (
	"github.com/iotaledger/wasp/contracts/wasm/testcore/go/testcore"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/coreaccounts"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/coregovernance"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

const (
	ContractNameDeployed = "exampleDeployTR"
	MsgCoreOnlyPanic     = "========== core only ========="

	MsgDoNothing     = "========== doing nothing"
	MsgFailOnPurpose = "failing on purpose"
	MsgFullPanic     = "========== panic FULL ENTRY POINT =========="
	MsgJustView      = "calling empty view entry point"

	MsgViewPanic = "========== panic VIEW =========="
)

func funcCallOnChain(ctx wasmlib.ScFuncContext, f *CallOnChainContext) {
	n := f.Params.N().Value()

	hnameContract := ctx.Contract()
	if f.Params.HnameContract().Exists() {
		hnameContract = f.Params.HnameContract().Value()
	}

	hnameEP := testcore.HFuncCallOnChain
	if f.Params.HnameEP().Exists() {
		hnameEP = f.Params.HnameEP().Value()
	}

	counter := f.State.Counter()

	ctx.Log("param IN = " + f.Params.N().String() +
		", hnameContract = " + hnameContract.String() +
		", hnameEP = " + hnameEP.String() +
		", counter = " + counter.String())

	counter.SetValue(counter.Value() + 1)

	params := wasmlib.NewScDict()
	key := []byte(testcore.ParamN)
	params.Set(key, wasmtypes.Uint64ToBytes(n))
	ret := ctx.Call(hnameContract, hnameEP, params, nil)
	retVal := wasmtypes.Uint64FromBytes(ret.Get(key))
	f.Results.N().SetValue(retVal)
}

func funcCheckContextFromFullEP(ctx wasmlib.ScFuncContext, f *CheckContextFromFullEPContext) {
	ctx.Require(f.Params.AgentID().Value() == ctx.AccountID(), "fail: agentID")
	ctx.Require(f.Params.Caller().Value() == ctx.Caller(), "fail: caller")
	ctx.Require(f.Params.ChainID().Value() == ctx.CurrentChainID(), "fail: chainID")
	ctx.Require(f.Params.ChainOwnerID().Value() == ctx.ChainOwnerID(), "fail: chainOwnerID")
}

func funcClaimAllowance(ctx wasmlib.ScFuncContext, _ *ClaimAllowanceContext) {
	allowance := ctx.Allowance()
	transfer := wasmlib.ScTransferFromBalances(allowance)
	ctx.TransferAllowed(ctx.AccountID(), transfer)
}

func funcDoNothing(ctx wasmlib.ScFuncContext, _ *DoNothingContext) {
	ctx.Log(MsgDoNothing)
}

func funcEstimateMinStorageDeposit(ctx wasmlib.ScFuncContext, _ *EstimateMinStorageDepositContext) {
	provided := ctx.Allowance().BaseTokens()
	dummy := testcore.ScFuncs.EstimateMinStorageDeposit(ctx)
	required := ctx.EstimateStorageDeposit(dummy.Func)
	ctx.Require(provided >= required, "not enough funds")
}

func funcIncCounter(_ wasmlib.ScFuncContext, f *IncCounterContext) {
	counter := f.State.Counter()
	counter.SetValue(counter.Value() + 1)
}

//nolint:revive
func funcInfiniteLoop(_ wasmlib.ScFuncContext, _ *InfiniteLoopContext) {
	for {
		// do nothing, just waste gas
	}
}

func funcInit(ctx wasmlib.ScFuncContext, f *InitContext) {
	if f.Params.Fail().Exists() {
		ctx.Panic(MsgFailOnPurpose)
	}
}

func funcPassTypesFull(ctx wasmlib.ScFuncContext, f *PassTypesFullContext) {
	hash := ctx.Utility().HashBlake2b([]byte(testcore.ParamHash))
	ctx.Require(f.Params.Hash().Value() == hash, "wrong hash")
	ctx.Require(f.Params.Hname().Value() == wasmtypes.NewScHname(testcore.ParamHname), "wrong hname")
	ctx.Require(f.Params.HnameZero().Value() == 0, "wrong hname-0")
	ctx.Require(f.Params.Int64().Value() == 42, "wrong int64")
	ctx.Require(f.Params.Int64Zero().Value() == 0, "wrong int64-0")
	ctx.Require(f.Params.String().Value() == testcore.ParamString, "wrong string")
	ctx.Require(f.Params.StringZero().Value() == "", "wrong string-0")
	// TODO more?
}

func funcPingAllowanceBack(ctx wasmlib.ScFuncContext, _ *PingAllowanceBackContext) {
	caller := ctx.Caller()
	ctx.Require(caller.IsAddress(), "pingAllowanceBack: caller expected to be a L1 address")
	transfer := wasmlib.ScTransferFromBalances(ctx.Allowance())
	ctx.TransferAllowed(ctx.AccountID(), transfer)
	ctx.Send(caller.Address(), transfer)
}

func funcRunRecursion(ctx wasmlib.ScFuncContext, f *RunRecursionContext) {
	depth := f.Params.N().Value()
	if depth <= 0 {
		return
	}

	callOnChain := testcore.ScFuncs.CallOnChain(ctx)
	callOnChain.Params.N().SetValue(depth - 1)
	callOnChain.Params.HnameEP().SetValue(testcore.HFuncRunRecursion)
	callOnChain.Func.Call()
	retVal := callOnChain.Results.N().Value()
	f.Results.N().SetValue(retVal)
}

func funcSendLargeRequest(_ wasmlib.ScFuncContext, _ *SendLargeRequestContext) {
}

func funcSendNFTsBack(ctx wasmlib.ScFuncContext, _ *SendNFTsBackContext) {
	address := ctx.Caller().Address()
	allowance := ctx.Allowance()
	transfer := wasmlib.ScTransferFromBalances(allowance)
	ctx.TransferAllowed(ctx.AccountID(), transfer)
	for _, nftID := range allowance.NftIDs() {
		transfer = wasmlib.ScTransferFromNFT(nftID)
		ctx.Send(address, transfer)
	}
}

func funcSendToAddress(_ wasmlib.ScFuncContext, _ *SendToAddressContext) {
	// transfer := wasmlib.ScTransferFromBalances(ctx.Balances())
	// ctx.Send(f.Params.Address().Value(), transfer)
}

func funcSetInt(_ wasmlib.ScFuncContext, f *SetIntContext) {
	f.State.Ints().GetInt64(f.Params.Name().Value()).SetValue(f.Params.IntValue().Value())
}

func funcSpawn(ctx wasmlib.ScFuncContext, f *SpawnContext) {
	programHash := f.Params.ProgHash().Value()
	spawnName := testcore.ScName + "_spawned"
	ctx.DeployContract(programHash, spawnName, nil)

	spawnHname := wasmtypes.NewScHname(spawnName)
	for i := 0; i < 5; i++ {
		ctx.Call(spawnHname, testcore.HFuncIncCounter, nil, nil)
	}
}

func funcSplitFunds(ctx wasmlib.ScFuncContext, _ *SplitFundsContext) {
	tokens := ctx.Allowance().BaseTokens()
	address := ctx.Caller().Address()
	tokensToTransfer := uint64(1_000_000)
	transfer := wasmlib.ScTransferFromBaseTokens(tokensToTransfer)
	for ; tokens >= tokensToTransfer; tokens -= tokensToTransfer {
		ctx.TransferAllowed(ctx.AccountID(), transfer)
		ctx.Send(address, transfer)
	}
}

func funcSplitFundsNativeTokens(ctx wasmlib.ScFuncContext, _ *SplitFundsNativeTokensContext) {
	tokens := ctx.Allowance().BaseTokens()
	address := ctx.Caller().Address()
	transfer := wasmlib.ScTransferFromBaseTokens(tokens)
	ctx.TransferAllowed(ctx.AccountID(), transfer)
	for _, token := range ctx.Allowance().TokenIDs() {
		one := wasmtypes.NewScBigInt(1)
		transfer = wasmlib.ScTransferFromTokens(token, one)
		tokens := ctx.Allowance().Balance(token)
		for ; tokens.Cmp(one) >= 0; tokens = tokens.Sub(one) {
			ctx.TransferAllowed(ctx.AccountID(), transfer)
			ctx.Send(address, transfer)
		}
	}
}

func funcTestBlockContext1(ctx wasmlib.ScFuncContext, _ *TestBlockContext1Context) {
	ctx.Panic(MsgCoreOnlyPanic)
}

func funcTestBlockContext2(ctx wasmlib.ScFuncContext, _ *TestBlockContext2Context) {
	ctx.Panic(MsgCoreOnlyPanic)
}

func funcTestCallPanicFullEP(ctx wasmlib.ScFuncContext, _ *TestCallPanicFullEPContext) {
	ctx.Log("will be calling entry point '" + testcore.FuncTestPanicFullEP + "' from full EP")
	testcore.ScFuncs.TestPanicFullEP(ctx).Func.Call()
}

func funcTestCallPanicViewEPFromFull(ctx wasmlib.ScFuncContext, _ *TestCallPanicViewEPFromFullContext) {
	ctx.Log("will be calling entry point '" + testcore.ViewTestPanicViewEP + "' from full EP")
	testcore.ScFuncs.TestPanicViewEP(ctx).Func.Call()
}

func funcTestChainOwnerIDFull(ctx wasmlib.ScFuncContext, f *TestChainOwnerIDFullContext) {
	f.Results.ChainOwnerID().SetValue(ctx.ChainOwnerID())
}

func funcTestEventLogDeploy(ctx wasmlib.ScFuncContext, _ *TestEventLogDeployContext) {
	// deploy the same contract with another name
	programHash := ctx.Utility().HashBlake2b([]byte(testcore.ScName))
	ctx.DeployContract(programHash, ContractNameDeployed, nil)
}

func funcTestEventLogEventData(ctx wasmlib.ScFuncContext, f *TestEventLogEventDataContext) {
	f.Events.Test()
}

func funcTestEventLogGenericData(ctx wasmlib.ScFuncContext, f *TestEventLogGenericDataContext) {
	f.Events.Counter(f.Params.Counter().Value())
}

func funcTestPanicFullEP(ctx wasmlib.ScFuncContext, _ *TestPanicFullEPContext) {
	ctx.Panic(MsgFullPanic)
}

func funcWithdrawFromChain(ctx wasmlib.ScFuncContext, f *WithdrawFromChainContext) {
	targetChain := f.Params.ChainID().Value()
	withdrawal := f.Params.BaseTokens().Value()

	// if it is not already present in the SC's account the caller should have
	// provided enough base tokens to cover the gas fees for the current call,
	// and for the storage deposit plus gas fees for the outgoing request to
	// accounts.transferAllowanceTo()
	transfer := wasmlib.ScTransferFromBalances(ctx.Allowance())
	ctx.TransferAllowed(ctx.AccountID(), transfer)

	gasReserveTransferAccountToChain := wasmlib.MinGasFee
	if f.Params.GasReserveTransferAccountToChain().Exists() {
		gasReserveTransferAccountToChain = f.Params.GasReserveTransferAccountToChain().Value()
	}
	gasReserve := wasmlib.MinGasFee
	if f.Params.GasReserve().Exists() {
		gasReserve = f.Params.GasReserve().Value()
	}
	const storageDeposit = wasmlib.StorageDeposit

	// note: gasReserve is the gas necessary to run accounts.transferAllowanceTo
	// on the other chain by the accounts.transferAccountToChain request

	// NOTE: make sure you READ THE DOCS before calling this function
	xfer := coreaccounts.ScFuncs.TransferAccountToChain(ctx)
	xfer.Params.GasReserve().SetValue(gasReserve)
	xfer.Func.TransferBaseTokens(storageDeposit + gasReserveTransferAccountToChain + gasReserve).
		AllowanceBaseTokens(withdrawal + storageDeposit + gasReserve).
		PostToChain(targetChain)
}

func viewCheckContextFromViewEP(ctx wasmlib.ScViewContext, f *CheckContextFromViewEPContext) {
	ctx.Require(f.Params.AgentID().Value() == ctx.AccountID(), "fail: agentID")
	ctx.Require(f.Params.ChainID().Value() == ctx.CurrentChainID(), "fail: chainID")
	ctx.Require(f.Params.ChainOwnerID().Value() == ctx.ChainOwnerID(), "fail: chainOwnerID")
}

func fibonacci(n uint64) uint64 {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func viewFibonacci(_ wasmlib.ScViewContext, f *FibonacciContext) {
	n := f.Params.N().Value()
	result := fibonacci(n)
	f.Results.N().SetValue(result)
}

func viewFibonacciIndirect(ctx wasmlib.ScViewContext, f *FibonacciIndirectContext) {
	n := f.Params.N().Value()
	if n == 0 || n == 1 {
		f.Results.N().SetValue(n)
		return
	}

	fib := testcore.ScFuncs.FibonacciIndirect(ctx)
	fib.Params.N().SetValue(n - 1)
	fib.Func.Call()
	n1 := fib.Results.N().Value()

	fib.Params.N().SetValue(n - 2)
	fib.Func.Call()
	n2 := fib.Results.N().Value()

	f.Results.N().SetValue(n1 + n2)
}

func viewGetCounter(_ wasmlib.ScViewContext, f *GetCounterContext) {
	f.Results.Counter().SetValue(f.State.Counter().Value())
}

func viewGetInt(ctx wasmlib.ScViewContext, f *GetIntContext) {
	name := f.Params.Name().Value()
	value := f.State.Ints().GetInt64(name)
	ctx.Require(value.Exists(), "param '"+name+"' not found")
	f.Results.Values().GetInt64(name).SetValue(value.Value())
}

func viewGetStringValue(ctx wasmlib.ScViewContext, _ *GetStringValueContext) {
	ctx.Panic(MsgCoreOnlyPanic)
	// varName := f.Params.VarName().Value()
	// value := f.State.Strings().GetString(varName).Value()
	// f.Results.Vars().GetString(varName).SetValue(value)
}

//nolint:revive
func viewInfiniteLoopView(_ wasmlib.ScViewContext, _ *InfiniteLoopViewContext) {
	for {
		// do nothing, just waste gas
	}
}

func viewJustView(ctx wasmlib.ScViewContext, _ *JustViewContext) {
	ctx.Log(MsgJustView)
}

func viewPassTypesView(ctx wasmlib.ScViewContext, f *PassTypesViewContext) {
	hash := ctx.Utility().HashBlake2b([]byte(testcore.ParamHash))
	ctx.Require(f.Params.Hash().Value() == hash, "wrong hash")
	ctx.Require(f.Params.Hname().Value() == wasmtypes.NewScHname(testcore.ParamHname), "wrong hname")
	ctx.Require(f.Params.HnameZero().Value() == 0, "wrong hname-0")
	ctx.Require(f.Params.Int64().Value() == 42, "wrong int64")
	ctx.Require(f.Params.Int64Zero().Value() == 0, "wrong int64-0")
	ctx.Require(f.Params.String().Value() == testcore.ParamString, "wrong string")
	ctx.Require(f.Params.StringZero().Value() == "", "wrong string-0")
	// TODO more?
}

func viewTestCallPanicViewEPFromView(ctx wasmlib.ScViewContext, _ *TestCallPanicViewEPFromViewContext) {
	ctx.Log("will be calling entry point '" + testcore.ViewTestPanicViewEP + "' from view EP")
	testcore.ScFuncs.TestPanicViewEP(ctx).Func.Call()
}

func viewTestChainOwnerIDView(ctx wasmlib.ScViewContext, f *TestChainOwnerIDViewContext) {
	f.Results.ChainOwnerID().SetValue(ctx.ChainOwnerID())
}

func viewTestPanicViewEP(ctx wasmlib.ScViewContext, _ *TestPanicViewEPContext) {
	ctx.Panic(MsgViewPanic)
}

func viewTestSandboxCall(ctx wasmlib.ScViewContext, f *TestSandboxCallContext) {
	getChainInfo := coregovernance.ScFuncs.GetChainInfo(ctx)
	getChainInfo.Func.Call()
	f.Results.SandboxCall().SetValue(getChainInfo.Results.ChainID().Value().String())
}
