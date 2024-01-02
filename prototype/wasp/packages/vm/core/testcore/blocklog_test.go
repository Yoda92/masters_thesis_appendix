package testcore

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/corecontracts"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
)

func TestBlockInfoLatest(t *testing.T) {
	corecontracts.PrintWellKnownHnames()
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	chain := env.NewChain()

	bi := chain.GetLatestBlockInfo()
	require.NotNil(t, bi)
	require.EqualValues(t, 1, bi.BlockIndex())
	require.EqualValues(t, 1, bi.TotalRequests)
	require.EqualValues(t, 1, bi.NumSuccessfulRequests)
	require.EqualValues(t, 0, bi.NumOffLedgerRequests)
	t.Logf("%s", bi.String())
}

func TestBlockInfo(t *testing.T) {
	corecontracts.PrintWellKnownHnames()
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	chain := env.NewChain()

	bi, err := chain.GetBlockInfo(0)
	require.NoError(t, err)
	require.NotNil(t, bi)
	require.EqualValues(t, 0, bi.BlockIndex())
	require.EqualValues(t, 1, bi.TotalRequests)
	require.EqualValues(t, 1, bi.NumSuccessfulRequests)
	require.EqualValues(t, 0, bi.NumOffLedgerRequests)
	t.Logf("%s", bi.String())

	bi, err = chain.GetBlockInfo(1)
	require.NoError(t, err)
	require.NotNil(t, bi)
	require.EqualValues(t, 1, bi.BlockIndex())
	require.EqualValues(t, 1, bi.TotalRequests)
	require.EqualValues(t, 1, bi.NumSuccessfulRequests)
	require.EqualValues(t, 0, bi.NumOffLedgerRequests)
	t.Logf("%s", bi.String())
}

func TestBlockInfoLatestWithRequest(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})

	ch := env.NewChain()

	err := ch.DepositBaseTokensToL2(100_000, nil)
	require.NoError(t, err)

	bi := ch.GetLatestBlockInfo()
	t.Logf("after ch deployment:\n%s", bi.String())
	// uploading one blob
	_, err = ch.UploadBlob(nil, "field", "dummy blob data")
	require.NoError(t, err)

	bi = ch.GetLatestBlockInfo()
	require.NotNil(t, bi)
	require.EqualValues(t, 3, bi.BlockIndex())
	require.EqualValues(t, 1, bi.TotalRequests)
	require.EqualValues(t, 1, bi.NumSuccessfulRequests)
	require.EqualValues(t, 1, bi.NumOffLedgerRequests)
	t.Logf("%s", bi.String())
}

func TestBlockInfoSeveral(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	ch := env.NewChain()

	err := ch.DepositBaseTokensToL2(100_000, nil)
	require.NoError(t, err)

	const numReqs = 5
	for i := 0; i < numReqs; i++ {
		_, err := ch.UploadBlob(nil, "field", fmt.Sprintf("dummy blob data #%d", i))
		require.NoError(t, err)
	}

	bi := ch.GetLatestBlockInfo()
	require.EqualValues(t, 2+numReqs, int(bi.BlockIndex()))

	for blockIndex := uint32(0); blockIndex <= bi.BlockIndex(); blockIndex++ {
		bi1, err := ch.GetBlockInfo(blockIndex)
		require.NoError(t, err)
		require.NotNil(t, bi1)
		t.Logf("%s", bi1.String())
		require.EqualValues(t, blockIndex, bi1.BlockIndex())
		require.EqualValues(t, 1, bi1.TotalRequests)
		require.EqualValues(t, 1, bi1.NumSuccessfulRequests)
		require.LessOrEqual(t, bi1.NumOffLedgerRequests, bi1.TotalRequests)
	}
}

func TestRequestIsProcessed(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	ch := env.NewChain()

	ch.MustDepositBaseTokensToL2(10_000, nil)

	req := solo.NewCallParams(governance.Contract.Name, governance.FuncSetMetadata.Name,
		governance.ParamPublicURL, []byte("foo"),
	).
		WithGasBudget(100_000)
	tx, _, err := ch.PostRequestSyncTx(req, nil)
	require.NoError(t, err)

	reqs, err := env.RequestsForChain(tx, ch.ChainID)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(reqs))

	bi := ch.GetLatestBlockInfo()
	require.NoError(t, err)
	require.True(t, ch.IsRequestProcessed(reqs[0].ID()))
	t.Logf("%s", bi.String())
}

func TestRequestReceipt(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	ch := env.NewChain()

	ch.MustDepositBaseTokensToL2(10_000, nil)

	req := solo.NewCallParams(governance.Contract.Name, governance.FuncSetMetadata.Name,
		governance.ParamPublicURL, []byte("foo"),
	).
		WithGasBudget(100_000)
	tx, _, err := ch.PostRequestSyncTx(req, nil)
	require.NoError(t, err)

	reqs, err := env.RequestsForChain(tx, ch.ChainID)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(reqs))
	require.True(t, ch.IsRequestProcessed(reqs[0].ID()))

	receipt, err := ch.GetRequestReceipt(reqs[0].ID())
	require.NoError(t, err)
	a := reqs[0].Bytes()
	b := receipt.Request.Bytes()
	require.Equal(t, a, b)
	require.Nil(t, receipt.Error)
	require.EqualValues(t, 3, int(receipt.BlockIndex))
	require.EqualValues(t, 0, receipt.RequestIndex)
	t.Logf("%s", receipt.String())
}

func TestRequestReceiptsForBlocks(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	ch := env.NewChain()

	ch.MustDepositBaseTokensToL2(10_000, nil)

	req := solo.NewCallParams(governance.Contract.Name, governance.FuncSetMetadata.Name,
		governance.ParamPublicURL, []byte("foo"),
	).
		WithGasBudget(100_000)
	tx, _, err := ch.PostRequestSyncTx(req, nil)
	require.NoError(t, err)

	reqs, err := env.RequestsForChain(tx, ch.ChainID)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(reqs))

	require.True(t, ch.IsRequestProcessed(reqs[0].ID()))

	recs := ch.GetRequestReceiptsForBlock(3)
	require.EqualValues(t, 1, len(recs))
	require.EqualValues(t, reqs[0].ID(), recs[0].Request.ID())
	t.Logf("%s\n", recs[0].String())
}

func TestRequestIDsForBlocks(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	ch := env.NewChain()

	ch.MustDepositBaseTokensToL2(10_000, nil)

	req := solo.NewCallParams(governance.Contract.Name, governance.FuncSetMetadata.Name,
		governance.ParamPublicURL, []byte("foo"),
	).
		WithGasBudget(100_000)
	tx, _, err := ch.PostRequestSyncTx(req, nil)
	require.NoError(t, err)

	reqs, err := env.RequestsForChain(tx, ch.ChainID)
	require.NoError(t, err)
	require.EqualValues(t, 1, len(reqs))

	require.True(t, ch.IsRequestProcessed(reqs[0].ID()))

	ids := ch.GetRequestIDsForBlock(3)
	require.EqualValues(t, 1, len(ids))
	require.EqualValues(t, reqs[0].ID(), ids[0])
}

func TestViewGetRequestReceipt(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true})
	ch := env.NewChain()
	// try to get a receipt for a request that does not exist
	receipt, err := ch.GetRequestReceipt(isc.RequestID{})
	require.Nil(t, receipt)
	require.NoError(t, err)
}

func TestBlocklogPruning(t *testing.T) {
	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true, Debug: true})
	ch, _ := env.NewChainExt(nil, 10*isc.Million, "chain1", dict.Dict{
		origin.ParamBlockKeepAmount: codec.EncodeInt32(10),
	})
	for i := 1; i <= 20; i++ {
		ch.DepositBaseTokensToL2(1000, nil)
	}
	// at this point blocks 0..10 have been pruned, and blocks 11..20 are available
	require.EqualValues(t, 20, ch.LatestBlock().StateIndex())
	require.EqualValues(t, 20, ch.EVM().BlockNumber().Uint64())
	for i := uint32(0); i <= 10; i++ {
		_, err := ch.GetBlockInfo(i)
		require.ErrorContains(t, err, "not found")
		_, err = ch.EVM().BlockByNumber(big.NewInt(int64(i)))
		require.ErrorContains(t, err, "not found")
	}
	for i := uint32(11); i <= 20; i++ {
		bi, err := ch.GetBlockInfo(i)
		require.NoError(t, err)
		require.EqualValues(t, i, bi.BlockIndex())
		evmBlock, err := ch.EVM().BlockByNumber(big.NewInt(int64(i)))
		require.NoError(t, err)
		require.EqualValues(t, i, evmBlock.Number().Uint64())
	}
}

func TestBlocklogFoundriesWithPruning(t *testing.T) {
	// test that foundries can be accessed even after the block is pruned

	env := solo.New(t, &solo.InitOptions{AutoAdjustStorageDeposit: true, Debug: true})
	ch, _ := env.NewChainExt(nil, 10*isc.Million, "chain1", dict.Dict{
		origin.ParamBlockKeepAmount: codec.EncodeInt32(10),
	})
	ch.DepositBaseTokensToL2(1*isc.Million, nil)

	sn, _, err := ch.NewFoundryParams(10).CreateFoundry()
	require.NoError(t, err)

	// provoke the block where the foundry was stored to be pruned
	for i := 1; i <= 20; i++ {
		ch.DepositBaseTokensToL2(1000, nil)
	}

	err = ch.DestroyFoundry(sn, ch.OriginatorPrivateKey)
	require.NoError(t, err)
}
