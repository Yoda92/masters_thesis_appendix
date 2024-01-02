package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/contracts/wasm/testcore/go/testcore"
)

func TestMainCallsFromFullEP(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w, true)
		user := ctx.Creator()

		f := testcore.ScFuncs.CheckContextFromFullEP(ctx.Sign(user))
		f.Params.ChainID().SetValue(ctx.CurrentChainID())
		f.Params.AgentID().SetValue(ctx.AccountID())
		f.Params.Caller().SetValue(user.ScAgentID())
		f.Params.ChainOwnerID().SetValue(ctx.Originator().ScAgentID())
		f.Func.Post()
		require.NoError(t, ctx.Err)
	})
}

func TestMainCallsFromViewEP(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w, true)

		f := testcore.ScFuncs.CheckContextFromViewEP(ctx)
		f.Params.ChainID().SetValue(ctx.CurrentChainID())
		f.Params.AgentID().SetValue(ctx.AccountID())
		f.Params.ChainOwnerID().SetValue(ctx.Originator().ScAgentID())
		f.Func.Call()
		require.NoError(t, ctx.Err)
	})
}

//func TestMintedSupplyOk(t *testing.T) {
//	// TODO no minting yet
//	t.SkipNow()
//	run2(t, func(t *testing.T, w bool) {
//		ctx := deployTestCore(t, w, true)
//		user := ctx.Creator()
//
//		f := testcore.ScFuncs.GetMintedSupply(ctx.Sign(user, 42))
//		f.Func.Post()
//		require.NoError(t, ctx.Err)
//
//		mintedColor, mintedAmount := ctx.Minted()
//
//		requests := int64(2)
//		if w {
//			requests++
//		}
//
//		require.EqualValues(t, solo.Saldo-42-requests, user.Balance())
//		require.EqualValues(t, 42, user.Balance(mintedColor))
//
//		require.EqualValues(t, mintedColor, f.Results.MintedColor().Value())
//		require.EqualValues(t, mintedAmount, f.Results.MintedSupply().Value())
//	})
//}
