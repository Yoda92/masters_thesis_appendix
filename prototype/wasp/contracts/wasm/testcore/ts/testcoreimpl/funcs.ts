// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmlib from "wasmlib"
import * as wasmtypes from "wasmlib/wasmtypes";
import * as coreaccounts from "wasmlib/coreaccounts"
import * as coregovernance from "wasmlib/coregovernance"
import * as sc from "../testcore/index";

const CONTRACT_NAME_DEPLOYED = "exampleDeployTR";
const MSG_CORE_ONLY_PANIC = "========== core only =========";
const MSG_FULL_PANIC = "========== panic FULL ENTRY POINT ==========";
const MSG_VIEW_PANIC = "========== panic VIEW ==========";

export function funcCallOnChain(ctx: wasmlib.ScFuncContext, f: sc.CallOnChainContext): void {
    let paramInt = f.params.n().value();

    let hnameContract = ctx.contract();
    if (f.params.hnameContract().exists()) {
        hnameContract = f.params.hnameContract().value();
    }

    let hnameEP = sc.HFuncCallOnChain;
    if (f.params.hnameEP().exists()) {
        hnameEP = f.params.hnameEP().value();
    }

    let counter = f.state.counter();
    ctx.log("call depth = " + f.params.n().toString() +
        ", hnameContract = " + hnameContract.toString() +
        ", hnameEP = " + hnameEP.toString() +
        ", counter = " + counter.toString())

    counter.setValue(counter.value() + 1);

    let params = new wasmlib.ScDict(null);
    const key = wasmtypes.stringToBytes(sc.ParamN);
    params.set(key, wasmtypes.uint64ToBytes(paramInt))
    let ret = ctx.call(hnameContract, hnameEP, params, null);
    let retVal = wasmtypes.uint64FromBytes(ret.get(key));
    f.results.n().setValue(retVal);
}

export function funcCheckContextFromFullEP(ctx: wasmlib.ScFuncContext, f: sc.CheckContextFromFullEPContext): void {
    ctx.require(f.params.agentID().value().equals(ctx.accountID()), "fail: agentID");
    ctx.require(f.params.caller().value().equals(ctx.caller()), "fail: caller");
    ctx.require(f.params.chainID().value().equals(ctx.currentChainID()), "fail: chainID");
    ctx.require(f.params.chainOwnerID().value().equals(ctx.chainOwnerID()), "fail: chainOwnerID");
}

export function funcClaimAllowance(ctx: wasmlib.ScFuncContext, f: sc.ClaimAllowanceContext): void {
    let allowance = ctx.allowance();
    let transfer = wasmlib.ScTransfer.fromBalances(allowance);
    ctx.transferAllowed(ctx.accountID(), transfer);
}

export function funcDoNothing(ctx: wasmlib.ScFuncContext, f: sc.DoNothingContext): void {
    ctx.log("doing nothing...");
}

export function funcEstimateMinStorageDeposit(ctx: wasmlib.ScFuncContext, f: sc.EstimateMinStorageDepositContext): void {
    const provided = ctx.allowance().baseTokens();
    let dummy = sc.ScFuncs.estimateMinStorageDeposit(ctx);
    const required = ctx.estimateStorageDeposit(dummy.func);
    ctx.require(provided >= required, "not enough funds");
}

export function funcIncCounter(ctx: wasmlib.ScFuncContext, f: sc.IncCounterContext): void {
    let counter = f.state.counter();
    counter.setValue(counter.value() + 1);
}

export function funcInfiniteLoop(ctx: wasmlib.ScFuncContext, f: sc.InfiniteLoopContext): void {
    for (; ;) {
        // do nothing, just waste gas
    }
}

export function funcInit(ctx: wasmlib.ScFuncContext, f: sc.InitContext): void {
    if (f.params.fail().exists()) {
        ctx.panic("failing on purpose");
    }
}

export function funcPassTypesFull(ctx: wasmlib.ScFuncContext, f: sc.PassTypesFullContext): void {
    let hash = ctx.utility().hashBlake2b(wasmtypes.stringToBytes(sc.ParamHash));
    ctx.require(f.params.hash().value().equals(hash), "Hash wrong");
    ctx.require(f.params.int64().value() == 42, "int64 wrong");
    ctx.require(f.params.int64Zero().value() == 0, "int64-0 wrong");
    ctx.require(f.params.string().value() == sc.ParamString, "string wrong");
    ctx.require(f.params.stringZero().value() == "", "string-0 wrong");
    ctx.require(f.params.hname().value().equals(wasmtypes.ScHname.fromName(sc.ParamHname)), "Hname wrong");
    ctx.require(f.params.hnameZero().value().equals(new wasmtypes.ScHname(0)), "Hname-0 wrong");
}

export function funcPingAllowanceBack(ctx: wasmlib.ScFuncContext, f: sc.PingAllowanceBackContext): void {
    const caller = ctx.caller();
    ctx.require(caller.isAddress(), "pingAllowanceBack: caller expected to be a L1 address");
    const transfer = wasmlib.ScTransfer.fromBalances(ctx.allowance());
    ctx.transferAllowed(ctx.accountID(), transfer);
    ctx.send(caller.address(), transfer);
}

export function funcRunRecursion(ctx: wasmlib.ScFuncContext, f: sc.RunRecursionContext): void {
    let depth = f.params.n().value();
    if (depth <= 0) {
        return;
    }

    let callOnChain = sc.ScFuncs.callOnChain(ctx);
    callOnChain.params.n().setValue(depth - 1);
    callOnChain.params.hnameEP().setValue(sc.HFuncRunRecursion);
    callOnChain.func.call();
    let retVal = callOnChain.results.n().value();
    f.results.n().setValue(retVal);
}

export function funcSendLargeRequest(ctx: wasmlib.ScFuncContext, f: sc.SendLargeRequestContext): void {
}

export function funcSendNFTsBack(ctx: wasmlib.ScFuncContext, f: sc.SendNFTsBackContext): void {
    let address = ctx.caller().address();
    let allowance = ctx.allowance();
    let transfer = wasmlib.ScTransfer.fromBalances(allowance);
    ctx.transferAllowed(ctx.accountID(), transfer);
    const nftIDs = allowance.nftIDs();
    for (let i = 0; i < nftIDs.length; i++) {
        let transfer = wasmlib.ScTransfer.nft(nftIDs[i]);
        ctx.send(address, transfer);
    }
}

export function funcSendToAddress(ctx: wasmlib.ScFuncContext, f: sc.SendToAddressContext): void {
    // let transfer = wasmlib.ScTransfers.fromBalances(ctx.balances());
    // ctx.send(f.params.address().value(), transfer);
}

export function funcSetInt(ctx: wasmlib.ScFuncContext, f: sc.SetIntContext): void {
    f.state.ints().getInt64(f.params.name().value()).setValue(f.params.intValue().value());
}

export function funcSpawn(ctx: wasmlib.ScFuncContext, f: sc.SpawnContext): void {
    let programHash = f.params.progHash().value();
    let spawnName = sc.ScName + "_spawned";
    ctx.deployContract(programHash, spawnName, null);

    let spawnHname = wasmtypes.ScHname.fromName(spawnName);
    for (let i = 0; i < 5; i++) {
        ctx.call(spawnHname, sc.HFuncIncCounter, null, null);
    }
}

export function funcSplitFunds(ctx: wasmlib.ScFuncContext, f: sc.SplitFundsContext): void {
    let tokens = ctx.allowance().baseTokens();
    const address = ctx.caller().address();
    let tokensToTransfer: u64 = 1_000_000;
    const transfer = wasmlib.ScTransfer.baseTokens(tokensToTransfer);
    for (; tokens >= tokensToTransfer; tokens -= tokensToTransfer) {
        ctx.transferAllowed(ctx.accountID(), transfer);
        ctx.send(address, transfer);
    }
}

export function funcSplitFundsNativeTokens(ctx: wasmlib.ScFuncContext, f: sc.SplitFundsNativeTokensContext): void {
    let tokens = ctx.allowance().baseTokens();
    const address = ctx.caller().address();
    let transfer = wasmlib.ScTransfer.baseTokens(tokens);
    ctx.transferAllowed(ctx.accountID(), transfer);
    const tokenIDs = ctx.allowance().tokenIDs();
    const one = wasmtypes.ScBigInt.fromUint64(1);
    for (let i = 0; i < tokenIDs.length; i++) {
        const token = tokenIDs[i];
        transfer = wasmlib.ScTransfer.tokens(token, one);
        let tokens = ctx.allowance().balance(token);
        for (; tokens.cmp(one) >= 0; tokens = tokens.sub(one)) {
            ctx.transferAllowed(ctx.accountID(), transfer);
            ctx.send(address, transfer);
        }
    }
}

export function funcTestBlockContext1(ctx: wasmlib.ScFuncContext, f: sc.TestBlockContext1Context): void {
    ctx.panic(MSG_CORE_ONLY_PANIC);
}

export function funcTestBlockContext2(ctx: wasmlib.ScFuncContext, f: sc.TestBlockContext2Context): void {
    ctx.panic(MSG_CORE_ONLY_PANIC);
}

export function funcTestCallPanicFullEP(ctx: wasmlib.ScFuncContext, f: sc.TestCallPanicFullEPContext): void {
    sc.ScFuncs.testPanicFullEP(ctx).func.call();
}

export function funcTestCallPanicViewEPFromFull(ctx: wasmlib.ScFuncContext, f: sc.TestCallPanicViewEPFromFullContext): void {
    sc.ScFuncs.testPanicViewEP(ctx).func.call();
}

export function funcTestChainOwnerIDFull(ctx: wasmlib.ScFuncContext, f: sc.TestChainOwnerIDFullContext): void {
    f.results.chainOwnerID().setValue(ctx.chainOwnerID());
}

export function funcTestEventLogDeploy(ctx: wasmlib.ScFuncContext, f: sc.TestEventLogDeployContext): void {
    // deploy the same contract with another name
    let programHash = ctx.utility().hashBlake2b(wasmtypes.stringToBytes("testcore"));
    ctx.deployContract(programHash, CONTRACT_NAME_DEPLOYED, null);
}

export function funcTestEventLogEventData(ctx: wasmlib.ScFuncContext, f: sc.TestEventLogEventDataContext): void {
    f.events.test();
}

export function funcTestEventLogGenericData(ctx: wasmlib.ScFuncContext, f: sc.TestEventLogGenericDataContext): void {
    f.events.counter(f.params.counter().value());
}

export function funcTestPanicFullEP(ctx: wasmlib.ScFuncContext, f: sc.TestPanicFullEPContext): void {
    ctx.panic(MSG_FULL_PANIC);
}

export function funcWithdrawFromChain(ctx: wasmlib.ScFuncContext, f: sc.WithdrawFromChainContext): void {
    const targetChain = f.params.chainID().value();
    const withdrawal = f.params.baseTokens().value();

    // if it is not already present in the SC's account the caller should have
    // provided enough base tokens to cover the gas fees for the current call,
    // and for the storage deposit plus gas fees for the outgoing request to
    // accounts.transferAllowanceTo()
    const transfer = wasmlib.ScTransfer.fromBalances(ctx.allowance());
    ctx.transferAllowed(ctx.accountID(), transfer);

    let gasReserveTransferAccountToChain: u64 = wasmlib.MinGasFee;
    if (f.params.gasReserveTransferAccountToChain().exists()) {
        gasReserveTransferAccountToChain = f.params.gasReserveTransferAccountToChain().value();
    }
    let gasReserve: u64 = wasmlib.MinGasFee;
    if (f.params.gasReserve().exists()) {
        gasReserve = f.params.gasReserve().value();
    }
    const storageDeposit: u64 = wasmlib.StorageDeposit;

    // note: gasReserve is the gas necessary to run accounts.transferAllowanceTo
    // on the other chain by the accounts.transferAccountToChain request

    // NOTE: make sure you READ THE DOCS before calling this function
    const xfer = coreaccounts.ScFuncs.transferAccountToChain(ctx);
    xfer.params.gasReserve().setValue(gasReserve);
    xfer.func.transferBaseTokens(storageDeposit + gasReserveTransferAccountToChain + gasReserve)
        .allowanceBaseTokens(withdrawal + storageDeposit + gasReserve)
        .postToChain(targetChain);
}

export function viewCheckContextFromViewEP(ctx: wasmlib.ScViewContext, f: sc.CheckContextFromViewEPContext): void {
    ctx.require(f.params.agentID().value().equals(ctx.accountID()), "fail: agentID");
    ctx.require(f.params.chainID().value().equals(ctx.currentChainID()), "fail: chainID");
    ctx.require(f.params.chainOwnerID().value().equals(ctx.chainOwnerID()), "fail: chainOwnerID");
}

function fibonacci(n: u64): u64 {
    if (n <= 1) {
        return n;
    }
    return fibonacci(n - 1) + fibonacci(n - 2);
}

export function viewFibonacci(ctx: wasmlib.ScViewContext, f: sc.FibonacciContext): void {
    const n = f.params.n().value();
    const result = fibonacci(n);
    f.results.n().setValue(result);
}

export function viewFibonacciIndirect(ctx: wasmlib.ScViewContext, f: sc.FibonacciIndirectContext): void {
    const n = f.params.n().value();
    if (n <= 1) {
        f.results.n().setValue(n);
        return;
    }

    const fib = sc.ScFuncs.fibonacciIndirect(ctx);
    fib.params.n().setValue(n - 1);
    fib.func.call();
    const n1 = fib.results.n().value();

    fib.params.n().setValue(n - 2);
    fib.func.call();
    const n2 = fib.results.n().value();

    f.results.n().setValue(n1 + n2);
}

export function viewGetCounter(ctx: wasmlib.ScViewContext, f: sc.GetCounterContext): void {
    f.results.counter().setValue(f.state.counter().value());
}

export function viewGetInt(ctx: wasmlib.ScViewContext, f: sc.GetIntContext): void {
    let name = f.params.name().value();
    let value = f.state.ints().getInt64(name);
    ctx.require(value.exists(), "param '" + name + "' not found");
    f.results.values().getInt64(name).setValue(value.value());
}

export function viewGetStringValue(ctx: wasmlib.ScViewContext, f: sc.GetStringValueContext): void {
    ctx.panic(MSG_CORE_ONLY_PANIC);
}

export function viewInfiniteLoopView(ctx: wasmlib.ScViewContext, f: sc.InfiniteLoopViewContext): void {
    for (; ;) {
        // do nothing, just waste gas
    }
}

export function viewJustView(ctx: wasmlib.ScViewContext, f: sc.JustViewContext): void {
    ctx.log("doing nothing...");
}

export function viewPassTypesView(ctx: wasmlib.ScViewContext, f: sc.PassTypesViewContext): void {
    let hash = ctx.utility().hashBlake2b(wasmtypes.stringToBytes(sc.ParamHash));
    ctx.require(f.params.hash().value().equals(hash), "Hash wrong");
    ctx.require(f.params.int64().value() == 42, "int64 wrong");
    ctx.require(f.params.int64Zero().value() == 0, "int64-0 wrong");
    ctx.require(f.params.string().value() == sc.ParamString, "string wrong");
    ctx.require(f.params.stringZero().value() == "", "string-0 wrong");
    ctx.require(f.params.hname().value().equals(wasmtypes.ScHname.fromName(sc.ParamHname)), "Hname wrong");
    ctx.require(f.params.hnameZero().value().equals(new wasmtypes.ScHname(0)), "Hname-0 wrong");
}

export function viewTestCallPanicViewEPFromView(ctx: wasmlib.ScViewContext, f: sc.TestCallPanicViewEPFromViewContext): void {
    sc.ScFuncs.testPanicViewEP(ctx).func.call();
}

export function viewTestChainOwnerIDView(ctx: wasmlib.ScViewContext, f: sc.TestChainOwnerIDViewContext): void {
    f.results.chainOwnerID().setValue(ctx.chainOwnerID());
}

export function viewTestPanicViewEP(ctx: wasmlib.ScViewContext, f: sc.TestPanicViewEPContext): void {
    ctx.panic(MSG_VIEW_PANIC);
}

export function viewTestSandboxCall(ctx: wasmlib.ScViewContext, f: sc.TestSandboxCallContext): void {
    let getChainInfo = coregovernance.ScFuncs.getChainInfo(ctx);
    getChainInfo.func.call();
    f.results.sandboxCall().setValue(getChainInfo.results.chainID().value().toString());
}
