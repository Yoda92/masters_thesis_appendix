package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/contracts/wasm/testcore/go/testcore"
)

func TestChainOwnerIDView(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		f := testcore.ScFuncs.TestChainOwnerIDView(ctx)
		f.Func.Call()
		require.NoError(t, ctx.Err)
		require.EqualValues(t, ctx.Originator().ScAgentID(), f.Results.ChainOwnerID().Value())
	})
}

func TestChainOwnerIDFull(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		f := testcore.ScFuncs.TestChainOwnerIDFull(ctx)
		f.Func.Post()
		require.NoError(t, ctx.Err)
		require.EqualValues(t, ctx.Originator().ScAgentID(), f.Results.ChainOwnerID().Value())
	})
}

func TestSandboxCall(t *testing.T) {
	run2(t, func(t *testing.T, w bool) {
		ctx := deployTestCore(t, w)

		f := testcore.ScFuncs.TestSandboxCall(ctx)
		f.Func.Call()
		require.NoError(t, ctx.Err)
	})
}
