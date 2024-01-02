package models

// DKSharesPostRequest is a POST request for creating new DKShare.
type DKSharesPostRequest struct {
	PeerPubKeysOrNames []string `json:"peerIdentities" swagger:"desc(Names or hex encoded public keys of trusted peers to run DKG on.),required"`
	Threshold          uint16   `json:"threshold" swagger:"desc(Should be =< len(PeerPublicIdentities)),required,min(1)"`
	TimeoutMS          uint32   `json:"timeoutMS" swagger:"desc(Timeout in milliseconds.),required,min(1)"`
}

// DKSharesInfo stands for the DKShare representation, returned by the GET and POST methods.
type DKSharesInfo struct {
	Address         string   `json:"address" swagger:"desc(New generated shared address.),required"`
	PeerIdentities  []string `json:"peerIdentities" swagger:"desc(Identities of the nodes sharing the key. (Hex)),required"`
	PeerIndex       *uint16  `json:"peerIndex" swagger:"desc(Index of the node returning the share, if it is a member of the sharing group.),required,min(1)"`
	PublicKey       string   `json:"publicKey" swagger:"desc(Used public key. (Hex)),required"`
	PublicKeyShares []string `json:"publicKeyShares" swagger:"desc(Public key shares for all the peers. (Hex)),required"`
	Threshold       uint16   `json:"threshold" swagger:"required,min(1)"`
}
