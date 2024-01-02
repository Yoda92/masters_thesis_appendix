package models

import (
	"time"

	"github.com/iotaledger/wasp/packages/chain"
	"github.com/iotaledger/wasp/packages/isc"
	"github.com/iotaledger/wasp/packages/webapi/dto"
)

type MetricItem[T interface{}] struct {
	Messages    uint32    `json:"messages" swagger:"required,min(1)"`
	Timestamp   time.Time `json:"timestamp" swagger:"required"`
	LastMessage T         `json:"lastMessage" swagger:"required"`
}

/*
ChainMetrics
Echo Swagger does not support generics such as MetricItem[Foo]
Creating separate types works instead.
To not create a mapper for each type, the actual service remains using MetricItem[Foo] but this model here is presented to the docs.
This can be removed if we change to swag/echo-swagger
*/
type (
	AliasOutputMetricItem         MetricItem[*Output]
	OnLedgerRequestMetricItem     MetricItem[*OnLedgerRequest]
	InOutputMetricItem            MetricItem[*InOutput]
	InStateOutputMetricItem       MetricItem[*InStateOutput]
	TxInclusionStateMsgMetricItem MetricItem[*TxInclusionStateMsg]
	TransactionMetricItem         MetricItem[*Transaction]
	TransactionIDMetricItem       MetricItem[*Transaction]
	UTXOInputMetricItem           MetricItem[*OutputID]
	InterfaceMetricItem           MetricItem[interface{}]
	PublisherStateTransactionItem MetricItem[*StateTransaction]
	RegisteredChainIDItems        []string
	MilestoneMetricItem           MetricItem[*MilestoneInfo]
)

type ChainMessageMetrics struct {
	InStateOutput      InStateOutputMetricItem       `json:"inStateOutput" swagger:"required"`
	InAliasOutput      AliasOutputMetricItem         `json:"inAliasOutput" swagger:"required"`
	InOutput           InOutputMetricItem            `json:"inOutput" swagger:"required"`
	InOnLedgerRequest  OnLedgerRequestMetricItem     `json:"inOnLedgerRequest" swagger:"required"`
	InTxInclusionState TxInclusionStateMsgMetricItem `json:"inTxInclusionState" swagger:"required"`

	OutPublishStateTransaction      PublisherStateTransactionItem `json:"outPublisherStateTransaction" swagger:"required"`
	OutPublishGovernanceTransaction TransactionMetricItem         `json:"outPublishGovernanceTransaction" swagger:"required"`
	OutPullLatestOutput             InterfaceMetricItem           `json:"outPullLatestOutput" swagger:"required"`
	OutPullTxInclusionState         TransactionIDMetricItem       `json:"outPullTxInclusionState" swagger:"required"`
	OutPullOutputByID               UTXOInputMetricItem           `json:"outPullOutputByID" swagger:"required"`
}

type NodeMessageMetrics struct {
	RegisteredChainIDs RegisteredChainIDItems `json:"registeredChainIDs" swagger:"required"`

	InMilestone        MilestoneMetricItem           `json:"inMilestone" swagger:"required"`
	InStateOutput      InStateOutputMetricItem       `json:"inStateOutput" swagger:"required"`
	InAliasOutput      AliasOutputMetricItem         `json:"inAliasOutput" swagger:"required"`
	InOutput           InOutputMetricItem            `json:"inOutput" swagger:"required"`
	InOnLedgerRequest  OnLedgerRequestMetricItem     `json:"inOnLedgerRequest" swagger:"required"`
	InTxInclusionState TxInclusionStateMsgMetricItem `json:"inTxInclusionState" swagger:"required"`

	OutPublishStateTransaction      PublisherStateTransactionItem `json:"outPublisherStateTransaction" swagger:"required"`
	OutPublishGovernanceTransaction TransactionMetricItem         `json:"outPublishGovernanceTransaction" swagger:"required"`
	OutPullLatestOutput             InterfaceMetricItem           `json:"outPullLatestOutput" swagger:"required"`
	OutPullTxInclusionState         TransactionIDMetricItem       `json:"outPullTxInclusionState" swagger:"required"`
	OutPullOutputByID               UTXOInputMetricItem           `json:"outPullOutputByID" swagger:"required"`
}

func MapMetricItem[T any, G any](metrics *dto.MetricItem[G], value T) MetricItem[T] {
	return MetricItem[T]{
		Messages:    metrics.Messages,
		Timestamp:   metrics.Timestamp,
		LastMessage: value,
	}
}

func MapRegisteredChainIDs(registered []isc.ChainID) []string {
	chainIDs := make([]string, len(registered))

	for k, v := range registered {
		chainIDs[k] = v.String()
	}

	return chainIDs
}

func MapChainMessageMetrics(metrics *dto.ChainMessageMetrics) *ChainMessageMetrics {
	return &ChainMessageMetrics{
		InStateOutput:      InStateOutputMetricItem(MapMetricItem(metrics.InStateOutput, InStateOutputFromISCInStateOutput(metrics.InStateOutput.LastMessage))),
		InAliasOutput:      AliasOutputMetricItem(MapMetricItem(metrics.InAliasOutput, OutputFromIotaGoOutput(metrics.InAliasOutput.LastMessage))),
		InOutput:           InOutputMetricItem(MapMetricItem(metrics.InOutput, InOutputFromISCInOutput(metrics.InOutput.LastMessage))),
		InOnLedgerRequest:  OnLedgerRequestMetricItem(MapMetricItem(metrics.InOnLedgerRequest, OnLedgerRequestFromISC(metrics.InOnLedgerRequest.LastMessage))),
		InTxInclusionState: TxInclusionStateMsgMetricItem(MapMetricItem(metrics.InTxInclusionState, TxInclusionStateMsgFromISCTxInclusionStateMsg(metrics.InTxInclusionState.LastMessage))),

		OutPublishStateTransaction:      PublisherStateTransactionItem(MapMetricItem(metrics.OutPublishStateTransaction, StateTransactionFromISCStateTransaction(metrics.OutPublishStateTransaction.LastMessage))),
		OutPublishGovernanceTransaction: TransactionMetricItem(MapMetricItem(metrics.OutPublishGovernanceTransaction, TransactionFromIotaGoTransaction(metrics.OutPublishGovernanceTransaction.LastMessage))),
		OutPullLatestOutput:             InterfaceMetricItem(MapMetricItem(metrics.OutPullLatestOutput, metrics.OutPullLatestOutput.LastMessage)),
		OutPullTxInclusionState:         TransactionIDMetricItem(MapMetricItem(metrics.OutPullTxInclusionState, TransactionFromIotaGoTransactionID(&metrics.OutPullTxInclusionState.LastMessage))),
		OutPullOutputByID:               UTXOInputMetricItem(MapMetricItem(metrics.OutPullOutputByID, OutputIDFromIotaGoOutputID(metrics.OutPullOutputByID.LastMessage))),
	}
}

func MapNodeMessageMetrics(metrics *dto.NodeMessageMetrics) *NodeMessageMetrics {
	return &NodeMessageMetrics{
		RegisteredChainIDs: MapRegisteredChainIDs(metrics.RegisteredChainIDs),

		InMilestone:        MilestoneMetricItem(MapMetricItem(metrics.InMilestone, MilestoneFromIotaGoMilestone(metrics.InMilestone.LastMessage))),
		InStateOutput:      InStateOutputMetricItem(MapMetricItem(metrics.InStateOutput, InStateOutputFromISCInStateOutput(metrics.InStateOutput.LastMessage))),
		InAliasOutput:      AliasOutputMetricItem(MapMetricItem(metrics.InAliasOutput, OutputFromIotaGoOutput(metrics.InAliasOutput.LastMessage))),
		InOutput:           InOutputMetricItem(MapMetricItem(metrics.InOutput, InOutputFromISCInOutput(metrics.InOutput.LastMessage))),
		InOnLedgerRequest:  OnLedgerRequestMetricItem(MapMetricItem(metrics.InOnLedgerRequest, OnLedgerRequestFromISC(metrics.InOnLedgerRequest.LastMessage))),
		InTxInclusionState: TxInclusionStateMsgMetricItem(MapMetricItem(metrics.InTxInclusionState, TxInclusionStateMsgFromISCTxInclusionStateMsg(metrics.InTxInclusionState.LastMessage))),

		OutPublishStateTransaction:      PublisherStateTransactionItem(MapMetricItem(metrics.OutPublishStateTransaction, StateTransactionFromISCStateTransaction(metrics.OutPublishStateTransaction.LastMessage))),
		OutPublishGovernanceTransaction: TransactionMetricItem(MapMetricItem(metrics.OutPublishGovernanceTransaction, TransactionFromIotaGoTransaction(metrics.OutPublishGovernanceTransaction.LastMessage))),
		OutPullLatestOutput:             InterfaceMetricItem(MapMetricItem(metrics.OutPullLatestOutput, metrics.OutPullLatestOutput.LastMessage)),
		OutPullTxInclusionState:         TransactionIDMetricItem(MapMetricItem(metrics.OutPullTxInclusionState, TransactionFromIotaGoTransactionID(&metrics.OutPullTxInclusionState.LastMessage))),
		OutPullOutputByID:               UTXOInputMetricItem(MapMetricItem(metrics.OutPullOutputByID, OutputIDFromIotaGoOutputID(metrics.OutPullOutputByID.LastMessage))),
	}
}

type ConsensusWorkflowMetrics struct {
	FlagStateReceived        bool `json:"flagStateReceived" swagger:"desc(Shows if state output is received in current consensus iteration),required"`
	FlagBatchProposalSent    bool `json:"flagBatchProposalSent" swagger:"desc(Shows if batch proposal is sent out in current consensus iteration),required"`
	FlagConsensusBatchKnown  bool `json:"flagConsensusBatchKnown" swagger:"desc(Shows if consensus on batch is reached and known in current consensus iteration),required"`
	FlagVMStarted            bool `json:"flagVMStarted" swagger:"desc(Shows if virtual machine is started in current consensus iteration),required"`
	FlagVMResultSigned       bool `json:"flagVMResultSigned" swagger:"desc(Shows if virtual machine has returned its results in current consensus iteration),required"`
	FlagTransactionFinalized bool `json:"flagTransactionFinalized" swagger:"desc(Shows if consensus on transaction is reached in current consensus iteration),required"`
	FlagTransactionPosted    bool `json:"flagTransactionPosted" swagger:"desc(Shows if transaction is posted to L1 in current consensus iteration),required"`
	FlagTransactionSeen      bool `json:"flagTransactionSeen" swagger:"desc(Shows if L1 reported that it has seen the transaction of current consensus iteration),required"`
	FlagInProgress           bool `json:"flagInProgress" swagger:"desc(Shows if consensus algorithm is still not completed in current consensus iteration),required"`

	TimeBatchProposalSent    time.Time `json:"timeBatchProposalSent" swagger:"desc(Shows when batch proposal was last sent out in current consensus iteration),required"`
	TimeConsensusBatchKnown  time.Time `json:"timeConsensusBatchKnown" swagger:"desc(Shows when ACS results of consensus on batch was last received in current consensus iteration),required"`
	TimeVMStarted            time.Time `json:"timeVMStarted" swagger:"desc(Shows when virtual machine was last started in current consensus iteration),required"`
	TimeVMResultSigned       time.Time `json:"timeVMResultSigned" swagger:"desc(Shows when virtual machine results were last received and signed in current consensus iteration),required"`
	TimeTransactionFinalized time.Time `json:"timeTransactionFinalized" swagger:"desc(Shows when algorithm last noted that all the data for consensus on transaction had been received in current consensus iteration),required"`
	TimeTransactionPosted    time.Time `json:"timeTransactionPosted" swagger:"desc(Shows when transaction was last posted to L1 in current consensus iteration),required"`
	TimeTransactionSeen      time.Time `json:"timeTransactionSeen" swagger:"desc(Shows when algorithm last noted that transaction had been seen by L1 in current consensus iteration),required"`
	TimeCompleted            time.Time `json:"timeCompleted" swagger:"desc(Shows when algorithm was last completed in current consensus iteration),required"`

	CurrentStateIndex uint32 `json:"currentStateIndex" swagger:"desc(Shows current state index of the consensus),min(1)"`
}

func MapConsensusWorkflowStatus(status chain.ConsensusWorkflowStatus) *ConsensusWorkflowMetrics {
	return &ConsensusWorkflowMetrics{
		FlagStateReceived:        status.IsStateReceived(),
		FlagBatchProposalSent:    status.IsBatchProposalSent(),
		FlagConsensusBatchKnown:  status.IsConsensusBatchKnown(),
		FlagVMStarted:            status.IsVMStarted(),
		FlagVMResultSigned:       status.IsVMResultSigned(),
		FlagTransactionFinalized: status.IsTransactionFinalized(),
		FlagTransactionPosted:    status.IsTransactionPosted(),
		FlagTransactionSeen:      status.IsTransactionSeen(),
		FlagInProgress:           status.IsInProgress(),

		TimeBatchProposalSent:    status.GetBatchProposalSentTime(),
		TimeConsensusBatchKnown:  status.GetConsensusBatchKnownTime(),
		TimeVMStarted:            status.GetVMStartedTime(),
		TimeVMResultSigned:       status.GetVMResultSignedTime(),
		TimeTransactionFinalized: status.GetTransactionFinalizedTime(),
		TimeTransactionPosted:    status.GetTransactionPostedTime(),
		TimeTransactionSeen:      status.GetTransactionSeenTime(),
		TimeCompleted:            status.GetCompletedTime(),

		CurrentStateIndex: status.GetCurrentStateIndex(),
	}
}

type ConsensusPipeMetrics struct {
	EventStateTransitionMsgPipeSize int `json:"eventStateTransitionMsgPipeSize" swagger:"required"`
	EventPeerLogIndexMsgPipeSize    int `json:"eventPeerLogIndexMsgPipeSize" swagger:"required"`
	EventACSMsgPipeSize             int `json:"eventACSMsgPipeSize" swagger:"required"`
	EventVMResultMsgPipeSize        int `json:"eventVMResultMsgPipeSize" swagger:"required"`
	EventTimerMsgPipeSize           int `json:"eventTimerMsgPipeSize" swagger:"required"`
}

func MapConsensusPipeMetrics(pipeMetrics chain.ConsensusPipeMetrics) *ConsensusPipeMetrics {
	return &ConsensusPipeMetrics{
		EventStateTransitionMsgPipeSize: pipeMetrics.GetEventStateTransitionMsgPipeSize(),
		EventPeerLogIndexMsgPipeSize:    pipeMetrics.GetEventPeerLogIndexMsgPipeSize(),
		EventACSMsgPipeSize:             pipeMetrics.GetEventACSMsgPipeSize(),
		EventVMResultMsgPipeSize:        pipeMetrics.GetEventVMResultMsgPipeSize(),
		EventTimerMsgPipeSize:           pipeMetrics.GetEventTimerMsgPipeSize(),
	}
}
