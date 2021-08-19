package exporter

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dsrvlabs/tezos-prometheus-exporter/mocks"
	"dsrvlabs/tezos-prometheus-exporter/rpc"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestTezosMetric(t *testing.T) {
	exporter := createTezosExporter("http://localhost:8732", 1).(*tezosExporter)

	tests := []struct {
		mockBootstrapStatus rpc.BootstrapStatus
		mockBlock           rpc.Block
		mockPeers           []rpc.Peer
	}{
		{
			mockBootstrapStatus: rpc.BootstrapStatus{
				IsBootstrapped: true,
				SyncState:      rpc.ChainStatusSynced,
			},

			mockBlock: rpc.Block{
				Header: &rpc.BlockHeader{
					Level: 100,
				},
			},
			mockPeers: []rpc.Peer{
				{ID: "peer-id-0", State: rpc.PeerStateRunning},
				{ID: "peer-id-1", State: rpc.PeerStateAccepted},
				{ID: "peer-id-2", State: rpc.PeerStateRunning},
				{ID: "peer-id-3", State: rpc.PeerStateDisconnected},
				{ID: "peer-id-4", State: rpc.PeerStateRunning},
			},
		},

		{
			mockBootstrapStatus: rpc.BootstrapStatus{
				IsBootstrapped: false,
				SyncState:      rpc.ChainStatusStuck,
			},

			mockBlock: rpc.Block{
				Header: &rpc.BlockHeader{
					Level: 123456,
				},
			},
			mockPeers: []rpc.Peer{
				{ID: "peer-id-10", State: rpc.PeerStateAccepted},
				{ID: "peer-id-11", State: rpc.PeerStateAccepted},
				{ID: "peer-id-22", State: rpc.PeerStateAccepted},
				{ID: "peer-id-33", State: rpc.PeerStateAccepted},
				{ID: "peer-id-44", State: rpc.PeerStateAccepted},
			},
		},
	}

	for _, test := range tests {
		// Mocks
		mockClient := mocks.Client{}
		exporter.rpcClient = &mockClient

		mockClient.On("GetBootstrapStatus").Return(&test.mockBootstrapStatus, nil)
		mockClient.On("GetHeadBlock").Return(&test.mockBlock, nil)
		mockClient.On("GetPeers").Return(test.mockPeers, nil)

		err := exporter.getInfo()

		// Test
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		rr := httptest.NewRecorder()
		c := e.NewContext(req, rr)

		h := echo.WrapHandler(promhttp.Handler())
		err = h(c)

		// Asserts
		runningMockPeers := make([]rpc.Peer, 0)
		for _, p := range test.mockPeers {
			if p.State == rpc.PeerStateRunning {
				runningMockPeers = append(runningMockPeers, p)
			}
		}

		rawBody, err := io.ReadAll(rr.Body)

		assert.Nil(t, err)

		value, err := findMetric("tezos_block_level", strings.NewReader(string(rawBody)))
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%d", test.mockBlock.Header.Level), value)

		value, err = findMetric("tezos_peer_count", strings.NewReader(string(rawBody)))
		assert.Nil(t, err)
		assert.Equal(t, fmt.Sprintf("%d", len(runningMockPeers)), value)

		value, err = findMetric("tezos_sync_status", strings.NewReader(string(rawBody)))
		assert.Nil(t, err)

		expectValue := bootstrapMap[test.mockBootstrapStatus.IsBootstrapped]
		assert.Equal(t, fmt.Sprintf("%d", int(expectValue)), value)

		value, err = findMetric("tezos_is_bootstrapped", strings.NewReader(string(rawBody)))
		assert.Nil(t, err)

		expectValue = syncMap[test.mockBootstrapStatus.SyncState]
		assert.Equal(t, fmt.Sprintf("%d", int(expectValue)), value)
	}
}
