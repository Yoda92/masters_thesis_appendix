// smart contract for testing
package sbtestsc

import (
	"github.com/iotaledger/wasp/packages/isc/coreutil"
)

var Contract = coreutil.NewContract("testcore")

var Processor = Contract.Processor(initialize,
	FuncChainOwnerIDView.WithHandler(testChainOwnerIDView),
	FuncChainOwnerIDFull.WithHandler(testChainOwnerIDFull),

	FuncEventLogGenericData.WithHandler(testEventLogGenericData),
	FuncEventLogEventData.WithHandler(testEventLogEventData),
	FuncEventLogDeploy.WithHandler(testEventLogDeploy),
	FuncSandboxCall.WithHandler(testSandboxCall),

	FuncPanicFullEP.WithHandler(testPanicFullEP),
	FuncPanicViewEP.WithHandler(testPanicViewEP),
	FuncCallPanicFullEP.WithHandler(testCallPanicFullEP),
	FuncCallPanicViewEPFromFull.WithHandler(testCallPanicViewEPFromFull),
	FuncCallPanicViewEPFromView.WithHandler(testCallPanicViewEPFromView),

	FuncDoNothing.WithHandler(doNothing),
	// FuncSendToAddress.WithHandler(sendToAddress),

	FuncWithdrawFromChain.WithHandler(withdrawFromChain),
	FuncCallOnChain.WithHandler(callOnChain),
	FuncSetInt.WithHandler(setInt),
	FuncGetInt.WithHandler(getInt),
	FuncGetFibonacci.WithHandler(getFibonacci),
	FuncGetFibonacciIndirect.WithHandler(getFibonacciIndirect),
	FuncCalcFibonacciIndirectStoreValue.WithHandler(calcFibonacciIndirectStoreValue),
	FuncViewCalcFibonacciResult.WithHandler(viewFibResult),
	FuncIncCounter.WithHandler(incCounter),
	FuncGetCounter.WithHandler(getCounter),
	FuncRunRecursion.WithHandler(runRecursion),

	FuncPassTypesFull.WithHandler(passTypesFull),
	FuncPassTypesView.WithHandler(passTypesView),
	FuncCheckContextFromFullEP.WithHandler(testCheckContextFromFullEP),
	FuncCheckContextFromViewEP.WithHandler(testCheckContextFromViewEP),

	FuncJustView.WithHandler(testJustView),

	FuncSpawn.WithHandler(spawn),

	FuncSplitFunds.WithHandler(testSplitFunds),
	FuncSplitFundsNativeTokens.WithHandler(testSplitFundsNativeTokens),
	FuncPingAllowanceBack.WithHandler(pingAllowanceBack),
	FuncSendLargeRequest.WithHandler(sendLargeRequest),
	FuncEstimateMinStorageDeposit.WithHandler(testEstimateMinimumStorageDeposit),
	FuncInfiniteLoop.WithHandler(infiniteLoop),
	FuncInfiniteLoopView.WithHandler(infiniteLoopView),
	FuncSendNFTsBack.WithHandler(sendNFTsBack),
	FuncClaimAllowance.WithHandler(claimAllowance),
)

var (
	// function eventlog test
	FuncEventLogGenericData = coreutil.Func("testEventLogGenericData")
	FuncEventLogEventData   = coreutil.Func("testEventLogEventData")
	FuncEventLogDeploy      = coreutil.Func("testEventLogDeploy")

	// Function sandbox test
	FuncChainOwnerIDView = coreutil.ViewFunc("testChainOwnerIDView")
	FuncChainOwnerIDFull = coreutil.Func("testChainOwnerIDFull")

	FuncSandboxCall            = coreutil.ViewFunc("testSandboxCall")
	FuncCheckContextFromFullEP = coreutil.Func("checkContextFromFullEP")
	FuncCheckContextFromViewEP = coreutil.ViewFunc("checkContextFromViewEP")

	FuncPanicFullEP             = coreutil.Func("testPanicFullEP")
	FuncPanicViewEP             = coreutil.ViewFunc("testPanicViewEP")
	FuncCallPanicFullEP         = coreutil.Func("testCallPanicFullEP")
	FuncCallPanicViewEPFromFull = coreutil.Func("testCallPanicViewEPFromFull")
	FuncCallPanicViewEPFromView = coreutil.ViewFunc("testCallPanicViewEPFromView")

	FuncWithdrawFromChain = coreutil.Func("withdrawFromChain")

	FuncDoNothing = coreutil.Func("doNothing")
	// FuncSendToAddress = coreutil.Func("sendToAddress")
	FuncJustView = coreutil.ViewFunc("justView")

	FuncCallOnChain                     = coreutil.Func("callOnChain")
	FuncSetInt                          = coreutil.Func("setInt")
	FuncGetInt                          = coreutil.ViewFunc("getInt")
	FuncGetFibonacci                    = coreutil.ViewFunc("fibonacci")
	FuncGetFibonacciIndirect            = coreutil.ViewFunc("fibonacciIndirect")
	FuncCalcFibonacciIndirectStoreValue = coreutil.Func("calcFibonacciIndirectStoreValue")
	FuncViewCalcFibonacciResult         = coreutil.ViewFunc("getFibCalcResult")
	FuncGetCounter                      = coreutil.ViewFunc("getCounter")
	FuncIncCounter                      = coreutil.Func("incCounter")
	FuncRunRecursion                    = coreutil.Func("runRecursion")

	FuncPassTypesFull = coreutil.Func("passTypesFull")
	FuncPassTypesView = coreutil.ViewFunc("passTypesView")

	FuncSpawn = coreutil.Func("spawn")

	FuncSplitFunds                = coreutil.Func("splitFunds")
	FuncSplitFundsNativeTokens    = coreutil.Func("splitFundsNativeTokens")
	FuncPingAllowanceBack         = coreutil.Func("pingAllowanceBack")
	FuncSendLargeRequest          = coreutil.Func("sendLargeRequest")
	FuncEstimateMinStorageDeposit = coreutil.Func("estimateMinStorageDeposit")
	FuncInfiniteLoop              = coreutil.Func("infiniteLoop")
	FuncInfiniteLoopView          = coreutil.ViewFunc("infiniteLoopView")
	FuncSendNFTsBack              = coreutil.Func("sendNFTsBack")
	FuncClaimAllowance            = coreutil.Func("claimAllowance")
)

const (
	// State variables
	VarCounter              = "counter"
	VarContractNameDeployed = "exampleDeployTR"

	// parameters
	ParamAddress                          = "address"
	ParamAgentID                          = "agentID"
	ParamCaller                           = "caller"
	ParamChainID                          = "chainID"
	ParamChainOwnerID                     = "chainOwnerID"
	ParamGasReserve                       = "gasReserve"
	ParamGasReserveTransferAccountToChain = "gasReserveTransferAccountToChain"
	ParamContractID                       = "contractID"
	ParamFail                             = "initFailParam"
	ParamHnameContract                    = "hnameContract"
	ParamHnameEP                          = "hnameEP"
	ParamIntParamName                     = "intParamName"
	ParamIntParamValue                    = "intParamValue"
	ParamBaseTokens                       = "baseTokens"
	ParamN                                = "n"
	ParamProgHash                         = "progHash"
	ParamSize                             = "size"

	// error fragments for testing
	MsgDoNothing = "========== doing nothing =========="
	MsgFullPanic = "========== panic FULL ENTRY POINT =========="
	MsgViewPanic = "========== panic VIEW =========="
)
