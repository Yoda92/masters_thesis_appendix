package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/contracts/wasm/testcore/go/testcore"
	"github.com/iotaledger/wasp/packages/testutil/utxodb"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmsolo"
)

func TestOffLedgerFailNoAccount(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		// note: create agent without depositing into L2
		user := wasmsolo.NewSoloAgent(ctx.Chain.Env, "user")
		require.EqualValues(t, utxodb.FundsFromFaucetAmount, user.Balance())
		require.EqualValues(t, 0, ctx.Balance(user))
		bal := ctx.Balances(user)

		// no deposit yet, so account is unverified

		f := testcore.ScFuncs.SetInt(ctx.OffLedger(user))
		f.Params.Name().SetValue("ppp")
		f.Params.IntValue().SetValue(314)
		f.Func.Post()
		require.Error(t, ctx.Err)
		require.Contains(t, ctx.Err.Error(), "gas budget exceeded")
		bal.VerifyBalances(t)
	})
}

func TestOffLedgerNoTransfer(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		user := ctx.NewSoloAgent("user")
		bal := ctx.Balances(user)
		userL1 := user.Balance()

		// we're using setInt() here to be able to verify the state update was done
		f := testcore.ScFuncs.SetInt(ctx.OffLedger(user))
		f.Params.Name().SetValue("ppp")
		f.Params.IntValue().SetValue(314)
		f.Func.Post()
		require.NoError(t, ctx.Err)

		bal.Add(user, -ctx.GasFee)
		bal.UpdateFeeBalances(ctx.GasFee)
		bal.VerifyBalances(t)
		require.EqualValues(t, userL1, user.Balance())

		// verify state update
		v := testcore.ScFuncs.GetInt(ctx)
		v.Params.Name().SetValue("ppp")
		v.Func.Call()
		require.NoError(t, ctx.Err)
		require.EqualValues(t, 314, v.Results.Values().GetInt64("ppp").Value())
	})
}

func TestOffLedgerTransferWhenEnoughBudget(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		user := ctx.NewSoloAgent("user")
		bal := ctx.Balances(user)
		userL1 := user.Balance()

		// Allow 4321 tokens to be transferred, there's enough budget
		// note that SetInt() will not try to grab them
		f := testcore.ScFuncs.SetInt(ctx.OffLedger(user))
		f.Params.Name().SetValue("ppp")
		f.Params.IntValue().SetValue(314)
		f.Func.TransferBaseTokens(4321).Post()
		require.NoError(t, ctx.Err)
		ctx.Balances(user)

		bal.UpdateFeeBalances(ctx.GasFee)
		bal.Add(user, -ctx.GasFee)
		bal.VerifyBalances(t)
		require.EqualValues(t, userL1, user.Balance())

		// verify state update
		v := testcore.ScFuncs.GetInt(ctx)
		v.Params.Name().SetValue("ppp")
		v.Func.Call()
		require.NoError(t, ctx.Err)
		require.EqualValues(t, 314, v.Results.Values().GetInt64("ppp").Value())
	})
}

func TestOffLedgerTransferWhenNotEnoughBudget(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		user := ctx.NewSoloAgent("user")
		bal := ctx.Balances(user)

		// Try to transfer everything from L2, which does preclude paying for gas
		// note that SetInt() will not try to grab the allowance
		f := testcore.ScFuncs.SetInt(ctx.OffLedger(user))
		f.Params.Name().SetValue("ppp")
		f.Params.IntValue().SetValue(314)
		f.Func.TransferBaseTokens(ctx.Balance(user)).Post()
		require.Error(t, ctx.Err)
		require.Contains(t, ctx.Err.Error(), "gas budget exceeded")

		bal.UpdateFeeBalances(ctx.GasFee)
		bal.Add(user, -ctx.GasFee)
		bal.VerifyBalances(t)

		// verify no state update
		v := testcore.ScFuncs.GetInt(ctx)
		v.Params.Name().SetValue("ppp")
		v.Func.Call()
		if w {
			require.Error(t, ctx.Err)
			require.Contains(t, ctx.Err.Error(), "param 'ppp' not found")
		} else {
			require.NoError(t, ctx.Err)
		}
		require.EqualValues(t, 0, v.Results.Values().GetInt64("ppp").Value())
	})
}
