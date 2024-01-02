// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmsolo

import (
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/wasp/packages/cryptolib"
	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/isc/coreutil"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmclient/go/wasmclient"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmhost"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmrequests"
	"github.com/iotaledger/wasp/packages/wasmvm/wasmlib/go/wasmlib/wasmtypes"
)

const (
	SoloDebug        = false
	SoloHostTracing  = false
	SoloStackTracing = false
)

var (
	// NoWasm / GoWasm / RsWasm / TsWasm are used to specify the Wasm language mode,
	// By default, SoloContext will try to run the Go SC code directly (NoWasm)
	// The 3 other flags can be used to cause Wasm code to be loaded and run instead.
	// They are checked in sequence and the first one set determines the Wasm language used.
	NoWasm = flag.Bool("nowasm", false, "use Go Wasm smart contract code without Wasm")
	GoWasm = flag.Bool("gowasm", false, "use Go Wasm smart contract code")
	RsWasm = flag.Bool("rswasm", false, "use Rust Wasm smart contract code")
	TsWasm = flag.Bool("tswasm", false, "use TypeScript Wasm smart contract code")

	// UseWasmEdge flag is kept here in case we decide to use WasmEdge again. Some tests
	// refer to this flag, so we keep it here instead of having to comment out a bunch
	// of code. To actually enable WasmEdge you need to uncomment the relevant lines in
	// NewSoloContextForChain(), and remove the go:build directives from wasmedge.go, so
	// that the linker can actually pull in the WasmEdge runtime.
	UseWasmEdge = flag.Bool("wasmedge", false, "use WasmEdge instead of WasmTime")
)

const (
	L2FundsAgent      = 10 * isc.Million
	L2FundsContract   = 10 * isc.Million
	L2FundsCreator    = 20 * isc.Million
	L2FundsOriginator = 30 * isc.Million
)

type SoloContext struct {
	Chain          *solo.Chain
	Cvt            wasmhost.WasmConvertor
	creator        *SoloAgent
	StorageDeposit uint64
	Err            error
	Gas            uint64
	GasFee         uint64
	Hprog          hashing.HashValue
	isRequest      bool
	IsWasm         bool
	keyPair        *cryptolib.KeyPair
	nfts           map[iotago.NFTID]*isc.NFT
	offLedger      bool
	scName         string
	Tx             *iotago.Transaction
	wc             *wasmhost.WasmContext
}

var (
	_ wasmlib.ScFuncClientContext = new(SoloContext)
	_ wasmlib.ScViewClientContext = new(SoloContext)
)

func contains(s []isc.AgentID, e isc.AgentID) bool {
	for _, a := range s {
		if a.Equals(e) {
			return true
		}
	}
	return false
}

// NewSoloContext can be used to create a SoloContext associated with a smart contract
// with minimal information and will verify successful creation before returning ctx.
// It will start a default chain "chain1" before initializing the smart contract.
// It takes the scName and onLoad() function associated with the contract.
// Optionally, an init.Func that has been initialized with the parameters to pass to
// the contract's init() function can be specified.
// Unless you want to use a different chain than the default "chain1" this will be your
// function of choice to set up a smart contract for your tests
func NewSoloContext(t testing.TB, scName string, onLoad wasmhost.ScOnloadFunc, init ...*wasmlib.ScInitFunc) *SoloContext {
	ctx := NewSoloContextForChain(t, nil, nil, scName, onLoad, init...)
	require.NoError(t, ctx.Err)
	return ctx
}

// NewSoloContextForChain can be used to create a SoloContext associated with a smart contract
// on a particular chain.  When chain is nil the function will start a default chain "chain1"
// before initializing the smart contract.
// When creator is nil the creator will be the chain originator
// It takes the scName and onLoad() function associated with the contract.
// Optionally, an init.Func that has been initialized with the parameters to pass to
// the contract's init() function can be specified.
// You can check for any error that occurred by checking the ctx.Err member.
func NewSoloContextForChain(t testing.TB, chain *solo.Chain, creator *SoloAgent, scName string,
	onLoad wasmhost.ScOnloadFunc, init ...*wasmlib.ScInitFunc,
) *SoloContext {
	ctx := soloContext(t, chain, scName, creator)

	ctx.Balances()

	var keyPair *cryptolib.KeyPair
	if creator != nil {
		keyPair = creator.Pair
		chain.MustDepositBaseTokensToL2(L2FundsCreator, creator.Pair)
	}
	ctx.uploadWasm(keyPair)
	if ctx.Err != nil {
		return ctx
	}

	ctx.Balances()

	var params []interface{}
	if len(init) != 0 {
		params = init[0].Params()
	}
	if !ctx.IsWasm {
		wasmhost.GoWasmVM = func() wasmhost.WasmVM {
			return wasmhost.NewWasmGoVM(ctx.scName, onLoad)
		}
	}
	//if ctx.IsWasm && *UseWasmEdge && wasmproc.GoWasmVM == nil {
	//	wasmproc.GoWasmVM = wasmhost.NewWasmEdgeVM
	//}
	ctx.Err = ctx.Chain.DeployContract(keyPair, ctx.scName, ctx.Hprog, params...)
	if !ctx.IsWasm {
		// just in case deploy failed we don't want to leave this around
		wasmhost.GoWasmVM = nil
	}
	if ctx.Err != nil {
		return ctx
	}

	ctx.Balances()

	scAccount := isc.NewContractAgentID(ctx.Chain.ChainID, isc.Hn(scName))
	ctx.Err = ctx.Chain.SendFromL1ToL2AccountBaseTokens(0, L2FundsContract, scAccount, ctx.Creator().Pair)

	ctx.Balances()

	if ctx.Err != nil {
		return ctx
	}
	return ctx.init(onLoad)
}

// NewSoloContextForNative can be used to create a SoloContext associated with a native smart contract
// on a particular chain. When chain is nil the function will start a default chain "chain1" before
// deploying and initializing the smart contract.
// When creator is nil the creator will be the chain originator
// It takes the scName, onLoad() function, and processor associated with the contract.
// Optionally, an init.Func that has been initialized with the parameters to pass to
// the contract's init() function can be specified.
// You can check for any error that occurred by checking the ctx.Err member.
func NewSoloContextForNative(t testing.TB, chain *solo.Chain, creator *SoloAgent, scName string, onLoad wasmhost.ScOnloadFunc,
	proc *coreutil.ContractProcessor, init ...*wasmlib.ScInitFunc,
) *SoloContext {
	ctx := soloContext(t, chain, scName, creator)
	ctx.Chain.Env.WithNativeContract(proc)
	ctx.Hprog = proc.Contract.ProgramHash

	var keyPair *cryptolib.KeyPair
	if creator != nil {
		keyPair = creator.Pair
		chain.MustDepositBaseTokensToL2(L2FundsCreator, creator.Pair)
	}
	var params []interface{}
	if len(init) != 0 {
		params = init[0].Params()
	}
	ctx.Err = ctx.Chain.DeployContract(keyPair, scName, ctx.Hprog, params...)
	if ctx.Err != nil {
		return ctx
	}

	scAccount := isc.NewContractAgentID(ctx.Chain.ChainID, isc.Hn(scName))
	ctx.Err = ctx.Chain.SendFromL1ToL2AccountBaseTokens(0, L2FundsContract, scAccount, ctx.Creator().Pair)
	if ctx.Err != nil {
		return ctx
	}

	return ctx.init(onLoad)
}

func soloContext(t testing.TB, chain *solo.Chain, scName string, creator *SoloAgent) *SoloContext {
	if chain == nil {
		chain = StartChain(t, "chain1")
	}
	err := wasmclient.SetSandboxWrappers(chain.ChainID.String())
	if err != nil {
		panic(err)
	}
	return &SoloContext{
		scName:         scName,
		Chain:          chain,
		creator:        creator,
		StorageDeposit: wasmlib.StorageDeposit,
	}
}

// StartChain starts a new chain named chainName.
func StartChain(t testing.TB, chainName string, env ...*solo.Solo) *solo.Chain {
	if SoloDebug {
		// avoid pesky timeouts during debugging
		wasmhost.DisableWasmTimeout = true
	}
	wasmhost.HostTracing = SoloHostTracing

	var soloEnv *solo.Solo
	if len(env) != 0 {
		soloEnv = env[0]
	}
	if soloEnv == nil {
		soloEnv = solo.New(t, &solo.InitOptions{
			Debug:                    SoloDebug,
			PrintStackTrace:          SoloStackTracing,
			AutoAdjustStorageDeposit: true,
		})
	}
	chain, _ := soloEnv.NewChainExt(nil, 0, chainName)
	chain.MustDepositBaseTokensToL2(L2FundsOriginator, chain.OriginatorPrivateKey)
	return chain
}

// Account returns a SoloAgent for the smart contract associated with ctx
func (ctx *SoloContext) Account() *SoloAgent {
	agentID := isc.NewContractAgentID(ctx.Chain.ChainID, isc.Hn(ctx.scName))
	return &SoloAgent{
		agentID: agentID,
		Env:     ctx.Chain.Env,
		ID:      agentID.String(),
		Name:    ctx.Chain.Name + "." + ctx.scName,
		Pair:    nil,
	}
}

func (ctx *SoloContext) AccountID() wasmtypes.ScAgentID {
	return ctx.Account().ScAgentID()
}

// AdvanceClockBy is used to forward the internal clock by the provided step duration.
func (ctx *SoloContext) AdvanceClockBy(step time.Duration) {
	ctx.Chain.Env.AdvanceClockBy(step)
}

// Balance returns the account balance of the specified agent on the chain associated with ctx.
// The optional nativeTokenID parameter can be used to retrieve the balance for the specific token.
// When nativeTokenID is omitted, the base tokens balance is assumed.
func (ctx *SoloContext) Balance(agent *SoloAgent, nativeTokenID ...wasmtypes.ScTokenID) uint64 {
	account := agent.AgentID()
	switch len(nativeTokenID) {
	case 0:
		baseTokens := ctx.Chain.L2BaseTokens(account)
		return baseTokens
	case 1:
		token := cvt.IscTokenID(&nativeTokenID[0])
		tokens := ctx.Chain.L2NativeTokens(account, token).Uint64()
		return tokens
	default:
		require.Fail(ctx.Chain.Env.T, "too many nativeTokenID arguments")
		return 0
	}
}

// Balances prints all known accounts, both L2 and L1.
// It uses the L2 ledger to enumerate the known accounts.
// Any newly created SoloAgents can be specified as extra accounts
func (ctx *SoloContext) Balances(agents ...*SoloAgent) *SoloBalances {
	return NewSoloBalances(ctx, agents...)
}

// CommonAccount returns a SoloAgent for the chain associated with ctx
func (ctx *SoloContext) CommonAccount() *SoloAgent {
	agentID := accounts.CommonAccount()
	return &SoloAgent{
		agentID: agentID,
		Env:     ctx.Chain.Env,
		ID:      agentID.String(),
		Name:    ctx.Chain.Name + ".Common",
		Pair:    nil,
	}
}

func (ctx *SoloContext) ChainOwnerID() wasmtypes.ScAgentID {
	return cvt.ScAgentID(ctx.Chain.OriginatorAgentID)
}

// ClientContract is a function that is required to use SoloContext as an ScViewClientContext
func (ctx *SoloContext) ClientContract(hContract wasmtypes.ScHname) wasmtypes.ScHname {
	_ = hContract
	return cvt.ScHname(isc.Hn(ctx.scName))
}

// ContractExists checks to see if the contract named scName exists in the chain associated with ctx.
func (ctx *SoloContext) ContractExists(scName string) error {
	_, err := ctx.Chain.FindContract(scName)
	return err
}

// Creator returns a SoloAgent representing the contract creator
func (ctx *SoloContext) Creator() *SoloAgent {
	if ctx.creator != nil {
		return ctx.creator
	}
	return ctx.Originator()
}

func (ctx *SoloContext) CurrentChainID() wasmtypes.ScChainID {
	return cvt.ScChainID(ctx.Chain.ChainID)
}

func (ctx *SoloContext) EnqueueRequest() {
	ctx.isRequest = true
}

func (ctx *SoloContext) existFile(lang string) string {
	fileName := ctx.scName + "_" + lang + ".wasm"
	pathName := "../../" + ctx.scName + "/" + lang + "/pkg/" + fileName
	if lang == "bg" {
		pathName = "../../" + ctx.scName + "/rs/" + ctx.scName + "wasm/pkg/" + ctx.scName + "wasm_bg.wasm"
	}

	// first check for new file in common build path
	exists, _ := util.ExistsFilePath(pathName)
	if exists {
		return pathName
	}

	// then check for file in current folder
	exists, _ = util.ExistsFilePath(fileName)
	if exists {
		return fileName
	}

	// file not found
	return ""
}

func (ctx *SoloContext) FnCall(req *wasmrequests.CallRequest) []byte {
	return NewSoloSandbox(ctx).FnCall(req)
}

func (ctx *SoloContext) FnChainID() wasmtypes.ScChainID {
	return ctx.CurrentChainID()
}

func (ctx *SoloContext) FnPost(req *wasmrequests.PostRequest) []byte {
	return NewSoloSandbox(ctx).FnPost(req)
}

func (ctx *SoloContext) Host() wasmlib.ScHost {
	return nil
}

// init further initializes the SoloContext.
func (ctx *SoloContext) init(onLoad wasmhost.ScOnloadFunc) *SoloContext {
	ctx.wc = wasmhost.NewWasmContextForSoloContext("-solo-", NewSoloSandbox(ctx))
	onLoad(-2).Export(ctx.wc.ExportName)
	return ctx
}

// MintNFT tells SoloContext to mint a new NFT issued/owned by the specified agent
// note that SoloContext will cache the NFT data to be able to use it
// in Post()s that go through the *SAME* SoloContext
func (ctx *SoloContext) MintNFT(agent *SoloAgent, metadata []byte) wasmtypes.ScNftID {
	addr, ok := isc.AddressFromAgentID(agent.AgentID())
	if !ok {
		ctx.Err = errors.New("agent should be an address")
		return wasmtypes.NftIDFromBytes(nil)
	}
	nft, _, err := ctx.Chain.Env.MintNFTL1(agent.Pair, addr, metadata)
	if err != nil {
		ctx.Err = err
		return wasmtypes.NftIDFromBytes(nil)
	}
	if ctx.nfts == nil {
		ctx.nfts = make(map[iotago.NFTID]*isc.NFT)
	}
	ctx.nfts[nft.ID] = nft
	return cvt.ScNftID(&nft.ID)
}

// NewSoloAgent creates a new SoloAgent with utxodb.FundsFromFaucetAmount (1 Gi)
// tokens in its address and pre-deposits 10Mi into the corresponding chain account
func (ctx *SoloContext) NewSoloAgent(name string) *SoloAgent {
	agent := NewSoloAgent(ctx.Chain.Env, name)
	ctx.Chain.MustDepositBaseTokensToL2(L2FundsAgent+wasmlib.MinGasFee, agent.Pair)
	return agent
}

// NewSoloFoundry creates a new SoloFoundry
func (ctx *SoloContext) NewSoloFoundry(maxSupply interface{}, agent ...*SoloAgent) (*SoloFoundry, error) {
	return NewSoloFoundry(ctx, maxSupply, agent...)
}

// NFTs returns the list of NFTs in the account of the specified agent on
// the chain associated with ctx.
func (ctx *SoloContext) NFTs(agent *SoloAgent) []wasmtypes.ScNftID {
	account := agent.AgentID()
	l2nfts := ctx.Chain.L2NFTs(account)
	nfts := make([]wasmtypes.ScNftID, 0, len(l2nfts))
	for _, l2nft := range l2nfts {
		theNft := l2nft
		nfts = append(nfts, cvt.ScNftID(&theNft))
	}
	return nfts
}

// OffLedger tells SoloContext to Post() the next request off-ledger
func (ctx *SoloContext) OffLedger(agent *SoloAgent) wasmlib.ScFuncClientContext {
	ctx.offLedger = true
	ctx.keyPair = agent.Pair
	return ctx
}

// Originator returns a SoloAgent representing the chain originator
func (ctx *SoloContext) Originator() *SoloAgent {
	agentID := ctx.Chain.OriginatorAgentID
	return &SoloAgent{
		agentID: agentID,
		Env:     ctx.Chain.Env,
		ID:      agentID.String(),
		Name:    ctx.Chain.Name + ".Originator",
		Pair:    ctx.Chain.OriginatorPrivateKey,
	}
}

// Sign is used to force a different agent for signing a Post() request
func (ctx *SoloContext) Sign(agent *SoloAgent) wasmlib.ScFuncClientContext {
	ctx.keyPair = agent.Pair
	return ctx
}

func (ctx *SoloContext) SoloContextForCore(t testing.TB, scName string, onLoad wasmhost.ScOnloadFunc) *SoloContext {
	return soloContext(t, ctx.Chain, scName, nil).init(onLoad)
}

func (ctx *SoloContext) UpdateGasFees() {
	receipt := ctx.Chain.LastReceipt()
	if receipt == nil {
		panic("UpdateGasFees: missing last receipt")
	}
	ctx.Gas = receipt.GasBurned
	ctx.GasFee = receipt.GasFeeCharged
}

func (ctx *SoloContext) uploadWasm(keyPair *cryptolib.KeyPair) {
	// default to use WasmGoVM to run Go SC code directly without Wasm VM
	wasmFile := "go"
	switch {
	case *NoWasm:
		// explicit default
	case *GoWasm:
		// find Go Wasm file
		wasmFile = ctx.existFile("go")
	case *RsWasm:
		// find Rust Wasm file
		wasmFile = ctx.existFile("bg")
	case *TsWasm:
		// find TypeScript Wasm file
		wasmFile = ctx.existFile("ts")
	}
	if wasmFile == "" {
		panic("cannot find Wasm file for: " + ctx.scName)
	}

	if wasmFile == "go" {
		// none of the Wasm modes selected, use WasmGoVM to run Go SC code directly
		ctx.Hprog, ctx.Err = ctx.Chain.UploadWasm(keyPair, []byte("go:"+ctx.scName))
		return
	}

	// upload the Wasm code into the core blob contract
	ctx.Hprog, ctx.Err = ctx.Chain.UploadWasmFromFile(keyPair, wasmFile)
	ctx.IsWasm = true
}

// WaitForPendingRequests waits for expectedRequests requests to be processed
// since the last call to WaitForPendingRequestsMark().
// The function will wait for maxWait (default 5 seconds per request) duration
// before giving up with a timeout. The function returns false in case of a timeout.
func (ctx *SoloContext) WaitForPendingRequests(expectedRequests int, maxWait ...time.Duration) bool {
	timeout := time.Duration(expectedRequests*5) * time.Second
	if len(maxWait) > 0 {
		timeout = maxWait[0]
	}
	return ctx.Chain.WaitForRequestsThrough(expectedRequests, timeout)
}

// WaitForPendingRequestsMark marks the current InPoolCounter to be used by
// a subsequent call to WaitForPendingRequests()
func (ctx *SoloContext) WaitForPendingRequestsMark() {
	ctx.Chain.WaitForRequestsMark()
}
