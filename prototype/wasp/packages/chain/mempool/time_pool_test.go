// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package mempool_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"

	"github.com/iotaledger/wasp/packages/chain/mempool"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/testutil/testlogger"
	"github.com/iotaledger/wasp/packages/vm/core/governance"
	"github.com/iotaledger/wasp/packages/vm/gas"
)

func TestTimePoolBasic(t *testing.T) {
	log := testlogger.NewLogger(t)
	kp := cryptolib.NewKeyPair()
	tp := mempool.NewTimePool(func(i int) {}, log)
	t0 := time.Now()
	t1 := t0.Add(17 * time.Nanosecond)
	t2 := t0.Add(17 * time.Minute)
	t3 := t0.Add(17 * time.Hour)
	r0 := isc.NewOffLedgerRequest(isc.RandomChainID(), governance.Contract.Hname(), governance.FuncAddCandidateNode.Hname(), nil, 0, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	r1 := isc.NewOffLedgerRequest(isc.RandomChainID(), governance.Contract.Hname(), governance.FuncAddCandidateNode.Hname(), nil, 1, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	r2 := isc.NewOffLedgerRequest(isc.RandomChainID(), governance.Contract.Hname(), governance.FuncAddCandidateNode.Hname(), nil, 2, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	r3 := isc.NewOffLedgerRequest(isc.RandomChainID(), governance.Contract.Hname(), governance.FuncAddCandidateNode.Hname(), nil, 3, gas.LimitsDefault.MaxGasPerRequest).Sign(kp)
	require.False(t, tp.Has(isc.RequestRefFromRequest(r0)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r1)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r2)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r3)))
	tp.AddRequest(t0, r0)
	tp.AddRequest(t1, r1)
	tp.AddRequest(t2, r2)
	tp.AddRequest(t3, r3)
	require.True(t, tp.Has(isc.RequestRefFromRequest(r0)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r1)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r2)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r3)))

	var taken []isc.Request

	taken = tp.TakeTill(t0)
	require.Len(t, taken, 1)
	require.Equal(t, r0, taken[0])
	require.False(t, tp.Has(isc.RequestRefFromRequest(r0)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r1)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r2)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r3)))

	taken = tp.TakeTill(t0)
	require.Len(t, taken, 0)

	taken = tp.TakeTill(t0.Add(30 * time.Minute))
	require.Len(t, taken, 2)
	require.Contains(t, taken, r1)
	require.Contains(t, taken, r2)
	require.False(t, tp.Has(isc.RequestRefFromRequest(r0)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r1)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r2)))
	require.True(t, tp.Has(isc.RequestRefFromRequest(r3)))

	taken = tp.TakeTill(t0.Add(30 * time.Hour))
	require.Len(t, taken, 1)
	require.Contains(t, taken, r3)
	require.False(t, tp.Has(isc.RequestRefFromRequest(r0)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r1)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r2)))
	require.False(t, tp.Has(isc.RequestRefFromRequest(r3)))
}

func TestTimePoolRapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		sm := newtimePoolSM(t)
		t.Repeat(rapid.StateMachineActions(sm))
	})
}

type timePoolSM struct {
	tp    mempool.TimePool
	kp    *cryptolib.KeyPair
	added int
	taken int
}

var _ rapid.StateMachine = &timePoolSM{}

func newtimePoolSM(t *rapid.T) *timePoolSM {
	sm := new(timePoolSM)
	log := testlogger.NewLogger(t)
	sm.tp = mempool.NewTimePool(func(i int) {}, log)
	sm.kp = cryptolib.NewKeyPair()
	sm.added = 0
	sm.taken = 0
	return sm
}

func (sm *timePoolSM) Check(t *rapid.T) {
	require.GreaterOrEqual(t, sm.added, sm.taken)
}

func (sm *timePoolSM) AddRequest(t *rapid.T) {
	ts := time.Unix(rapid.Int64().Draw(t, "req.ts"), 0)
	req := isc.NewOffLedgerRequest(isc.RandomChainID(), governance.Contract.Hname(), governance.FuncAddCandidateNode.Hname(), nil, 0, gas.LimitsDefault.MaxGasPerRequest).Sign(sm.kp)
	sm.tp.AddRequest(ts, req)
	sm.added++
}

func (sm *timePoolSM) TakeTill(t *rapid.T) {
	ts := time.Unix(rapid.Int64().Draw(t, "take.ts"), 0)
	res := sm.tp.TakeTill(ts)
	sm.taken += len(res)
}
