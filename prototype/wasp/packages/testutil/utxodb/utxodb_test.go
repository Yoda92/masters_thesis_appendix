package utxodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/builder"
	"github.com/iotaledger/iota.go/v3/tpkg"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/parameters"
)

func TestBasic(t *testing.T) {
	u := New()
	require.EqualValues(t, u.Supply(), u.GetAddressBalanceBaseTokens(u.GenesisAddress()))
	gtime := u.GlobalTime()
	expectedTime := time.Unix(1, 0).Add(1 * time.Millisecond)
	require.EqualValues(t, expectedTime, gtime)

	u.AdvanceClockBy(10 * time.Second)
	gtime1 := u.GlobalTime()
	expectedTime = gtime.Add(10 * time.Second)
	require.EqualValues(t, expectedTime, gtime1)
}

func TestRequestFunds(t *testing.T) {
	u := New()
	addr := tpkg.RandEd25519Address()
	tx, err := u.GetFundsFromFaucet(addr)
	require.NoError(t, err)
	require.EqualValues(t, u.Supply()-FundsFromFaucetAmount, u.GetAddressBalanceBaseTokens(u.GenesisAddress()))
	require.EqualValues(t, FundsFromFaucetAmount, u.GetAddressBalanceBaseTokens(addr))

	txID, err := tx.ID()
	require.NoError(t, err)
	require.Same(t, tx, u.MustGetTransaction(txID))

	gtime := u.GlobalTime()
	expectedTime := time.Unix(1, 0).Add(2 * time.Millisecond)
	require.EqualValues(t, expectedTime, gtime)
}

func TestAddTransactionFail(t *testing.T) {
	u := New()

	addr := tpkg.RandEd25519Address()
	tx, err := u.GetFundsFromFaucet(addr)
	require.NoError(t, err)

	err = u.AddToLedger(tx)
	require.Error(t, err)
}

func TestDoubleSpend(t *testing.T) {
	keyPair1 := cryptolib.NewKeyPair()

	addr1 := keyPair1.GetPublicKey().AsEd25519Address()
	key1Signer := iotago.NewInMemoryAddressSigner(keyPair1.GetPrivateKey().AddressKeysForEd25519Address(addr1))

	addr2 := tpkg.RandEd25519Address()
	addr3 := tpkg.RandEd25519Address()

	u := New()

	tx1, err := u.GetFundsFromFaucet(addr1)
	require.NoError(t, err)
	tx1ID, err := tx1.ID()
	require.NoError(t, err)

	spend2, err := builder.NewTransactionBuilder(tpkg.TestNetworkID).
		AddInput(&builder.TxInput{
			UnlockTarget: addr1,
			Input:        tx1.Essence.Outputs[0],
			InputID:      iotago.OutputIDFromTransactionIDAndIndex(tx1ID, 0),
		}).
		AddOutput(&iotago.BasicOutput{
			Amount: FundsFromFaucetAmount,
			Conditions: iotago.UnlockConditions{
				&iotago.AddressUnlockCondition{Address: addr2},
			},
		}).
		Build(parameters.L1().Protocol, key1Signer)
	require.NoError(t, err)
	err = u.AddToLedger(spend2)
	require.NoError(t, err)

	spend3, err := builder.NewTransactionBuilder(tpkg.TestNetworkID).
		AddInput(&builder.TxInput{
			UnlockTarget: addr1,
			Input:        tx1.Essence.Outputs[0],
			InputID:      iotago.OutputIDFromTransactionIDAndIndex(tx1ID, 0),
		}).
		AddOutput(&iotago.BasicOutput{
			Amount: FundsFromFaucetAmount,
			Conditions: iotago.UnlockConditions{
				&iotago.AddressUnlockCondition{Address: addr3},
			},
		}).
		Build(parameters.L1().Protocol, key1Signer)
	require.NoError(t, err)
	err = u.AddToLedger(spend3)
	require.Error(t, err)
}

func TestGetOutput(t *testing.T) {
	u := New()
	addr := tpkg.RandEd25519Address()
	tx, err := u.GetFundsFromFaucet(addr)
	require.NoError(t, err)

	txID, err := tx.ID()
	require.NoError(t, err)

	outid0 := iotago.OutputIDFromTransactionIDAndIndex(txID, 0)
	out0 := u.GetOutput(outid0)
	require.EqualValues(t, FundsFromFaucetAmount, out0.Deposit())

	outid1 := iotago.OutputIDFromTransactionIDAndIndex(txID, 1)
	out1 := u.GetOutput(outid1)
	require.EqualValues(t, u.Supply()-FundsFromFaucetAmount, out1.Deposit())

	outidFail := iotago.OutputIDFromTransactionIDAndIndex(txID, 5)
	require.Nil(t, u.GetOutput(outidFail))
}
