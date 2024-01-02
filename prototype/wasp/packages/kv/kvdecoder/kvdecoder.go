package kvdecoder

import (
	"fmt"
	"math/big"
	"time"

	"github.com/iotaledger/hive.go/serializer/v2"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
)

type kvdecoder struct {
	kv.KVStoreReader
	log isc.LogInterface
}

var _ isc.KVDecoder = &kvdecoder{}

func New(kvReader kv.KVStoreReader, log ...isc.LogInterface) isc.KVDecoder {
	var l isc.LogInterface
	if len(log) > 0 {
		l = log[0]
	}
	return &kvdecoder{kvReader, l}
}

func (p *kvdecoder) check(err error) {
	if err == nil {
		return
	}
	if p.log == nil {
		panic(err)
	}
	p.log.Panicf("%v", err)
}

func (p *kvdecoder) wrapError(key kv.Key, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("cannot decode key '%s': %w", key, err)
}

func (p *kvdecoder) GetInt16(key kv.Key, def ...int16) (int16, error) {
	v, err := codec.DecodeInt16(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetInt16(key kv.Key, def ...int16) int16 {
	ret, err := p.GetInt16(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetUint16(key kv.Key, def ...uint16) (uint16, error) {
	v, err := codec.DecodeUint16(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetUint16(key kv.Key, def ...uint16) uint16 {
	ret, err := p.GetUint16(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetInt32(key kv.Key, def ...int32) (int32, error) {
	v, err := codec.DecodeInt32(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetInt32(key kv.Key, def ...int32) int32 {
	ret, err := p.GetInt32(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetUint32(key kv.Key, def ...uint32) (uint32, error) {
	v, err := codec.DecodeUint32(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetUint32(key kv.Key, def ...uint32) uint32 {
	ret, err := p.GetUint32(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetInt64(key kv.Key, def ...int64) (int64, error) {
	v, err := codec.DecodeInt64(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetInt64(key kv.Key, def ...int64) int64 {
	ret, err := p.GetInt64(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetUint64(key kv.Key, def ...uint64) (uint64, error) {
	v, err := codec.DecodeUint64(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetUint64(key kv.Key, def ...uint64) uint64 {
	ret, err := p.GetUint64(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetBool(key kv.Key, def ...bool) (bool, error) {
	v, err := codec.DecodeBool(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetBool(key kv.Key, def ...bool) bool {
	ret, err := p.GetBool(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetTime(key kv.Key, def ...time.Time) (time.Time, error) {
	v, err := codec.DecodeTime(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetTime(key kv.Key, def ...time.Time) time.Time {
	ret, err := p.GetTime(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetString(key kv.Key, def ...string) (string, error) {
	v, err := codec.DecodeString(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetString(key kv.Key, def ...string) string {
	ret, err := p.GetString(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetHname(key kv.Key, def ...isc.Hname) (isc.Hname, error) {
	v, err := codec.DecodeHname(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetHname(key kv.Key, def ...isc.Hname) isc.Hname {
	ret, err := p.GetHname(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetHashValue(key kv.Key, def ...hashing.HashValue) (hashing.HashValue, error) {
	v, err := codec.DecodeHashValue(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetHashValue(key kv.Key, def ...hashing.HashValue) hashing.HashValue {
	ret, err := p.GetHashValue(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetAddress(key kv.Key, def ...iotago.Address) (iotago.Address, error) {
	v, err := codec.DecodeAddress(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetAddress(key kv.Key, def ...iotago.Address) iotago.Address {
	ret, err := p.GetAddress(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetRequestID(key kv.Key, def ...isc.RequestID) (isc.RequestID, error) {
	v, err := codec.DecodeRequestID(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetRequestID(key kv.Key, def ...isc.RequestID) isc.RequestID {
	ret, err := p.GetRequestID(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetAgentID(key kv.Key, def ...isc.AgentID) (isc.AgentID, error) {
	v, err := codec.DecodeAgentID(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetAgentID(key kv.Key, def ...isc.AgentID) isc.AgentID {
	ret, err := p.GetAgentID(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetChainID(key kv.Key, def ...isc.ChainID) (isc.ChainID, error) {
	v, err := codec.DecodeChainID(p.Get(key), def...)
	return v, p.wrapError(key, err)
}

func (p *kvdecoder) MustGetChainID(key kv.Key, def ...isc.ChainID) isc.ChainID {
	ret, err := p.GetChainID(key, def...)
	p.check(err)
	return ret
}

// nil means does not exist
func (p *kvdecoder) GetBytes(key kv.Key, def ...[]byte) ([]byte, error) {
	v := p.Get(key)
	if v != nil {
		return v, nil
	}
	if len(def) == 0 {
		return nil, fmt.Errorf("GetBytes: mandatory parameter '%s' does not exist", key)
	}
	return def[0], nil
}

func (p *kvdecoder) MustGetBytes(key kv.Key, def ...[]byte) []byte {
	ret, err := p.GetBytes(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetTokenScheme(key kv.Key, def ...iotago.TokenScheme) (iotago.TokenScheme, error) {
	v := p.Get(key)
	if len(v) > 1 {
		ts, err := iotago.TokenSchemeSelector(uint32(v[0]))
		if err != nil {
			return nil, err
		}
		_, err = ts.Deserialize(v, serializer.DeSeriModeNoValidation, nil)
		if err != nil {
			return nil, err
		}
		return ts, nil
	}
	if len(def) == 0 {
		return nil, fmt.Errorf("GetTokenScheme: mandatory parameter '%s' does not exist", key)
	}
	return def[0], nil
}

func (p *kvdecoder) MustGetTokenScheme(key kv.Key, def ...iotago.TokenScheme) iotago.TokenScheme {
	ret, err := p.GetTokenScheme(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetBigInt(key kv.Key, def ...*big.Int) (*big.Int, error) {
	v := p.Get(key)
	if v == nil {
		if len(def) != 0 {
			return def[0], nil
		}
		return nil, fmt.Errorf("GetBigInt: mandatory parameter '%s' does not exist", key)
	}
	return codec.DecodeBigIntAbs(v)
}

func (p *kvdecoder) MustGetBigInt(key kv.Key, def ...*big.Int) *big.Int {
	ret, err := p.GetBigInt(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetNativeTokenID(key kv.Key, def ...iotago.NativeTokenID) (iotago.NativeTokenID, error) {
	v := p.Get(key)
	if v == nil {
		if len(def) != 0 {
			return def[0], nil
		}
		return iotago.NativeTokenID{}, fmt.Errorf("GetNativeTokenID: mandatory parameter %q does not exist", key)
	}
	return codec.DecodeNativeTokenID(v)
}

func (p *kvdecoder) MustGetNativeTokenID(key kv.Key, def ...iotago.NativeTokenID) iotago.NativeTokenID {
	ret, err := p.GetNativeTokenID(key, def...)
	p.check(err)
	return ret
}

func (p *kvdecoder) GetNFTID(key kv.Key, def ...iotago.NFTID) (iotago.NFTID, error) {
	v := p.Get(key)
	if v == nil {
		if len(def) != 0 {
			return def[0], nil
		}
		return iotago.NFTID{}, fmt.Errorf("GetNFTID: mandatory parameter %q does not exist", key)
	}
	return codec.DecodeNFTID(v)
}

func (p *kvdecoder) MustGetNFTID(key kv.Key, def ...iotago.NFTID) iotago.NFTID {
	ret, err := p.GetNFTID(key, def...)
	p.check(err)
	return ret
}
