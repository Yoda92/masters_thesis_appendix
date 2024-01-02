package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/clients/chainclient"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/testutil/utxodb"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
)

func TestDepositWithdraw(t *testing.T) {
	e := setupWithNoChain(t)

	chain, err := e.Clu.DeployDefaultChain()
	require.NoError(t, err)

	chEnv := newChainEnv(t, e.Clu, chain)

	myWallet, myAddress, err := e.Clu.NewKeyPairWithFunds()
	require.NoError(e.t, err)

	require.True(t,
		e.Clu.AssertAddressBalances(myAddress, isc.NewAssetsBaseTokens(utxodb.FundsFromFaucetAmount)),
	)
	chEnv.checkLedger()

	myAgentID := isc.NewAgentID(myAddress)
	// origAgentID := isc.NewAgentID(chain.OriginatorAddress(), 0)

	// chEnv.checkBalanceOnChain(origAgentID, isc.BaseTokenID, 0)
	chEnv.checkBalanceOnChain(myAgentID, isc.BaseTokenID, 0)
	chEnv.checkLedger()

	// deposit some base tokens to the chain
	depositBaseTokens := 10 * isc.Million
	chClient := chainclient.New(e.Clu.L1Client(), e.Clu.WaspClient(0), chain.ChainID, myWallet)

	par := chainclient.NewPostRequestParams().WithBaseTokens(depositBaseTokens)
	reqTx, err := chClient.Post1Request(accounts.Contract.Hname(), accounts.FuncDeposit.Hname(), *par)
	require.NoError(t, err)

	receipts, err := chain.CommitteeMultiClient().WaitUntilAllRequestsProcessedSuccessfully(chain.ChainID, reqTx, true, 30*time.Second)
	require.NoError(t, err)
	chEnv.checkLedger()

	// chEnv.checkBalanceOnChain(origAgentID, isc.BaseTokenID, 0)
	gasFees1, err := iotago.DecodeUint64(receipts[0].GasFeeCharged)
	require.NoError(t, err)

	onChainBalance := depositBaseTokens - gasFees1
	chEnv.checkBalanceOnChain(myAgentID, isc.BaseTokenID, onChainBalance)

	require.True(t,
		e.Clu.AssertAddressBalances(myAddress, isc.NewAssetsBaseTokens(utxodb.FundsFromFaucetAmount-depositBaseTokens)),
	)

	// withdraw some base tokens back
	baseTokensToWithdraw := 1 * isc.Million
	req, err := chClient.PostOffLedgerRequest(context.Background(), accounts.Contract.Hname(), accounts.FuncWithdraw.Hname(),
		chainclient.PostRequestParams{
			Allowance: isc.NewAssetsBaseTokens(baseTokensToWithdraw),
		},
	)
	require.NoError(t, err)
	receipt, err := chain.CommitteeMultiClient().WaitUntilRequestProcessedSuccessfully(chain.ChainID, req.ID(), true, 30*time.Second)
	require.NoError(t, err)

	chEnv.checkLedger()
	gasFees2, err := iotago.DecodeUint64(receipt.GasFeeCharged)
	require.NoError(t, err)

	chEnv.checkBalanceOnChain(myAgentID, isc.BaseTokenID, onChainBalance-baseTokensToWithdraw-gasFees2)
	require.True(t,
		e.Clu.AssertAddressBalances(myAddress, isc.NewAssetsBaseTokens(utxodb.FundsFromFaucetAmount-depositBaseTokens+baseTokensToWithdraw)),
	)

	// TODO use "withdraw all base tokens" entrypoint to withdraw all remaining base tokens
}
