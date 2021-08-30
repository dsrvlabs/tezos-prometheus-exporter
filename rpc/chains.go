package rpc

import (
	"time"
)

// BootstrapStatus
const (
	ChainStatusStuck    ChainStatus = "stuck"
	ChainStatusSynced   ChainStatus = "synced"
	ChainStatusUnsynced ChainStatus = "unsynced"
)

// ChainStatus describes status of network connection of peer.
type ChainStatus string

// BootstrapStatus describes status of node.
type BootstrapStatus struct {
	IsBootstrapped bool        `json:"bootstrapped"`
	SyncState      ChainStatus `json:"sync_state"`
}

// Block descibes single block.
type Block struct {
	Protocol string       `json:"protocol"`
	ChainID  string       `json:"chain_id"`
	Hash     string       `json:"hash"`
	Header   *BlockHeader `json:"header"`
}

// BlockHeader decrives block header.
type BlockHeader struct {
	Level          uint      `json:"level"`
	Proto          uint      `json:"proto"`
	Predecessor    string    `json:"predecessor"`
	Timestamp      time.Time `json:"timestamp"`
	ValidationPass uint      `json:"validation_pass"`
	OperationsHash string    `json:"operations_hash"`
}
