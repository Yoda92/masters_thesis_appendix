// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Package peering provides an overlay network for communicating
// between nodes in a peer-to-peer style with low overhead
// encoding and persistent connections. The network provides only
// the asynchronous communication.
//
// It is intended to use for the committee consensus protocol.
package peering

import (
	"crypto/rand"
	"io"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/util/rwutil"
)

// PeeringID is relates peers in different nodes for a particular
// communication group. E.g. PeeringID identifies a committee in
// the consensus, etc.
type PeeringID [iotago.Ed25519AddressBytesLength]byte

func RandomPeeringID(seed ...[]byte) PeeringID {
	var pid PeeringID
	_, _ = rand.Read(pid[:])
	return pid
}

// HashPeeringIDFromBytes generates a PeeringID by concatenating all the given data and hash with Blake2b 256.
func HashPeeringIDFromBytes(src []byte, additional ...[]byte) PeeringID {
	hashed := hashing.HashDataBlake2b(append([][]byte{src}, additional...)...)
	pid := PeeringID{}
	copy(pid[:], hashed[:])
	return pid
}

func (pid *PeeringID) String() string {
	return iotago.EncodeHex(pid[:])
}

func (pid *PeeringID) Read(r io.Reader) error {
	rr := rwutil.NewReader(r)
	rr.ReadN(pid[:])
	return rr.Err
}

func (pid *PeeringID) Write(w io.Writer) error {
	ww := rwutil.NewWriter(w)
	ww.WriteN(pid[:])
	return ww.Err
}
