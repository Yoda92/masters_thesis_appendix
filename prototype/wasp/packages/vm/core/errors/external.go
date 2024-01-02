package errors

import (
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/core/errors/coreerrors"
)

// ViewCaller is a generic interface for any function that can call views
type ViewCaller func(contractName string, funcName string, params dict.Dict) (dict.Dict, error)

func GetMessageFormat(code isc.VMErrorCode, callView ViewCaller) (string, error) {
	ret, err := callView(Contract.Name, ViewGetErrorMessageFormat.Name, dict.Dict{
		ParamErrorCode: codec.EncodeVMErrorCode(code),
	})
	if err != nil {
		return "", err
	}
	return codec.DecodeString(ret.Get(ParamErrorMessageFormat))
}

func Resolve(e *isc.UnresolvedVMError, callView ViewCaller) (*isc.VMError, error) {
	if e == nil {
		return nil, nil
	}

	messageFormat, err := GetMessageFormat(e.Code(), callView)
	if err != nil {
		return nil, err
	}

	return isc.NewVMErrorTemplate(e.Code(), messageFormat).Create(e.Params...), nil
}

func ResolveFromState(state kv.KVStoreReader, e *isc.UnresolvedVMError) (*isc.VMError, error) {
	if e == nil {
		return nil, nil
	}
	template, ok := getErrorMessageFormat(state, e.Code())
	if !ok {
		return nil, coreerrors.ErrErrorNotFound
	}
	return isc.NewVMErrorTemplate(e.Code(), template.MessageFormat()).Create(e.Params...), nil
}
