package blob

import (
	"encoding/binary"
	"fmt"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/vm/vmtypes"
)

const DirectoryPrefix = "d"

func valuesMapName(blobHash hashing.HashValue) string {
	return "v" + string(blobHash[:])
}

func sizesMapName(blobHash hashing.HashValue) string {
	return "s" + string(blobHash[:])
}

func mustGetBlobHash(fields dict.Dict) (hashing.HashValue, []kv.Key, [][]byte) {
	sorted := fields.KeysSorted() // mind determinism
	values := make([][]byte, 0, len(sorted))
	all := make([][]byte, 0, 2*len(sorted))

	// hashBlob = hash(KeyLen0|Key0|Val0 | KeyLen1|Key1|Val1 | ... | KeyLenN|KeyN|ValN)
	// by prepend the key length we can avoid the possible collision
	for _, key := range sorted {
		var prefix [4]byte
		v := fields.Get(key)
		values = append(values, v)
		binary.LittleEndian.PutUint32(prefix[:], uint32(len(key)))
		all = append(all, prefix[:], []byte(key), v)
	}
	return hashing.HashData(all...), sorted, values
}

// MustGetBlobHash deterministically hashes map of binary values
func MustGetBlobHash(fields dict.Dict) hashing.HashValue {
	ret, _, _ := mustGetBlobHash(fields)
	return ret
}

// GetDirectory retrieves the blob directory from the state
func GetDirectory(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, DirectoryPrefix)
}

// GetDirectoryR retrieves the blob directory from the read-only state
func GetDirectoryR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, DirectoryPrefix)
}

// GetBlobValues retrieves the blob field-value map from the state
func GetBlobValues(state kv.KVStore, blobHash hashing.HashValue) *collections.Map {
	return collections.NewMap(state, valuesMapName(blobHash))
}

// GetBlobValuesR retrieves the blob field-value map from the read-only state
func GetBlobValuesR(state kv.KVStoreReader, blobHash hashing.HashValue) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, valuesMapName(blobHash))
}

// GetBlobSizes retrieves the writeable blob field-size map from the state
func GetBlobSizes(state kv.KVStore, blobHash hashing.HashValue) *collections.Map {
	return collections.NewMap(state, sizesMapName(blobHash))
}

// GetBlobSizesR retrieves the blob field-size map from the read-only state
func GetBlobSizesR(state kv.KVStoreReader, blobHash hashing.HashValue) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, sizesMapName(blobHash))
}

func LocateProgram(state kv.KVStoreReader, programHash hashing.HashValue) (string, []byte, error) {
	blbValues := GetBlobValuesR(state, programHash)
	programBinary := blbValues.GetAt([]byte(VarFieldProgramBinary))
	if programBinary == nil {
		return "", nil, fmt.Errorf("can't find program binary for hash %s", programHash.String())
	}
	v := blbValues.GetAt([]byte(VarFieldVMType))
	vmType := vmtypes.WasmTime
	if v != nil {
		vmType = string(v)
	}
	return vmType, programBinary, nil
}

func EncodeSize(size uint32) []byte {
	return codec.EncodeUint32(size)
}

func DecodeSize(size []byte) (uint32, error) {
	return codec.DecodeUint32(size)
}

func DecodeSizesMap(sizes dict.Dict) (map[string]uint32, error) {
	ret := make(map[string]uint32)
	for field, size := range sizes {
		v, err := DecodeSize(size)
		if err != nil {
			return nil, err
		}
		ret[string(field)] = v
	}
	return ret, nil
}

func DecodeDirectory(blobs dict.Dict) (map[hashing.HashValue]uint32, error) {
	ret := make(map[hashing.HashValue]uint32)
	for hash, size := range blobs {
		v, err := DecodeSize(size)
		if err != nil {
			return nil, err
		}
		h, err := codec.DecodeHashValue([]byte(hash))
		if err != nil {
			return nil, err
		}
		ret[h] = v
	}
	return ret, nil
}
