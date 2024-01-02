package examples

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	"github.com/iotaledger/wasp/packages/isc"
// 	"github.com/iotaledger/wasp/packages/solo"
// 	"github.com/iotaledger/wasp/packages/vm/core"
// )

// func TestExample1(t *testing.T) {
// 	env := solo.New(t, false, false)
// 	chain := env.NewChain(nil, "ex1")

// 	chainID, chainOwner, coreContracts := chain.GetInfo()              // calls view root::GetInfo
// 	require.EqualValues(t, len(corecontracts.All), len(coreContracts)) // 5 core contracts deployed by default

// 	t.Logf("chainID: %s", chainID.String())
// 	t.Logf("chain owner ID: %s", chainOwner.String())
// 	for hname, rec := range coreContracts {
// 		cid := isc.NewAgentID(chain.ChainID.AsAddress(), hname)
// 		t.Logf("    Core contract '%s': %s", rec.Name, cid)
// 	}
// }

// func TestExample2(t *testing.T) {
// 	env := solo.New(t, false, false)
// 	_, userAddress := env.NewKeyPair()
// 	t.Logf("Address of the userWallet is: %s", userAddress.String())
// 	numBaseTokens := env.L1NativeTokens(userAddress, colored.IOTA)
// 	t.Logf("balance of the userWallet is: %d iota", numBaseTokens)
// 	env.AssertAddressNativeTokenBalance(userAddress, colored.IOTA, 0)
// }

// func TestExample3(t *testing.T) {
// 	env := solo.New(t, false, false)
// 	_, userAddress := env.NewKeyPairWithFunds()
// 	t.Logf("Address of the userWallet is: %s", userAddress.String())
// 	numBaseTokens := env.L1NativeTokens(userAddress, colored.IOTA)
// 	t.Logf("balance of the userWallet is: %d iota", numBaseTokens)
// 	env.AssertAddressNativeTokenBalance(userAddress, colored.IOTA, utxodb.FundsFromFaucetAmount)
// }
