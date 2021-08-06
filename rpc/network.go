package rpc

// Peer states
var (
	PeerStateRunning      PeerState = "running"
	PeerStateAccepted     PeerState = "accepted"
	PeerStateDisconnected PeerState = "disconnected"
)

// PeerState describes connection state of peer.
type PeerState string

// Peer describes single peer.
type Peer struct {
	ID      string    `json:"id"`
	Score   uint      `json:"score"`
	Trusted bool      `json:"trusted"`
	State   PeerState `json:"state"`
}
