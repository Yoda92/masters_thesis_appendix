package vmtxbuilder

import (
	"fmt"
	"math/big"
	"sort"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/parameters"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm"
)

func (txb *AnchorTransactionBuilder) CreateNewFoundry(
	scheme iotago.TokenScheme,
	metadata []byte,
) (uint32, uint64) {
	// must check for valid token scheme to not create invalid transactions
	simpleTokenScheme := util.MustTokenScheme(scheme)
	maxSupply := simpleTokenScheme.MaximumSupply
	if maxSupply.Cmp(util.Big0) <= 0 {
		panic(vm.ErrCreateFoundryMaxSupplyMustBePositive)
	}
	if maxSupply.Cmp(util.MaxUint256) > 0 {
		panic(vm.ErrCreateFoundryMaxSupplyTooBig)
	}

	f := &iotago.FoundryOutput{
		Amount:       0,
		NativeTokens: nil,
		SerialNumber: txb.nextFoundrySerialNumber(),
		TokenScheme:  scheme,
		Conditions: iotago.UnlockConditions{
			&iotago.ImmutableAliasUnlockCondition{
				Address: util.AliasIDFromAliasOutput(txb.anchorOutput, txb.anchorOutputID).ToAddress().(*iotago.AliasAddress),
			},
		},
		Features: nil,
	}
	if len(metadata) > 0 {
		f.Features = iotago.Features{&iotago.MetadataFeature{
			Data: metadata,
		}}
	}
	f.Amount = parameters.L1().Protocol.RentStructure.MinRent(f)
	txb.invokedFoundries[f.SerialNumber] = &foundryInvoked{
		serialNumber:     f.SerialNumber,
		accountingInput:  nil,
		accountingOutput: f,
	}
	return f.SerialNumber, f.Amount
}

// ModifyNativeTokenSupply inflates the supply is delta > 0, shrinks if delta < 0
// returns adjustment of the storage deposit.
func (txb *AnchorTransactionBuilder) ModifyNativeTokenSupply(nativeTokenID iotago.NativeTokenID, delta *big.Int) int64 {
	sn := nativeTokenID.FoundrySerialNumber()
	f := txb.ensureFoundry(sn)
	if f == nil {
		panic(vm.ErrFoundryDoesNotExist)
	}
	// check if the loaded foundry matches the nativeTokenID
	if nativeTokenID != f.accountingInput.MustNativeTokenID() {
		panic(fmt.Errorf("%v: requested token ID: %s, foundry token id: %s",
			vm.ErrCantModifySupplyOfTheToken, nativeTokenID.String(), f.accountingInput.MustNativeTokenID().String()))
	}

	defer txb.mustCheckTotalNativeTokensExceeded()

	simpleTokenScheme := util.MustTokenScheme(f.accountingOutput.TokenScheme)

	// check the supply bounds
	var newMinted, newMelted *big.Int
	if delta.Cmp(util.Big0) >= 0 {
		newMinted = big.NewInt(0).Add(simpleTokenScheme.MintedTokens, delta)
		newMelted = simpleTokenScheme.MeltedTokens
	} else {
		newMinted = simpleTokenScheme.MintedTokens
		newMelted = big.NewInt(0).Sub(simpleTokenScheme.MeltedTokens, delta)
	}
	if newMinted.Cmp(util.Big0) < 0 || newMinted.Cmp(simpleTokenScheme.MaximumSupply) > 0 {
		panic(vm.ErrNativeTokenSupplyOutOffBounds)
	}
	// accrue/adjust this token balance in the internal outputs
	adjustment := txb.addNativeTokenBalanceDelta(nativeTokenID, delta)
	// update the supply and foundry record in the builder
	simpleTokenScheme.MintedTokens = newMinted
	simpleTokenScheme.MeltedTokens = newMelted
	txb.invokedFoundries[sn] = f

	adjustment += int64(f.accountingInput.Amount) - int64(f.accountingOutput.Amount)
	return adjustment
}

func (txb *AnchorTransactionBuilder) ensureFoundry(sn uint32) *foundryInvoked {
	if foundryOutput, exists := txb.invokedFoundries[sn]; exists {
		return foundryOutput
	}

	// load foundry output from the state
	foundryOutput, outputID := txb.accountsView.FoundryOutput(sn)
	if foundryOutput == nil {
		return nil
	}
	f := &foundryInvoked{
		serialNumber:      foundryOutput.SerialNumber,
		accountingInputID: outputID,
		accountingInput:   foundryOutput,
		accountingOutput:  cloneFoundryOutput(foundryOutput),
	}
	txb.invokedFoundries[sn] = f
	return f
}

// DestroyFoundry destroys existing foundry. Return storage deposit
func (txb *AnchorTransactionBuilder) DestroyFoundry(sn uint32) uint64 {
	f := txb.ensureFoundry(sn)
	if f == nil {
		panic(vm.ErrFoundryDoesNotExist)
	}
	if f.accountingInput == nil {
		panic(vm.ErrCantDestroyFoundryBeingCreated)
	}

	defer txb.mustCheckTotalNativeTokensExceeded()

	f.accountingOutput = nil
	return f.accountingInput.Amount
}

func (txb *AnchorTransactionBuilder) nextFoundrySerialNumber() uint32 {
	return txb.nextFoundryCounter() + 1
}

func (txb *AnchorTransactionBuilder) nextFoundryCounter() uint32 {
	numNew := uint32(0)
	for _, f := range txb.invokedFoundries {
		if f.isNewCreated() {
			numNew++
		}
	}
	return txb.anchorOutput.FoundryCounter + numNew
}

func (txb *AnchorTransactionBuilder) foundriesSorted() []*foundryInvoked {
	ret := make([]*foundryInvoked, 0, len(txb.invokedFoundries))
	for _, f := range txb.invokedFoundries {
		if !f.requiresExistingAccountingUTXOAsInput() && !f.producesAccountingOutput() {
			continue
		}
		ret = append(ret, f)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].serialNumber < ret[j].serialNumber
	})
	return ret
}

func (txb *AnchorTransactionBuilder) FoundriesToBeUpdated() ([]uint32, []uint32) {
	toBeUpdated := make([]uint32, 0, len(txb.invokedFoundries))
	toBeRemoved := make([]uint32, 0, len(txb.invokedFoundries))
	for _, f := range txb.foundriesSorted() {
		if f.producesAccountingOutput() {
			toBeUpdated = append(toBeUpdated, f.serialNumber)
		} else if f.requiresExistingAccountingUTXOAsInput() {
			toBeRemoved = append(toBeRemoved, f.serialNumber)
		}
	}
	return toBeUpdated, toBeRemoved
}

func (txb *AnchorTransactionBuilder) FoundryOutputsBySN(serNums []uint32) map[uint32]*iotago.FoundryOutput {
	ret := make(map[uint32]*iotago.FoundryOutput)
	for _, sn := range serNums {
		ret[sn] = txb.invokedFoundries[sn].accountingOutput
	}
	return ret
}

type foundryInvoked struct {
	serialNumber      uint32
	accountingInputID iotago.OutputID       // if in != nil
	accountingInput   *iotago.FoundryOutput // nil if created
	accountingOutput  *iotago.FoundryOutput // nil if destroyed
}

func (f *foundryInvoked) Clone() *foundryInvoked {
	outputID := iotago.OutputID{}
	copy(outputID[:], f.accountingInputID[:])

	return &foundryInvoked{
		serialNumber:      f.serialNumber,
		accountingInputID: outputID,
		accountingInput:   cloneFoundryOutput(f.accountingInput),
		accountingOutput:  cloneFoundryOutput(f.accountingOutput),
	}
}

func (f *foundryInvoked) isNewCreated() bool {
	return !f.requiresExistingAccountingUTXOAsInput() && f.producesAccountingOutput()
}

func (f *foundryInvoked) requiresExistingAccountingUTXOAsInput() bool {
	if f.accountingInput == nil {
		return false
	}
	if identicalFoundries(f.accountingInput, f.accountingOutput) {
		return false
	}
	return true
}

func (f *foundryInvoked) producesAccountingOutput() bool {
	if f.accountingOutput == nil {
		return false
	}
	if identicalFoundries(f.accountingInput, f.accountingOutput) {
		return false
	}
	return true
}

func cloneFoundryOutput(f *iotago.FoundryOutput) *iotago.FoundryOutput {
	if f == nil {
		return nil
	}
	return f.Clone().(*iotago.FoundryOutput)
}

// identicalFoundries assumes use case and does consistency checks
func identicalFoundries(f1, f2 *iotago.FoundryOutput) bool {
	if f1 == nil || f2 == nil {
		return false
	}
	simpleTokenSchemeF1 := util.MustTokenScheme(f1.TokenScheme)
	simpleTokenSchemeF2 := util.MustTokenScheme(f2.TokenScheme)

	switch {
	case f1 == f2:
		return true
	case f1.SerialNumber != f2.SerialNumber:
		return false
	case simpleTokenSchemeF1.MintedTokens.Cmp(simpleTokenSchemeF2.MintedTokens) != 0:
		return false
	case simpleTokenSchemeF1.MeltedTokens.Cmp(simpleTokenSchemeF2.MeltedTokens) != 0:
		return false
	case f1.Amount != f2.Amount:
		panic("identicalFoundries: inconsistency, amount is assumed immutable")
	case len(f1.NativeTokens) > 0 || len(f2.NativeTokens) > 0:
		panic("identicalFoundries: inconsistency, foundry is not expected not contain native tokens")
	case simpleTokenSchemeF1.MaximumSupply.Cmp(simpleTokenSchemeF2.MaximumSupply) != 0:
		panic("identicalFoundries: inconsistency, maximum supply is immutable")
	case !f1.Ident().Equal(f2.Ident()):
		panic("identicalFoundries: inconsistency, addresses must always be equal")
	case !equalTokenScheme(simpleTokenSchemeF1, simpleTokenSchemeF2):
		panic("identicalFoundries: inconsistency, if serial numbers are equal, token schemes must be equal")
	case len(f1.Features) != 0 || len(f2.Features) != 0:
		panic("identicalFoundries: inconsistency, feat blocks are not expected in the foundry")
	}
	return true
}

func equalTokenScheme(a, b *iotago.SimpleTokenScheme) bool {
	return a.MintedTokens.Cmp(b.MintedTokens) == 0 &&
		a.MeltedTokens.Cmp(b.MeltedTokens) == 0 &&
		a.MaximumSupply.Cmp(b.MaximumSupply) == 0
}
