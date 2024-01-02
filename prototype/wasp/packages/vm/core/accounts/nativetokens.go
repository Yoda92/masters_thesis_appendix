package accounts

import (
	"math/big"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
)

func nativeTokensMapKey(accountKey kv.Key) string {
	return PrefixNativeTokens + string(accountKey)
}

func nativeTokensMapR(state kv.KVStoreReader, accountKey kv.Key) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, nativeTokensMapKey(accountKey))
}

func nativeTokensMap(state kv.KVStore, accountKey kv.Key) *collections.Map {
	return collections.NewMap(state, nativeTokensMapKey(accountKey))
}

func getNativeTokenAmount(state kv.KVStoreReader, accountKey kv.Key, tokenID iotago.NativeTokenID) *big.Int {
	r := new(big.Int)
	b := nativeTokensMapR(state, accountKey).GetAt(tokenID[:])
	if len(b) > 0 {
		r.SetBytes(b)
	}
	return r
}

func setNativeTokenAmount(state kv.KVStore, accountKey kv.Key, tokenID iotago.NativeTokenID, n *big.Int) {
	if n.Sign() == 0 {
		nativeTokensMap(state, accountKey).DelAt(tokenID[:])
	} else {
		nativeTokensMap(state, accountKey).SetAt(tokenID[:], codec.EncodeBigIntAbs(n))
	}
}

func GetNativeTokenBalance(state kv.KVStoreReader, agentID isc.AgentID, nativeTokenID iotago.NativeTokenID, chainID isc.ChainID) *big.Int {
	return getNativeTokenAmount(state, accountKey(agentID, chainID), nativeTokenID)
}

func GetNativeTokenBalanceTotal(state kv.KVStoreReader, nativeTokenID iotago.NativeTokenID) *big.Int {
	return getNativeTokenAmount(state, l2TotalsAccount, nativeTokenID)
}

func GetNativeTokens(state kv.KVStoreReader, agentID isc.AgentID, chainID isc.ChainID) iotago.NativeTokens {
	ret := iotago.NativeTokens{}
	nativeTokensMapR(state, accountKey(agentID, chainID)).Iterate(func(idBytes []byte, val []byte) bool {
		ret = append(ret, &iotago.NativeToken{
			ID:     isc.MustNativeTokenIDFromBytes(idBytes),
			Amount: new(big.Int).SetBytes(val),
		})
		return true
	})
	return ret
}
