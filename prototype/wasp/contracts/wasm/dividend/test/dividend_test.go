// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/contracts/wasm/dividend/go/dividend"
	"github.com/iotaledger/wasp/contracts/wasm/dividend/go/dividendimpl"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
)

func dividendMember(ctx *wasmsolo.SoloContext, agent *wasmsolo.SoloAgent, factor uint64) {
	member := dividend.ScFuncs.Member(ctx)
	member.Params.Address().SetValue(agent.ScAgentID().Address())
	member.Params.Factor().SetValue(factor)
	member.Func.Post()
}

func dividendDivide(ctx *wasmsolo.SoloContext, amount uint64) {
	divide := dividend.ScFuncs.Divide(ctx)
	divide.Func.TransferBaseTokens(amount).Post()
}

func dividendGetFactor(ctx *wasmsolo.SoloContext, member *wasmsolo.SoloAgent) uint64 {
	getFactor := dividend.ScFuncs.GetFactor(ctx)
	getFactor.Params.Address().SetValue(member.ScAgentID().Address())
	getFactor.Func.Call()
	value := getFactor.Results.Factor().Value()
	return value
}

func TestDeploy(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)
	require.NoError(t, ctx.ContractExists(dividend.ScName))
}

func TestAddMemberOk(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member1 := ctx.NewSoloAgent("member1")
	dividendMember(ctx, member1, 100)
	require.NoError(t, ctx.Err)
}

func TestAddMemberFailMissingAddress(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member := dividend.ScFuncs.Member(ctx)
	member.Params.Factor().SetValue(100)
	member.Func.Post()
	require.Error(t, ctx.Err)
	require.Contains(t, ctx.Err.Error(), "missing mandatory param: address")
}

func TestAddMemberFailMissingFactor(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member1 := ctx.NewSoloAgent("member1")
	member := dividend.ScFuncs.Member(ctx)
	member.Params.Address().SetValue(member1.ScAgentID().Address())
	member.Func.Post()
	require.Error(t, ctx.Err)
	require.Contains(t, ctx.Err.Error(), "missing mandatory param: factor")
}

func TestDivide1Member(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member1 := ctx.NewSoloAgent("member1")
	bal := ctx.Balances(member1)

	dividendMember(ctx, member1, 1000)
	require.NoError(t, ctx.Err)

	bal.Originator += ctx.StorageDeposit
	bal.VerifyBalances(t)

	const dividendToDivide = 1*isc.Million + 1
	dividendDivide(ctx, dividendToDivide)
	require.NoError(t, ctx.Err)

	bal.Add(member1, dividendToDivide)
	bal.VerifyBalances(t)
}

func TestDivide2Members(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member1 := ctx.NewSoloAgent("member1")
	bal := ctx.Balances(member1)

	dividendMember(ctx, member1, 250)
	require.NoError(t, ctx.Err)

	bal.Originator += ctx.StorageDeposit
	bal.VerifyBalances(t)

	member2 := ctx.NewSoloAgent("member2")
	bal = ctx.Balances(member1, member2)

	dividendMember(ctx, member2, 750)
	require.NoError(t, ctx.Err)

	bal.Originator += ctx.StorageDeposit
	bal.VerifyBalances(t)

	const dividendToDivide = 2*isc.Million - 1
	dividendDivide(ctx, dividendToDivide)
	require.NoError(t, ctx.Err)

	remain := dividendToDivide - dividendToDivide*250/1000 - dividendToDivide*750/1000
	bal.Originator += remain
	bal.Add(member1, dividendToDivide*250/1000)
	bal.Add(member2, dividendToDivide*750/1000)
	bal.VerifyBalances(t)
}

func TestDivide3Members(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member1 := ctx.NewSoloAgent("member1")
	bal := ctx.Balances(member1)

	dividendMember(ctx, member1, 250)
	require.NoError(t, ctx.Err)

	bal.Originator += ctx.StorageDeposit
	bal.VerifyBalances(t)

	member2 := ctx.NewSoloAgent("member2")
	bal = ctx.Balances(member1, member2)

	dividendMember(ctx, member2, 500)
	require.NoError(t, ctx.Err)

	bal.Originator += ctx.StorageDeposit
	bal.VerifyBalances(t)

	member3 := ctx.NewSoloAgent("member3")
	bal = ctx.Balances(member1, member2, member3)

	dividendMember(ctx, member3, 750)
	require.NoError(t, ctx.Err)

	bal.Originator += ctx.StorageDeposit
	bal.VerifyBalances(t)

	const dividendToDivide = 2*isc.Million - 1
	dividendDivide(ctx, dividendToDivide)
	require.NoError(t, ctx.Err)

	remain := dividendToDivide - dividendToDivide*250/1500 - dividendToDivide*500/1500 - dividendToDivide*750/1500
	bal.Originator += remain
	bal.Add(member1, dividendToDivide*250/1500)
	bal.Add(member2, dividendToDivide*500/1500)
	bal.Add(member3, dividendToDivide*750/1500)
	bal.VerifyBalances(t)

	const dividendToDivide2 = 2*isc.Million + 234
	dividendDivide(ctx, dividendToDivide2)
	require.NoError(t, ctx.Err)

	remain = dividendToDivide2 - dividendToDivide2*250/1500 - dividendToDivide2*500/1500 - dividendToDivide2*750/1500
	bal.Originator += remain
	bal.Add(member1, dividendToDivide2*250/1500)
	bal.Add(member2, dividendToDivide2*500/1500)
	bal.Add(member3, dividendToDivide2*750/1500)
	bal.VerifyBalances(t)
}

func TestGetFactor(t *testing.T) {
	ctx := wasmsolo.NewSoloContext(t, dividend.ScName, dividendimpl.OnDispatch)

	member1 := ctx.NewSoloAgent("member1")
	dividendMember(ctx, member1, 250)
	require.NoError(t, ctx.Err)

	member2 := ctx.NewSoloAgent("member2")
	dividendMember(ctx, member2, 500)
	require.NoError(t, ctx.Err)

	member3 := ctx.NewSoloAgent("member3")
	dividendMember(ctx, member3, 750)
	require.NoError(t, ctx.Err)

	value := dividendGetFactor(ctx, member3)
	require.NoError(t, ctx.Err)
	require.EqualValues(t, 750, value)

	value = dividendGetFactor(ctx, member2)
	require.NoError(t, ctx.Err)
	require.EqualValues(t, 500, value)

	value = dividendGetFactor(ctx, member1)
	require.NoError(t, ctx.Err)
	require.EqualValues(t, 250, value)
}
