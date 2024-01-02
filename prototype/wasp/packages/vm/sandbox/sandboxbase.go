// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package sandbox

import (
	"math/big"
	"time"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/assert"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/execution"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

type SandboxBase struct {
	Ctx       execution.WaspCallContext
	assertObj *assert.Assert
}

var _ isc.SandboxBase = &SandboxBase{}

func (s *SandboxBase) assert() *assert.Assert {
	if s.assertObj == nil {
		s.assertObj = assert.NewAssert(s.Ctx)
	}
	return s.assertObj
}

func (s *SandboxBase) AccountID() isc.AgentID {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.CurrentContractAccountID()
}

func (s *SandboxBase) Caller() isc.AgentID {
	s.Ctx.GasBurn(gas.BurnCodeGetCallerData)
	return s.Ctx.Caller()
}

func (s *SandboxBase) BalanceBaseTokens() uint64 {
	s.Ctx.GasBurn(gas.BurnCodeGetBalance)
	return s.Ctx.GetBaseTokensBalance(s.AccountID())
}

func (s *SandboxBase) BalanceNativeToken(nativeTokenID iotago.NativeTokenID) *big.Int {
	s.Ctx.GasBurn(gas.BurnCodeGetBalance)
	return s.Ctx.GetNativeTokenBalance(s.AccountID(), nativeTokenID)
}

func (s *SandboxBase) BalanceNativeTokens() iotago.NativeTokens {
	s.Ctx.GasBurn(gas.BurnCodeGetBalance)
	return s.Ctx.GetNativeTokens(s.AccountID())
}

func (s *SandboxBase) OwnedNFTs() []iotago.NFTID {
	s.Ctx.GasBurn(gas.BurnCodeGetBalance)
	return s.Ctx.GetAccountNFTs(s.AccountID())
}

func (s *SandboxBase) HasInAccount(agentID isc.AgentID, assets *isc.Assets) bool {
	s.Ctx.GasBurn(gas.BurnCodeGetBalance)
	accountAssets := isc.Assets{
		BaseTokens:   s.Ctx.GetBaseTokensBalance(agentID),
		NativeTokens: s.Ctx.GetNativeTokens(agentID),
		NFTs:         s.Ctx.GetAccountNFTs(agentID),
	}
	return accountAssets.Spend(assets)
}

func (s *SandboxBase) GetNFTData(nftID iotago.NFTID) *isc.NFT {
	s.Ctx.GasBurn(gas.BurnCodeGetNFTData)
	return s.Ctx.GetNFTData(nftID)
}

func (s *SandboxBase) ChainID() isc.ChainID {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.ChainID()
}

func (s *SandboxBase) ChainOwnerID() isc.AgentID {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.ChainOwnerID()
}

func (s *SandboxBase) ChainInfo() *isc.ChainInfo {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.ChainInfo()
}

func (s *SandboxBase) Contract() isc.Hname {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.CurrentContractHname()
}

func (s *SandboxBase) Timestamp() time.Time {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.Timestamp()
}

func (s *SandboxBase) Log() isc.LogInterface {
	// TODO should Log be disabled for wasm contracts? not much of a point in exposing internal logging
	return s.Ctx
}

func (s *SandboxBase) Params() *isc.Params {
	s.Ctx.GasBurn(gas.BurnCodeGetContext)
	return s.Ctx.Params()
}

func (s *SandboxBase) Utils() isc.Utils {
	return NewUtils(s.Gas())
}

func (s *SandboxBase) Gas() isc.Gas {
	return s
}

func (s *SandboxBase) Burned() uint64 {
	return s.Ctx.GasBurned()
}

func (s *SandboxBase) Burn(burnCode gas.BurnCode, par ...uint64) {
	s.Ctx.GasBurn(burnCode, par...)
}

func (s *SandboxBase) Budget() uint64 {
	return s.Ctx.GasBudgetLeft()
}

func (s *SandboxBase) EstimateGasMode() bool {
	return s.Ctx.GasEstimateMode()
}

// -- helper methods
func (s *SandboxBase) Requiref(cond bool, format string, args ...interface{}) {
	s.assert().Requiref(cond, format, args...)
}

func (s *SandboxBase) RequireNoError(err error, str ...string) {
	s.assert().RequireNoError(err, str...)
}

func (s *SandboxBase) CallView(contractHname, entryPoint isc.Hname, params dict.Dict) dict.Dict {
	s.Ctx.GasBurn(gas.BurnCodeCallContract)
	if params == nil {
		params = make(dict.Dict)
	}
	return s.Ctx.Call(contractHname, entryPoint, params, nil)
}

func (s *SandboxBase) StateR() kv.KVStoreReader {
	return s.Ctx.ContractStateReaderWithGasBurn()
}
