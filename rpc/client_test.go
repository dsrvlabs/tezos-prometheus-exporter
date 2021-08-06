package rpc

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRPC_GetHeadBlock(t *testing.T) {
	tests := []struct {
		ExpectURL   string
		ExpectBlock *Block
	}{
		{
			ExpectURL: "http://localhost:8732/chains/main/blocks/head",
			ExpectBlock: &Block{
				Protocol: "PtGRANADsDU8R9daYKAgWnQYAJ64omN1o3KMGVCykShA97vQbvV",
				ChainID:  "NetXz969SFaFn8k",
				Hash:     "BMVZwd6CDft9DSCpRoYSC422J8f9h8HRXY5DwqQDotVbp4PgXVT",
				Header: &BlockHeader{
					Level:          289076,
					Proto:          2,
					Predecessor:    "BLuhg29MDsYv7FerMVwKXb77gRwzfikPRZXk454MniBP4fV9tTd",
					Timestamp:      time.Date(2021, time.August, 6, 9, 53, 56, 0, time.UTC),
					ValidationPass: 4,
					OperationsHash: "LLoZxhd8gn6Y2usCRkuj6Hpb7KUgTs9cTSUUf3MRTBZZh26L4yqB3",
				},
			},
		},
	}

	for _, test := range tests {
		// Mocks
		f, _ := os.Open("fixtures/chains_main_blocks_head.json")
		rawBody, _ := io.ReadAll(f)
		mockResponder := httpmock.NewStringResponder(http.StatusOK, string(rawBody))
		httpmock.RegisterResponder(http.MethodGet, test.ExpectURL, mockResponder)
		httpmock.Activate()

		// Test
		cli := NewClient("http://localhost:8732")
		respBlock, err := cli.GetHeadBlock()

		// Asserts
		assert.Nil(t, err)
		assert.Equal(t, test.ExpectBlock.Protocol, respBlock.Protocol)
		assert.Equal(t, test.ExpectBlock.ChainID, respBlock.ChainID)
		assert.Equal(t, test.ExpectBlock.Hash, respBlock.Hash)
		assert.Equal(t, test.ExpectBlock.Header.Level, respBlock.Header.Level)
		assert.Equal(t, test.ExpectBlock.Header.Proto, respBlock.Header.Proto)
		assert.Equal(t, test.ExpectBlock.Header.Predecessor, respBlock.Header.Predecessor)
		assert.Equal(t, test.ExpectBlock.Header.Timestamp.Unix(), respBlock.Header.Timestamp.Unix())
		assert.Equal(t, test.ExpectBlock.Header.ValidationPass, respBlock.Header.ValidationPass)
		assert.Equal(t, test.ExpectBlock.Header.OperationsHash, respBlock.Header.OperationsHash)
	}
}

func TestRPC_GetPeers(t *testing.T) {
	expectURL := "http://localhost:8732/network/peers"
	expectPeers := []Peer{
		{
			ID:      "idqkcJVTXNzpyHegWqwSGj7s4K4Qmt",
			Score:   1,
			Trusted: false,
			State:   PeerStateDisconnected,
		},
		{
			ID:      "idsX66KveZbkqLcgpKN8L3oBypuuSu",
			Score:   1,
			Trusted: false,
			State:   PeerStateDisconnected,
		},
		{
			ID:      "idtCyCFkvG2quNh84m67dzwonSg8qb",
			Score:   1,
			Trusted: false,
			State:   PeerStateDisconnected,
		},
		{
			ID:      "idqSFdRtFXL7Lo25QdR42gXquV4KGo",
			Score:   1,
			Trusted: false,
			State:   PeerStateRunning,
		},
	}

	// Mocks
	f, _ := os.Open("fixtures/network_peers.json")
	rawBody, err := io.ReadAll(f)
	mockResponder := httpmock.NewStringResponder(http.StatusOK, string(rawBody))
	httpmock.RegisterResponder(http.MethodGet, expectURL, mockResponder)
	httpmock.Activate()

	cli := NewClient("http://localhost:8732")
	peers, err := cli.GetPeers()

	assert.Nil(t, err)

	// Asserts
	assert.Nil(t, err)
	assert.Equal(t, 4, len(peers))

	for i, expect := range expectPeers {
		assert.Equal(t, expect.ID, peers[i].ID)
		assert.Equal(t, expect.Score, peers[i].Score)
		assert.Equal(t, expect.Trusted, peers[i].Trusted)
		assert.Equal(t, expect.State, peers[i].State)
	}
}

func TestRPC_IsBootstrapped(t *testing.T) {
	expectURL := "http://localhost:8732/chains/main/is_bootstrapped"
	expectStatus := BootstrapStatus{
		IsBootstrapped: true,
		SyncState:      ChainStatusSynced,
	}

	f, _ := os.Open("fixtures/chains_main_is_bootstrapped.json")
	rawBody, err := io.ReadAll(f)
	mockResponder := httpmock.NewStringResponder(http.StatusOK, string(rawBody))
	httpmock.RegisterResponder(http.MethodGet, expectURL, mockResponder)
	httpmock.Activate()

	cli := NewClient("http://localhost:8732")
	status, err := cli.GetBootstrapStatus()

	assert.Nil(t, err)
	assert.Equal(t, expectStatus.IsBootstrapped, status.IsBootstrapped)
	assert.Equal(t, expectStatus.SyncState, status.SyncState)
}
