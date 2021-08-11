package exporter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"dsrvlabs/tezos-prometheus-exporter/mocks"
	"dsrvlabs/tezos-prometheus-exporter/rpc"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

func TestTezosMetric(t *testing.T) {
	exporter := createTezosExporter("http://localhost:8732", 1).(*tezosExporter)

	// Mocks
	mockClient := mocks.Client{}
	exporter.rpcClient = &mockClient

	mockBootstrapStatus := rpc.BootstrapStatus{
		IsBootstrapped: true,
		SyncState:      rpc.ChainStatusSynced,
	}

	mockBlock := rpc.Block{
		Header: &rpc.BlockHeader{
			Level: 100,
		},
	}
	mockPeers := []rpc.Peer{
		{ID: "peer-id-0", State: rpc.PeerStateRunning},
		{ID: "peer-id-1", State: rpc.PeerStateAccepted},
	}

	mockClient.On("GetBootstrapStatus").Return(&mockBootstrapStatus, nil)
	mockClient.On("GetHeadBlock").Return(&mockBlock, nil)
	mockClient.On("GetPeers").Return(mockPeers, nil)

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
	for _, p := range mockPeers {
		if p.State == rpc.PeerStateRunning {
			runningMockPeers = append(runningMockPeers, p)
		}
	}

	assert.Nil(t, err)
	value, err := findMetric("block_level", rr.Body)
	assert.Equal(t, "100", value)

	value, err = findMetric("peer_count", rr.Body)
	assert.Equal(t, fmt.Sprintf("%d", len(runningMockPeers)), value)
}
