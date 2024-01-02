// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmsolo

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

type SoloFoundry struct {
	sn            uint32
	nativeTokenID iotago.NativeTokenID
	agent         *SoloAgent
	ctx           *SoloContext
}

func NewSoloFoundry(ctx *SoloContext, maxSupply interface{}, agent ...*SoloAgent) (sf *SoloFoundry, err error) {
	sf = &SoloFoundry{ctx: ctx}
	fp := ctx.Chain.NewFoundryParams(cvt.ToBigInt(maxSupply))
	if len(agent) == 1 {
		sf.agent = agent[0]
		fp.WithUser(sf.agent.Pair)
	}
	sf.sn, sf.nativeTokenID, err = fp.CreateFoundry()
	if err != nil {
		return nil, err
	}
	return sf, nil
}

func (sf *SoloFoundry) Destroy() error {
	return sf.ctx.Chain.DestroyFoundry(sf.sn, sf.agent.Pair)
}

func (sf *SoloFoundry) DestroyTokens(amount interface{}) error {
	return sf.ctx.Chain.DestroyTokensOnL2(sf.nativeTokenID, amount, sf.agent.Pair)
}

func (sf *SoloFoundry) Mint(amount interface{}) error {
	return sf.ctx.Chain.MintTokens(sf.sn, cvt.ToBigInt(amount), sf.agent.Pair)
}

func (sf *SoloFoundry) SN() uint32 {
	return sf.sn
}

func (sf *SoloFoundry) TokenID() wasmtypes.ScTokenID {
	return cvt.ScTokenID(sf.nativeTokenID)
}
