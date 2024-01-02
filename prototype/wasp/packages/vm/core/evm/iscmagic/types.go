// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package iscmagic

import (
	"math/big"
	"time"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

// ISCChainID matches the type definition in ISCTypes.sol
type ISCChainID [isc.ChainIDLength]byte

func init() {
	if isc.ChainIDLength != 32 {
		panic("static check: ChainID length does not match bytes32 in ISCTypes.sol")
	}
}

func WrapISCChainID(c isc.ChainID) (ret ISCChainID) {
	copy(ret[:], c.Bytes())
	return
}

func (c ISCChainID) Unwrap() (isc.ChainID, error) {
	return isc.ChainIDFromBytes(c[:])
}

func (c ISCChainID) MustUnwrap() isc.ChainID {
	ret, err := c.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// NativeTokenID matches the struct definition in ISCTypes.sol
type NativeTokenID struct {
	Data []byte
}

func WrapNativeTokenID(nativeTokenID iotago.NativeTokenID) NativeTokenID {
	return NativeTokenID{Data: nativeTokenID[:]}
}

func (a NativeTokenID) Unwrap() (ret iotago.NativeTokenID) {
	copy(ret[:], a.Data)
	return
}

func (a NativeTokenID) MustUnwrap() (ret iotago.NativeTokenID) {
	copy(ret[:], a.Data)
	return
}

// NativeToken matches the struct definition in ISCTypes.sol
type NativeToken struct {
	ID     NativeTokenID
	Amount *big.Int
}

func WrapNativeToken(nativeToken *iotago.NativeToken) NativeToken {
	return NativeToken{
		ID:     WrapNativeTokenID(nativeToken.ID),
		Amount: nativeToken.Amount,
	}
}

func (nt NativeToken) Unwrap() *iotago.NativeToken {
	return &iotago.NativeToken{
		ID:     nt.ID.Unwrap(),
		Amount: nt.Amount,
	}
}

// L1Address matches the struct definition in ISCTypes.sol
type L1Address struct {
	Data []byte
}

func WrapL1Address(a iotago.Address) L1Address {
	if a == nil {
		return L1Address{Data: []byte{}}
	}
	return L1Address{Data: isc.AddressToBytes(a)}
}

func (a L1Address) Unwrap() (iotago.Address, error) {
	ret, err := isc.AddressFromBytes(a.Data)
	return ret, err
}

func (a L1Address) MustUnwrap() iotago.Address {
	ret, err := a.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// ISCAgentID matches the struct definition in ISCTypes.sol
type ISCAgentID struct {
	Data []byte
}

func WrapISCAgentID(a isc.AgentID) ISCAgentID {
	return ISCAgentID{Data: a.Bytes()}
}

func (a ISCAgentID) Unwrap() (isc.AgentID, error) {
	return isc.AgentIDFromBytes(a.Data)
}

func (a ISCAgentID) MustUnwrap() isc.AgentID {
	ret, err := a.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// ISCRequestID matches the struct definition in ISCTypes.sol
type ISCRequestID struct {
	Data []byte
}

func WrapISCRequestID(rid isc.RequestID) ISCRequestID {
	return ISCRequestID{Data: rid.Bytes()}
}

func (rid ISCRequestID) Unwrap() (isc.RequestID, error) {
	return isc.RequestIDFromBytes(rid.Data)
}

func (rid ISCRequestID) MustUnwrap() isc.RequestID {
	ret, err := rid.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// NFTID matches the type definition in ISCTypes.sol
type NFTID [iotago.NFTIDLength]byte

func init() {
	if iotago.NFTIDLength != 32 {
		panic("static check: NFTID length does not match bytes32 in ISCTypes.sol")
	}
}

func WrapNFTID(c iotago.NFTID) (ret NFTID) {
	copy(ret[:], c[:])
	return
}

func (n NFTID) Unwrap() (ret iotago.NFTID) {
	copy(ret[:], n[:])
	return
}

// TokenID returns the uint256 tokenID for ERC721
func (n NFTID) TokenID() *big.Int {
	return new(big.Int).SetBytes(n[:])
}

// ISCNFT matches the struct definition in ISCTypes.sol
type ISCNFT struct {
	ID       NFTID
	Issuer   L1Address
	Metadata []byte
	Owner    ISCAgentID
}

func WrapISCNFT(n *isc.NFT) ISCNFT {
	r := ISCNFT{
		ID:       WrapNFTID(n.ID),
		Issuer:   WrapL1Address(n.Issuer),
		Metadata: n.Metadata,
	}
	if n.Owner != nil {
		r.Owner = WrapISCAgentID(n.Owner)
	}
	return r
}

func (n ISCNFT) Unwrap() (*isc.NFT, error) {
	issuer, err := n.Issuer.Unwrap()
	if err != nil {
		return nil, err
	}
	return &isc.NFT{
		ID:       n.ID.Unwrap(),
		Issuer:   issuer,
		Metadata: n.Metadata,
		Owner:    n.Owner.MustUnwrap(),
	}, nil
}

func (n ISCNFT) MustUnwrap() *isc.NFT {
	ret, err := n.Unwrap()
	if err != nil {
		panic(err)
	}
	return ret
}

// IRC27NFTMetadata matches the struct definition in ISCTypes.sol
type IRC27NFTMetadata struct {
	Standard string
	Version  string
	MimeType string
	Uri      string //nolint:revive // false positive
	Name     string
}

func WrapIRC27NFTMetadata(m *isc.IRC27NFTMetadata) IRC27NFTMetadata {
	return IRC27NFTMetadata{
		Standard: m.Standard,
		Version:  m.Version,
		MimeType: m.MIMEType,
		Uri:      m.URI,
		Name:     m.Name,
	}
}

// IRC27NFT matches the struct definition in ISCTypes.sol
type IRC27NFT struct {
	Nft      ISCNFT
	Metadata IRC27NFTMetadata
}

// ISCAssets matches the struct definition in ISCTypes.sol
type ISCAssets struct {
	BaseTokens   uint64
	NativeTokens []NativeToken
	Nfts         []NFTID
}

func WrapISCAssets(a *isc.Assets) ISCAssets {
	if a == nil {
		return WrapISCAssets(isc.NewEmptyAssets())
	}
	tokens := make([]NativeToken, len(a.NativeTokens))
	for i, nativeToken := range a.NativeTokens {
		tokens[i] = WrapNativeToken(nativeToken)
	}
	nfts := make([]NFTID, len(a.NFTs))
	for i, id := range a.NFTs {
		nfts[i] = WrapNFTID(id)
	}
	return ISCAssets{
		BaseTokens:   a.BaseTokens,
		NativeTokens: tokens,
		Nfts:         nfts,
	}
}

func (a ISCAssets) Unwrap() *isc.Assets {
	tokens := make(iotago.NativeTokens, len(a.NativeTokens))
	for i, nativeToken := range a.NativeTokens {
		tokens[i] = nativeToken.Unwrap()
	}
	nfts := make([]iotago.NFTID, len(a.Nfts))
	for i, id := range a.Nfts {
		nfts[i] = id.Unwrap()
	}
	return isc.NewAssets(a.BaseTokens, tokens, nfts...)
}

// ISCDictItem matches the struct definition in ISCTypes.sol
type ISCDictItem struct {
	Key   []byte
	Value []byte
}

// ISCDict matches the struct definition in ISCTypes.sol
type ISCDict struct {
	Items []ISCDictItem
}

func WrapISCDict(d dict.Dict) ISCDict {
	items := make([]ISCDictItem, 0, len(d))
	for k, v := range d {
		items = append(items, ISCDictItem{Key: []byte(k), Value: v})
	}
	return ISCDict{Items: items}
}

func (d ISCDict) Unwrap() dict.Dict {
	ret := dict.Dict{}
	for _, item := range d.Items {
		ret[kv.Key(item.Key)] = item.Value
	}
	return ret
}

type ISCSendMetadata struct {
	TargetContract uint32
	Entrypoint     uint32
	Params         ISCDict
	Allowance      ISCAssets
	GasBudget      uint64
}

func WrapISCSendMetadata(metadata isc.SendMetadata) ISCSendMetadata {
	ret := ISCSendMetadata{
		GasBudget:      metadata.GasBudget,
		Entrypoint:     uint32(metadata.EntryPoint),
		TargetContract: uint32(metadata.TargetContract),
		Allowance:      WrapISCAssets(metadata.Allowance),
		Params:         WrapISCDict(metadata.Params),
	}

	return ret
}

func (i ISCSendMetadata) Unwrap() *isc.SendMetadata {
	ret := isc.SendMetadata{
		TargetContract: isc.Hname(i.TargetContract),
		EntryPoint:     isc.Hname(i.Entrypoint),
		Params:         i.Params.Unwrap(),
		Allowance:      i.Allowance.Unwrap(),
		GasBudget:      i.GasBudget,
	}

	return &ret
}

type ISCExpiration struct {
	Time          int64
	ReturnAddress L1Address
}

func (i *ISCExpiration) Unwrap() *isc.Expiration {
	if i == nil {
		return nil
	}

	if i.Time == 0 {
		return nil
	}

	address := i.ReturnAddress.MustUnwrap()

	ret := isc.Expiration{
		ReturnAddress: address,
		Time:          time.UnixMilli(i.Time),
	}

	return &ret
}

type ISCSendOptions struct {
	Timelock   int64
	Expiration ISCExpiration
}

func (i *ISCSendOptions) Unwrap() isc.SendOptions {
	var timeLock time.Time

	if i.Timelock > 0 {
		timeLock = time.UnixMilli(i.Timelock)
	}

	ret := isc.SendOptions{
		Timelock:   timeLock,
		Expiration: i.Expiration.Unwrap(),
	}

	return ret
}

type ISCTokenProperties struct {
	Name         string
	TickerSymbol string
	Decimals     uint8
	TotalSupply  *big.Int
}
