package inccounter

import (
	"fmt"
	"math"
	"time"

	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"
)

var Contract = coreutil.NewContract("inccounter")

var Processor = Contract.Processor(initialize,
	FuncIncCounter.WithHandler(incCounter),
	FuncIncAndRepeatOnceAfter2s.WithHandler(incCounterAndRepeatOnce),
	FuncIncAndRepeatMany.WithHandler(incCounterAndRepeatMany),
	FuncSpawn.WithHandler(spawn),
	ViewGetCounter.WithHandler(getCounter),
)

var (
	FuncIncCounter              = coreutil.Func("incCounter")
	FuncIncAndRepeatOnceAfter2s = coreutil.Func("incAndRepeatOnceAfter5s")
	FuncIncAndRepeatMany        = coreutil.Func("incAndRepeatMany")
	FuncSpawn                   = coreutil.Func("spawn")
	ViewGetCounter              = coreutil.ViewFunc("getCounter")
)

const (
	VarNumRepeats = "numRepeats"
	VarCounter    = "counter"
	VarName       = "name"
)

func initialize(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("inccounter.init in %s", ctx.Contract().String())
	params := ctx.Params()
	val := codec.MustDecodeInt64(params.Get(VarCounter), 0)
	ctx.State().Set(VarCounter, codec.EncodeInt64(val))
	eventCounter(ctx, val)
	return nil
}

func incCounter(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("inccounter.incCounter in %s", ctx.Contract().String())
	params := ctx.Params()
	inc := params.MustGetInt64(VarCounter, 1)

	state := kvdecoder.New(ctx.State(), ctx.Log())
	val := state.MustGetInt64(VarCounter, 0)
	ctx.Log().Infof("incCounter: increasing counter value %d by %d, anchor index: #%d",
		val, inc, ctx.StateAnchor().StateIndex)
	tra := "(empty)"
	if ctx.AllowanceAvailable() != nil {
		tra = ctx.AllowanceAvailable().String()
	}
	ctx.Log().Infof("incCounter: allowance available: %s", tra)
	ctx.State().Set(VarCounter, codec.EncodeInt64(val+inc))
	eventCounter(ctx, val+inc)
	return nil
}

func incCounterAndRepeatOnce(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("inccounter.incCounterAndRepeatOnce")
	state := ctx.State()
	val := codec.MustDecodeInt64(state.Get(VarCounter), 0)

	ctx.Log().Debugf(fmt.Sprintf("incCounterAndRepeatOnce: increasing counter value: %d", val))
	state.Set(VarCounter, codec.EncodeInt64(val+1))
	eventCounter(ctx, val+1)
	allowance := ctx.AllowanceAvailable()
	ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Send(isc.RequestParameters{
		TargetAddress:                 ctx.ChainID().AsAddress(),
		Assets:                        isc.NewAssets(allowance.BaseTokens, nil),
		AdjustToMinimumStorageDeposit: true,
		Metadata: &isc.SendMetadata{
			TargetContract: ctx.Contract(),
			EntryPoint:     FuncIncCounter.Hname(),
			GasBudget:      math.MaxUint64,
		},
		Options: isc.SendOptions{
			Timelock: ctx.Timestamp().Add(2 * time.Second),
		},
	})
	ctx.Log().Debugf("incCounterAndRepeatOnce: PostRequestToSelfWithDelay RequestInc 2 sec")
	return nil
}

func incCounterAndRepeatMany(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("inccounter.incCounterAndRepeatMany")

	state := ctx.State()
	params := ctx.Params()

	val := codec.MustDecodeInt64(state.Get(VarCounter), 0)

	state.Set(VarCounter, codec.EncodeInt64(val+1))
	eventCounter(ctx, val+1)
	ctx.Log().Debugf("inccounter.incCounterAndRepeatMany: increasing counter value: %d", val)

	var numRepeats int64
	if params.Has(VarNumRepeats) {
		numRepeats = codec.MustDecodeInt64(params.Get(VarNumRepeats), 0)
	} else {
		numRepeats = codec.MustDecodeInt64(state.Get(VarNumRepeats), 0)
	}
	if numRepeats == 0 {
		ctx.Log().Debugf("inccounter.incCounterAndRepeatMany: finished chain of requests. counter value: %d", val)
		return nil
	}

	ctx.Log().Debugf("chain of %d requests ahead", numRepeats)

	state.Set(VarNumRepeats, codec.EncodeInt64(numRepeats-1))
	ctx.TransferAllowedFunds(ctx.AccountID())
	ctx.Send(isc.RequestParameters{
		TargetAddress:                 ctx.ChainID().AsAddress(),
		Assets:                        isc.NewAssets(1000, nil),
		AdjustToMinimumStorageDeposit: true,
		Metadata: &isc.SendMetadata{
			TargetContract: ctx.Contract(),
			EntryPoint:     FuncIncAndRepeatMany.Hname(),
			GasBudget:      math.MaxUint64,
			Allowance:      isc.NewAssetsBaseTokens(1000),
		},
		Options: isc.SendOptions{
			Timelock: ctx.Timestamp().Add(2 * time.Second),
		},
	})

	ctx.Log().Debugf("incCounterAndRepeatMany. remaining repeats = %d", numRepeats-1)
	return nil
}

// spawn deploys new contract and calls it
func spawn(ctx isc.Sandbox) dict.Dict {
	ctx.Log().Debugf("inccounter.spawn")

	state := kvdecoder.New(ctx.State(), ctx.Log())
	val := state.MustGetInt64(VarCounter)
	params := ctx.Params()
	name := params.MustGetString(VarName)

	callPar := dict.New()
	callPar.Set(VarCounter, codec.EncodeInt64(val+1))
	eventCounter(ctx, val+1)
	ctx.DeployContract(Contract.ProgramHash, name, callPar)

	// increase counter in newly spawned contract
	hname := isc.Hn(name)
	ctx.Call(hname, FuncIncCounter.Hname(), nil, nil)

	return nil
}

func getCounter(ctx isc.SandboxView) dict.Dict {
	state := ctx.StateR()
	val := codec.MustDecodeInt64(state.Get(VarCounter), 0)
	return dict.Dict{VarCounter: codec.EncodeInt64(val)}
}
