// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package cons

import (
	"fmt"
	"strings"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/state"
	"github.com/iotaledger/wasp/packages/vm"
)

type SyncTX interface {
	VMResultReceived(vmResult *vm.VMTaskResult) gpa.OutMessages
	SignatureReceived(signature []byte) gpa.OutMessages
	BlockSaved(block state.Block) gpa.OutMessages
	String() string
}

type syncTXImpl struct {
	vmResult   *vm.VMTaskResult
	signature  []byte
	blockSaved bool
	block      state.Block

	inputsReady   bool
	inputsReadyCB func(vmResult *vm.VMTaskResult, block state.Block, signature []byte) gpa.OutMessages
}

func NewSyncTX(inputsReadyCB func(vmResult *vm.VMTaskResult, block state.Block, signature []byte) gpa.OutMessages) SyncTX {
	return &syncTXImpl{inputsReadyCB: inputsReadyCB}
}

func (sub *syncTXImpl) VMResultReceived(vmResult *vm.VMTaskResult) gpa.OutMessages {
	if sub.vmResult != nil || vmResult == nil {
		return nil
	}
	sub.vmResult = vmResult
	return sub.tryCompleteInputs()
}

func (sub *syncTXImpl) SignatureReceived(signature []byte) gpa.OutMessages {
	if sub.signature != nil || signature == nil {
		return nil
	}
	sub.signature = signature
	return sub.tryCompleteInputs()
}

func (sub *syncTXImpl) BlockSaved(block state.Block) gpa.OutMessages {
	if sub.blockSaved {
		return nil
	}
	sub.blockSaved = true
	sub.block = block
	return sub.tryCompleteInputs()
}

func (sub *syncTXImpl) tryCompleteInputs() gpa.OutMessages {
	if sub.inputsReady || sub.vmResult == nil || sub.signature == nil || !sub.blockSaved {
		return nil
	}
	sub.inputsReady = true
	return sub.inputsReadyCB(sub.vmResult, sub.block, sub.signature)
}

// Try to provide useful human-readable compact status.
func (sub *syncTXImpl) String() string {
	str := "TX"
	if sub.inputsReady {
		str += statusStrOK
	} else {
		wait := []string{}
		if sub.vmResult == nil {
			wait = append(wait, "VMResult")
		}
		if sub.signature == nil {
			wait = append(wait, "Signature")
		}
		if !sub.blockSaved {
			wait = append(wait, "SavedBlock")
		}
		str += fmt.Sprintf("/WAIT[%v]", strings.Join(wait, ","))
	}
	return str
}
