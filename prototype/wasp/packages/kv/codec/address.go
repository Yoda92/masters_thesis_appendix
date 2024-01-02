package codec

import (
	"errors"

	"github.com/iotaledger/hive.go/serializer/v2"
	iotago "github.com/iotaledger/iota.go/v3"
)

func DecodeAddress(b []byte, def ...iotago.Address) (iotago.Address, error) {
	if b == nil {
		if len(def) == 0 {
			return nil, errors.New("cannot decode nil Address")
		}
		return def[0], nil
	}
	if len(b) == 0 {
		return nil, errors.New("invalid Address size")
	}
	typeByte := b[0]
	addr, err := iotago.AddressSelector(uint32(typeByte))
	if err != nil {
		return nil, err
	}
	_, err = addr.Deserialize(b, serializer.DeSeriModePerformValidation, nil)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func MustDecodeAddress(b []byte, def ...iotago.Address) iotago.Address {
	a, err := DecodeAddress(b, def...)
	if err != nil {
		panic(err)
	}
	return a
}

func EncodeAddress(addr iotago.Address) []byte {
	addressInBytes, err := addr.Serialize(serializer.DeSeriModeNoValidation, nil)
	if err != nil {
		panic("cannot encode address")
	}
	return addressInBytes
}
