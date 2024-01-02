// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package chainmanager

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/gpa"
	"github.com/iotaledger/wasp/packages/isc"
)

type inputAliasOutputConfirmed struct {
	aliasOutput *isc.AliasOutputWithID
}

func NewInputAliasOutputConfirmed(aliasOutput *isc.AliasOutputWithID) gpa.Input {
	return &inputAliasOutputConfirmed{
		aliasOutput: aliasOutput,
	}
}

func (inp *inputAliasOutputConfirmed) String() string {
	return fmt.Sprintf("{chainMgr.inputAliasOutputConfirmed, %v}", inp.aliasOutput)
}
