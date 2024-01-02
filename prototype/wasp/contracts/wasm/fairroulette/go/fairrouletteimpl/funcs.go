// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// This example implements 'fairroulette', a simple smart contract that can automatically handle
// an unlimited amount of bets on a number during a timed betting round. Once a betting round
// is over the contract will automatically pay out the winners proportionally to their bet amount.
// The intent is to showcase basic functionality of WasmLib and timed calling of functions
// through a minimal implementation and not to come up with a complete real-world solution.

package fairrouletteimpl

import (
	"github.com/iotaledger/wasp/contracts/wasm/fairroulette/go/fairroulette"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
)

// Define some default configuration parameters.

// The maximum number one can bet on. The range of numbers starts at 1.
const MaxNumber = 8

// The default playing period of one betting round in seconds.
const DefaultPlayPeriod = 60

// Enable this if you deploy the contract to an actual node. It will pay out the prize after a certain timeout.
const EnableSelfPost = true

// The number to divide nano seconds to seconds.
const NanoTimeDivider = 1000_000_000

// 'placeBet' is used by betters to place a bet on a number from 1 to MAX_NUMBER. The first
// incoming bet triggers a betting round of configurable duration. After the playing period
// expires the smart contract will automatically pay out any winners and start a new betting
// round upon arrival of a new bet.
// The 'placeBet' function takes 1 mandatory parameter:
// - 'number', which must be an Int64 number from 1 to MAX_NUMBER
// The 'member' function will save the number together with the address of the better and
// the amount of incoming tokens as the bet amount in its state.
func funcPlaceBet(ctx wasmlib.ScFuncContext, f *PlaceBetContext) {
	// Get the array of current bets from state storage.
	bets := f.State.Bets()

	nrOfBets := bets.Length()
	for i := uint32(0); i < nrOfBets; i++ {
		bet := bets.GetBet(i).Value()

		if bet.Better.Address() == ctx.Caller().Address() {
			ctx.Panic("Bet already placed for this round")
		}
	}

	// Since we are sure that the 'number' parameter actually exists we can
	// retrieve its actual value into an i64.
	number := f.Params.Number().Value()

	// Require that the number is a valid number to bet on, otherwise panic out.
	ctx.Require(number >= 1 && number <= MaxNumber, "invalid number")

	// Create ScBalances proxy to the allowance balances for this request.
	// Note that ScBalances wraps an ScImmutableMap of token color/amount combinations
	// in a simpler to use interface.
	allowance := ctx.Allowance()

	// Retrieve the amount of plain iota tokens that are part of the allowance balance.
	amount := allowance.BaseTokens()

	// Require that there are actually some base tokens there
	ctx.Require(amount > 0, "empty bet")

	// Now we gather all information together into a single serializable struct
	// Note that we use the caller() method of the function context to determine
	// the agent id of the better. This is where a potential pay-out will be sent.
	bet := &fairroulette.Bet{
		Better: ctx.Caller(),
		Amount: amount,
		Number: number,
	}

	// Append the bet data to the bets array. The bet array will automatically take care
	// of serializing the bet struct into a bytes representation.
	bets.AppendBet().SetValue(bet)

	f.Events.Bet(bet.Better.Address(), bet.Amount, bet.Number)

	// Was this the first bet of this round?
	if nrOfBets == 0 {
		// Yes it was, query the state for the length of the playing period in seconds by
		// retrieving the playPeriod value from state storage
		playPeriod := f.State.PlayPeriod().Value()

		// if the play period is less than 10 seconds we override it with the default duration.
		// Note that this will also happen when the play period was not set yet because in that
		// case a zero value was returned.
		if playPeriod < 10 {
			playPeriod = DefaultPlayPeriod
		}

		if EnableSelfPost {
			f.State.RoundStatus().SetValue(1)

			// timestamp is nanotime, divide by NanoTimeDivider to get seconds => common unix timestamp
			timestamp := uint32(ctx.Timestamp() / NanoTimeDivider)
			f.State.RoundStartedAt().SetValue(timestamp)

			f.Events.Start()

			roundNumber := f.State.RoundNumber()
			roundNumber.SetValue(roundNumber.Value() + 1)

			f.Events.Round(roundNumber.Value())

			// And now for our next trick we post a delayed request to ourselves on the Tangle.
			// We are requesting to call the 'payWinners' function, but delay it for the playPeriod
			// amount of seconds. This will lock in the playing period, during which more bets can
			// be placed. Once the 'payWinners' function gets triggered by the ISC it will gather
			// all bets up to that moment as the ones to consider for determining the winner.
			fairroulette.ScFuncs.PayWinners(ctx).Func.Delay(playPeriod).Post()
		}
	}
}

// 'payWinners' is a function whose execution gets initiated by the 'placeBet' function.
// It collects a list of all bets, generates a random number, sorts out the winners and transfers
// the calculated winning sum to each attendee.
//
//nolint:funlen
func funcPayWinners(ctx wasmlib.ScFuncContext, f *PayWinnersContext) {
	// Use the built-in random number generator which has been automatically initialized by
	// using the transaction hash as initial entropy data. Note that the pseudo-random number
	// generator will use the next 8 bytes from the hash as its random Int64 number and once
	// it runs out of data it simply hashes the previous hash for a next pseudo-random sequence.
	// Here we determine the winning number for this round in the range of 1 thru MaxNumber.
	winningNumber := uint16(ctx.Random(MaxNumber-1) + 1)

	// Save the last winning number in state storage under 'lastWinningNumber' so that there
	// is (limited) time for people to call the 'getLastWinningNumber' View to verify the last
	// winning number if they wish. Note that this is just a silly example. We could log much
	// more extensive statistics information about each playing round in state storage and
	// make that data available through views for anyone to see.
	f.State.LastWinningNumber().SetValue(winningNumber)

	// Gather all winners and calculate some totals at the same time.
	// Keep track of the total bet amount, the total win amount, and all the winners.
	// Note how we decided to keep the winners in a local vector instead of creating
	// yet another array in state storage or having to go through lockedBets again.
	totalBetAmount := uint64(0)
	totalWinAmount := uint64(0)
	winners := make([]*fairroulette.Bet, 0)

	// Get the 'bets' array in state storage.
	bets := f.State.Bets()

	// Determine the amount of bets in the 'bets' array.
	nrOfBets := bets.Length()

	// Loop through all indexes of the 'bets' array.
	for i := uint32(0); i < nrOfBets; i++ {
		// Retrieve the bet stored at the next index
		bet := bets.GetBet(i).Value()

		// Add this bet's amount to the running total bet amount
		totalBetAmount += bet.Amount

		// Did this better bet on the winning number?
		if bet.Number == winningNumber {
			// Yes, add this bet amount to the running total win amount.
			totalWinAmount += bet.Amount

			// And save this bet in the winners vector.
			winners = append(winners, bet)
		}
	}

	// Now that we preprocessed all bets we can get rid of the data in state storage
	// so that the 'bets' array becomes available for when the next betting round ends.
	bets.Clear()

	f.Events.Winner(winningNumber)
	// Did we have any winners at all?
	if len(winners) == 0 {
		// No winners, log this fact to the log on the host.
		ctx.Log("Nobody wins!")
	}

	// Pay out the winners proportionally to their bet amount. Note that we could configure
	// a small percentage that would go to the owner of the smart contract as hosting payment.

	// Keep track of the total payout so we can calculate the remainder after truncation.
	totalPayout := uint64(0)

	// Loop through all winners.
	size := len(winners)
	for i := 0; i < size; i++ {
		// Get the next winner.
		bet := winners[i]

		// Determine the proportional win amount (we could take our percentage here)
		payout := totalBetAmount * bet.Amount / totalWinAmount

		// Anything to pay to the winner?
		if payout != 0 {
			// Yep, keep track of the running total payout
			totalPayout += payout

			// Set up an ScTransfer proxy that transfers the correct amount of tokens.
			// Note that ScTransfer wraps an ScMutableMap of token color/amount combinations
			// in a simpler to use interface. The constructor we use here creates and initializes
			// a single token color transfer in a single statement. The actual color and amount
			// values passed in will be stored in a new map on the host.
			transfers := wasmlib.ScTransferFromBaseTokens(payout)

			// Perform the actual transfer of tokens from the smart contract to the address
			// of the winner. The transfer_to_address() method receives the address value and
			// the proxy to the new transfers map on the host, and will call the corresponding
			// host sandbox function with these values.
			ctx.Send(bet.Better.Address(), transfers)
		}

		// Announce who got sent what as event.
		f.Events.Payout(bet.Better.Address(), payout)
	}

	// This is where we transfer the remainder after payout to the creator of the smart contract.
	// The bank always wins :-P
	remainder := totalBetAmount - totalPayout
	if remainder != 0 {
		// We have a remainder. First create a transfer for the remainder.
		transfers := wasmlib.ScTransferFromBaseTokens(remainder)

		// Send the remainder to the contract owner.
		ctx.Send(f.State.Owner().Value().Address(), transfers)
	}

	// Set round status to 0, send out event to notify that the round has ended
	f.State.RoundStatus().SetValue(0)
	f.Events.Stop()
}

func funcForceReset(_ wasmlib.ScFuncContext, f *ForceResetContext) {
	// Get the 'bets' array in state storage.
	bets := f.State.Bets()

	// Clear all bets.
	bets.Clear()

	// Set round status to 0, send out event to notify that the round has ended
	f.State.RoundStatus().SetValue(0)
	f.Events.Stop()
}

// 'playPeriod' can be used by the contract creator to set the length of a betting round
// to a different value than the default value, which is 120 seconds.
func funcPlayPeriod(ctx wasmlib.ScFuncContext, f *PlayPeriodContext) {
	// Since we are sure that the 'playPeriod' parameter actually exists we can
	// retrieve its actual value into an i32 value.
	playPeriod := f.Params.PlayPeriod().Value()

	// Require that the play period (in seconds) is not ridiculously low.
	// Otherwise, panic out with an error message.
	ctx.Require(playPeriod >= 10, "invalid play period")

	// Now we set the corresponding variable 'playPeriod' in state storage.
	f.State.PlayPeriod().SetValue(playPeriod)
}

func viewLastWinningNumber(_ wasmlib.ScViewContext, f *LastWinningNumberContext) {
	// Get the 'lastWinningNumber' int64 value from state storage.
	lastWinningNumber := f.State.LastWinningNumber().Value()

	// Set the 'lastWinningNumber' in results to the value from state storage.
	f.Results.LastWinningNumber().SetValue(lastWinningNumber)
}

func viewRoundNumber(_ wasmlib.ScViewContext, f *RoundNumberContext) {
	// Get the 'roundNumber' int64 value from state storage.
	roundNumber := f.State.RoundNumber().Value()

	// Set the 'roundNumber' in results to the value from state storage.
	f.Results.RoundNumber().SetValue(roundNumber)
}

func viewRoundStatus(_ wasmlib.ScViewContext, f *RoundStatusContext) {
	// Get the 'roundStatus' int16 value from state storage.
	roundStatus := f.State.RoundStatus().Value()

	// Set the 'roundStatus' in results to the value from state storage.
	f.Results.RoundStatus().SetValue(roundStatus)
}

func viewRoundStartedAt(_ wasmlib.ScViewContext, f *RoundStartedAtContext) {
	// Get the 'roundStartedAt' int32 value from state storage.
	roundStartedAt := f.State.RoundStartedAt().Value()

	// Set the 'roundStartedAt' in results to the value from state storage.
	f.Results.RoundStartedAt().SetValue(roundStartedAt)
}

func funcForcePayout(ctx wasmlib.ScFuncContext, _ *ForcePayoutContext) {
	fairroulette.ScFuncs.PayWinners(ctx).Func.Call()
}

func funcInit(ctx wasmlib.ScFuncContext, f *InitContext) {
	if f.Params.Owner().Exists() {
		f.State.Owner().SetValue(f.Params.Owner().Value())
		return
	}
	f.State.Owner().SetValue(ctx.RequestSender())
}
