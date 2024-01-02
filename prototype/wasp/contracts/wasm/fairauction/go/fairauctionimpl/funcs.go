// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package fairauctionimpl

import (
	"github.com/iotaledger/wasp/contracts/wasm/fairauction/go/fairauction"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

const (
	// default duration is 60 min
	DurationDefault = 60
	// minimum duration is 1 min
	DurationMin = 1
	// maximum duration is 120 min
	DurationMax          = 120
	MaxDescriptionLength = 150
	OwnerMarginDefault   = 50
	OwnerMarginMin       = 5
	OwnerMarginMax       = 100
)

func funcStartAuction(ctx wasmlib.ScFuncContext, f *StartAuctionContext) {
	allowance := ctx.Allowance()
	nfts := allowance.NftIDs()
	ctx.Require(len(nfts) == 1, "single NFT allowance expected")
	auctionNFT := *nfts[0]

	minimumBid := f.Params.MinimumBid().Value()

	// duration in minutes
	duration := f.Params.Duration().Value()
	if duration == 0 {
		duration = DurationDefault
	}
	if duration < DurationMin {
		duration = DurationMin
	}
	if duration > DurationMax {
		duration = DurationMax
	}

	description := f.Params.Description().Value()
	if description == "" {
		description = "N/A"
	}
	if len(description) > MaxDescriptionLength {
		ss := description[:MaxDescriptionLength]
		description = ss + "[...]"
	}

	ownerMargin := f.State.OwnerMargin().Value()
	if ownerMargin == 0 {
		ownerMargin = OwnerMarginDefault
	}

	// need at least 1 base token (storage deposit) to run SC
	margin := minimumBid * ownerMargin / 1000
	if margin == 0 {
		margin = 1
	}
	deposit := allowance.BaseTokens()
	if deposit < margin {
		ctx.Panic("Insufficient deposit")
	}

	currentAuction := f.State.Auctions().GetAuction(auctionNFT)
	if currentAuction.Exists() {
		ctx.Panic("Auction for this nft already exists")
	}

	auction := &fairauction.Auction{
		Creator:       ctx.Caller(),
		Deposit:       deposit,
		Description:   description,
		Duration:      duration,
		HighestBid:    0,
		HighestBidder: ctx.Caller(),
		MinimumBid:    minimumBid,
		OwnerMargin:   ownerMargin,
		Nft:           auctionNFT,
		WhenStarted:   ctx.Timestamp(),
	}
	currentAuction.SetValue(auction)

	// take custody of deposit and NFT
	transfer := wasmlib.ScTransferFromBaseTokens(deposit)
	transfer.AddNFT(&auctionNFT)
	ctx.TransferAllowed(ctx.AccountID(), transfer)

	fa := fairauction.ScFuncs.FinalizeAuction(ctx)
	fa.Params.Nft().SetValue(auction.Nft)
	fa.Func.Delay(duration * 60).Post()
}

func funcPlaceBid(ctx wasmlib.ScFuncContext, f *PlaceBidContext) {
	bidAmount := ctx.Allowance().BaseTokens()
	ctx.Require(bidAmount > 0, "Missing bid amount")

	token := f.Params.Nft().Value()
	currentAuction := f.State.Auctions().GetAuction(token)
	ctx.Require(currentAuction.Exists(), "Missing auction info")

	auction := currentAuction.Value()
	bids := f.State.Bids().GetBids(token)
	bidderList := f.State.BidderList().GetBidderList(token)
	caller := ctx.Caller()
	currentBid := bids.GetBid(caller)
	if currentBid.Exists() {
		ctx.Log("Upped bid from: " + caller.String())
		bid := currentBid.Value()
		bidAmount += bid.Amount
		bid.Amount = bidAmount
		bid.Timestamp = ctx.Timestamp()
		currentBid.SetValue(bid)
	} else {
		ctx.Require(bidAmount >= auction.MinimumBid, "Insufficient bid amount")
		ctx.Log("New bid from: " + caller.String())
		index := bidderList.Length()
		bidderList.AppendAgentID().SetValue(caller)
		bid := &fairauction.Bid{
			Index:     index,
			Amount:    bidAmount,
			Timestamp: ctx.Timestamp(),
		}
		currentBid.SetValue(bid)
	}
	if bidAmount > auction.HighestBid {
		ctx.Log("New highest bidder")
		auction.HighestBid = bidAmount
		auction.HighestBidder = caller
		currentAuction.SetValue(auction)
	}
}

func funcFinalizeAuction(ctx wasmlib.ScFuncContext, f *FinalizeAuctionContext) {
	auctionNFT := f.Params.Nft().Value()
	currentAuction := f.State.Auctions().GetAuction(auctionNFT)
	ctx.Require(currentAuction.Exists(), "Missing auction info")
	auction := currentAuction.Value()
	if auction.HighestBid == 0 {
		ctx.Log("No one bid on " + auctionNFT.String())
		ownerFee := auction.MinimumBid * auction.OwnerMargin / 1000
		if ownerFee == 0 {
			ownerFee = 1
		}
		// finalizeAuction request token was probably not confirmed yet
		transferTokens(ctx, f.State.Owner().Value(), ownerFee-1)
		transferNFT(ctx, auction.Creator, auction.Nft)
		transferTokens(ctx, auction.Creator, auction.Deposit-ownerFee)
		return
	}

	ownerFee := auction.HighestBid * auction.OwnerMargin / 1000
	if ownerFee == 0 {
		ownerFee = 1
	}

	// return staked bids to losers
	bids := f.State.Bids().GetBids(auctionNFT)
	bidderList := f.State.BidderList().GetBidderList(auctionNFT)
	size := bidderList.Length()
	for i := uint32(0); i < size; i++ {
		loser := bidderList.GetAgentID(i).Value()
		if loser != auction.HighestBidder {
			bid := bids.GetBid(loser).Value()
			transferTokens(ctx, loser, bid.Amount)
		}
	}

	// finalizeAuction request token was probably not confirmed yet
	transferTokens(ctx, f.State.Owner().Value(), ownerFee-1)
	transferNFT(ctx, auction.HighestBidder, auction.Nft)
	transferTokens(ctx, auction.Creator, auction.Deposit+auction.HighestBid-ownerFee)
}

func funcSetOwnerMargin(_ wasmlib.ScFuncContext, f *SetOwnerMarginContext) {
	ownerMargin := f.Params.OwnerMargin().Value()
	if ownerMargin < OwnerMarginMin {
		ownerMargin = OwnerMarginMin
	}
	if ownerMargin > OwnerMarginMax {
		ownerMargin = OwnerMarginMax
	}
	f.State.OwnerMargin().SetValue(ownerMargin)
}

func viewGetAuctionInfo(ctx wasmlib.ScViewContext, f *GetAuctionInfoContext) {
	token := f.Params.Nft().Value()
	currentAuction := f.State.Auctions().GetAuction(token)
	if !currentAuction.Exists() {
		ctx.Panic("Missing auction info")
	}

	auction := currentAuction.Value()
	f.Results.Creator().SetValue(auction.Creator)
	f.Results.Deposit().SetValue(auction.Deposit)
	f.Results.Description().SetValue(auction.Description)
	f.Results.Duration().SetValue(auction.Duration)
	f.Results.HighestBid().SetValue(auction.HighestBid)
	f.Results.HighestBidder().SetValue(auction.HighestBidder)
	f.Results.MinimumBid().SetValue(auction.MinimumBid)
	f.Results.OwnerMargin().SetValue(auction.OwnerMargin)
	f.Results.Nft().SetValue(auction.Nft)
	f.Results.WhenStarted().SetValue(auction.WhenStarted)

	bidderList := f.State.BidderList().GetBidderList(token)
	f.Results.Bidders().SetValue(bidderList.Length())
}

func transferTokens(ctx wasmlib.ScFuncContext, agent wasmtypes.ScAgentID, amount uint64) {
	if agent.IsAddress() {
		// send back to original Tangle address
		ctx.Send(agent.Address(), wasmlib.ScTransferFromBaseTokens(amount))
		return
	}

	// TODO not an address, deposit into account on chain
	ctx.Send(agent.Address(), wasmlib.ScTransferFromBaseTokens(amount))
}

func transferNFT(ctx wasmlib.ScFuncContext, agent wasmtypes.ScAgentID, nft wasmtypes.ScNftID) {
	if agent.IsAddress() {
		// send back to original Tangle address
		ctx.Send(agent.Address(), wasmlib.ScTransferFromNFT(&nft))
		return
	}

	// TODO not an address, deposit into account on chain
	ctx.Send(agent.Address(), wasmlib.ScTransferFromNFT(&nft))
}

func funcInit(ctx wasmlib.ScFuncContext, f *InitContext) {
	if f.Params.Owner().Exists() {
		f.State.Owner().SetValue(f.Params.Owner().Value())
		return
	}
	f.State.Owner().SetValue(ctx.RequestSender())
}
