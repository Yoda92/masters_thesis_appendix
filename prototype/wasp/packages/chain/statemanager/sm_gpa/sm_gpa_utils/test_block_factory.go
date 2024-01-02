// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package sm_gpa_utils

import (
	"crypto/rand"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iotaledger/hive.go/kvstore/mapdb"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/tpkg"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/origin"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/testutil"
)

type BlockFactory struct {
	t                   require.TestingT
	store               state.Store
	chainID             isc.ChainID
	chainInitParams     dict.Dict
	lastBlockCommitment *state.L1Commitment
	aliasOutputs        map[state.BlockHash]*isc.AliasOutputWithID
}

func NewBlockFactory(t require.TestingT, chainInitParamsOpt ...dict.Dict) *BlockFactory {
	var chainInitParams dict.Dict
	if len(chainInitParamsOpt) > 0 {
		chainInitParams = chainInitParamsOpt[0]
	} else {
		chainInitParams = nil
	}
	aliasOutput0ID := iotago.OutputIDFromTransactionIDAndIndex(getRandomTxID(t), 0)
	chainID := isc.ChainIDFromAliasID(iotago.AliasIDFromOutputID(aliasOutput0ID))
	stateAddress := cryptolib.NewKeyPair().GetPublicKey().AsEd25519Address()
	originCommitment := origin.L1Commitment(chainInitParams, 0)
	aliasOutput0 := &iotago.AliasOutput{
		Amount:        tpkg.TestTokenSupply,
		AliasID:       chainID.AsAliasID(), // NOTE: not very correct: origin output's AliasID should be empty; left here to make mocking transitions easier
		StateMetadata: testutil.DummyStateMetadata(originCommitment).Bytes(),
		Conditions: iotago.UnlockConditions{
			&iotago.StateControllerAddressUnlockCondition{Address: stateAddress},
			&iotago.GovernorAddressUnlockCondition{Address: stateAddress},
		},
		Features: iotago.Features{
			&iotago.SenderFeature{
				Address: stateAddress,
			},
		},
	}
	aliasOutputs := make(map[state.BlockHash]*isc.AliasOutputWithID)
	originOutput := isc.NewAliasOutputWithID(aliasOutput0, aliasOutput0ID)
	aliasOutputs[originCommitment.BlockHash()] = originOutput
	chainStore := state.NewStoreWithUniqueWriteMutex(mapdb.NewMapDB())
	origin.InitChain(chainStore, chainInitParams, 0)
	return &BlockFactory{
		t:                   t,
		store:               chainStore,
		chainID:             chainID,
		chainInitParams:     chainInitParams,
		lastBlockCommitment: originCommitment,
		aliasOutputs:        aliasOutputs,
	}
}

func (bfT *BlockFactory) GetChainID() isc.ChainID {
	return bfT.chainID
}

func (bfT *BlockFactory) GetChainInitParameters() dict.Dict {
	return bfT.chainInitParams
}

func (bfT *BlockFactory) GetOriginOutput() *isc.AliasOutputWithID {
	return bfT.GetAliasOutput(origin.L1Commitment(nil, 0))
}

func (bfT *BlockFactory) GetBlocks(
	count,
	branchingFactor int,
) []state.Block {
	blocks := bfT.GetBlocksFrom(count, branchingFactor, bfT.lastBlockCommitment)
	require.Equal(bfT.t, count, len(blocks))
	bfT.lastBlockCommitment = blocks[count-1].L1Commitment()
	return blocks
}

func (bfT *BlockFactory) GetBlocksFrom(
	count,
	branchingFactor int,
	commitment *state.L1Commitment,
	incrementFactorOpt ...uint64,
) []state.Block {
	var incrementFactor uint64
	if len(incrementFactorOpt) > 0 {
		incrementFactor = incrementFactorOpt[0]
	} else {
		incrementFactor = 1
	}
	result := make([]state.Block, count+1)
	var err error
	result[0], err = bfT.store.BlockByTrieRoot(commitment.TrieRoot())
	require.NoError(bfT.t, err)
	for i := 1; i < len(result); i++ {
		baseIndex := (i + branchingFactor - 2) / branchingFactor
		increment := uint64(1+i%branchingFactor) * incrementFactor
		result[i] = bfT.GetNextBlock(result[baseIndex].L1Commitment(), increment)
	}
	return result[1:]
}

func (bfT *BlockFactory) GetNextBlock(
	commitment *state.L1Commitment,
	incrementOpt ...uint64,
) state.Block {
	stateDraft, err := bfT.store.NewStateDraft(time.Now(), commitment)
	require.NoError(bfT.t, err)
	counterKey := kv.Key(coreutil.StateVarBlockIndex + "counter")
	counterBin := stateDraft.Get(counterKey)
	counter, err := codec.DecodeUint64(counterBin, 0)
	require.NoError(bfT.t, err)
	var increment uint64
	if len(incrementOpt) > 0 {
		increment = incrementOpt[0]
	} else {
		increment = 1
	}
	counterBin = codec.EncodeUint64(counter + increment)
	stateDraft.Mutations().Set(counterKey, counterBin)
	block := bfT.store.Commit(stateDraft)
	// require.EqualValues(t, stateDraft.BlockIndex(), block.BlockIndex())
	newCommitment := block.L1Commitment()

	consumedAliasOutput := bfT.GetAliasOutput(commitment).GetAliasOutput()

	aliasOutput := &iotago.AliasOutput{
		Amount:         consumedAliasOutput.Amount,
		NativeTokens:   consumedAliasOutput.NativeTokens,
		AliasID:        consumedAliasOutput.AliasID,
		StateIndex:     consumedAliasOutput.StateIndex + 1,
		StateMetadata:  testutil.DummyStateMetadata(newCommitment).Bytes(),
		FoundryCounter: consumedAliasOutput.FoundryCounter,
		Conditions:     consumedAliasOutput.Conditions,
		Features:       consumedAliasOutput.Features,
	}
	aliasOutputID := iotago.OutputIDFromTransactionIDAndIndex(getRandomTxID(bfT.t), 0)
	aliasOutputWithID := isc.NewAliasOutputWithID(aliasOutput, aliasOutputID)
	bfT.aliasOutputs[newCommitment.BlockHash()] = aliasOutputWithID

	return block
}

func (bfT *BlockFactory) GetStore() state.Store {
	return NewReadOnlyStore(bfT.store)
}

func (bfT *BlockFactory) GetStateDraft(block state.Block) state.StateDraft {
	result, err := bfT.store.NewEmptyStateDraft(block.PreviousL1Commitment())
	require.NoError(bfT.t, err)
	block.Mutations().ApplyTo(result)
	return result
}

func (bfT *BlockFactory) GetAliasOutput(commitment *state.L1Commitment) *isc.AliasOutputWithID {
	result, ok := bfT.aliasOutputs[commitment.BlockHash()]
	require.True(bfT.t, ok)
	return result
}

func getRandomTxID(t require.TestingT) iotago.TransactionID {
	var result iotago.TransactionID
	_, err := rand.Read(result[:])
	require.NoError(t, err)
	return result
}
