package coreutil

import (
	"github.com/iotaledger/wasp/packages/isc"
)

// names of core contracts
const (
	CoreContractRoot            = "root"
	CoreContractAccounts        = "accounts"
	CoreContractBlob            = "blob"
	CoreContractBlocklog        = "blocklog"
	CoreContractGovernance      = "governance"
	CoreContractErrors          = "errors"
	CoreContractEVM             = "evm"
	CoreEPRotateStateController = "rotateStateController"
)

var (
	CoreContractRootHname            = isc.Hn(CoreContractRoot)
	CoreContractAccountsHname        = isc.Hn(CoreContractAccounts)
	CoreContractBlobHname            = isc.Hn(CoreContractBlob)
	CoreContractBlocklogHname        = isc.Hn(CoreContractBlocklog)
	CoreContractGovernanceHname      = isc.Hn(CoreContractGovernance)
	CoreContractErrorsHname          = isc.Hn(CoreContractErrors)
	CoreContractEVMHname             = isc.Hn(CoreContractEVM)
	CoreEPRotateStateControllerHname = isc.Hn(CoreEPRotateStateController)

	hnames = map[string]isc.Hname{
		CoreContractRoot:       CoreContractRootHname,
		CoreContractAccounts:   CoreContractAccountsHname,
		CoreContractBlob:       CoreContractBlobHname,
		CoreContractBlocklog:   CoreContractBlocklogHname,
		CoreContractGovernance: CoreContractGovernanceHname,
		CoreContractEVM:        CoreContractEVMHname,
		CoreContractErrors:     CoreContractErrorsHname,
	}
)

// the global names used in 'blocklog' contract and in 'state' package
const (
	StateVarTimestamp           = "T"
	StateVarBlockIndex          = "I"
	StateVarPrevL1Commitment    = "H"
	ParamStateControllerAddress = "S"
)

// used in 'state' package as key for timestamp and block index
var (
	StatePrefixTimestamp        = string(CoreContractBlocklogHname.Bytes()) + StateVarTimestamp
	StatePrefixBlockIndex       = string(CoreContractBlocklogHname.Bytes()) + StateVarBlockIndex
	StatePrefixPrevL1Commitment = string(CoreContractBlocklogHname.Bytes()) + StateVarPrevL1Commitment
)

func CoreHname(name string) isc.Hname {
	if ret, ok := hnames[name]; ok {
		return ret
	}
	return isc.Hn(name)
}
