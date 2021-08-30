package exporter

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	"dsrvlabs/tezos-prometheus-exporter/mocks"
)

func TestDaemonList(t *testing.T) {
	processNames := []string{
		"tezos-node",
		"tezos-baker",
		"tezos-endorser",
		"tezos-accuser",
	}
	exporter := createDaemonExporter(processNames, 1).(*daemonExporter)

	mockProcess := mocks.Process{}

	exporter.manager = &mockProcess

	// Mocks

	mockProcess.On("IsRunning", "tezos-node").Return(true, nil)
	mockProcess.On("IsRunning", "tezos-baker").Return(false, nil)
	mockProcess.On("IsRunning", "tezos-endorser").Return(true, nil)
	mockProcess.On("IsRunning", "tezos-accuser").Return(false, nil)

	// Test
	err := exporter.checkDaemons()

	// Asserts
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rr)

	h := echo.WrapHandler(promhttp.Handler())
	err = h(c)

	assert.Nil(t, err)

	rawBody, err := io.ReadAll(rr.Body)

	assert.Nil(t, err)

	value, err := findMetric(fmt.Sprintf("tezos_daemon{name=\"tezos-node\"}"), strings.NewReader(string(rawBody)))
	assert.Nil(t, err)
	assert.Equal(t, "1", value)

	value, err = findMetric(fmt.Sprintf("tezos_daemon{name=\"tezos-baker\"}"), strings.NewReader(string(rawBody)))
	assert.Nil(t, err)
	assert.Equal(t, "0", value)

	value, err = findMetric(fmt.Sprintf("tezos_daemon{name=\"tezos-endorser\"}"), strings.NewReader(string(rawBody)))
	assert.Nil(t, err)
	assert.Equal(t, "1", value)

	value, err = findMetric(fmt.Sprintf("tezos_daemon{name=\"tezos-accuser\"}"), strings.NewReader(string(rawBody)))
	assert.Nil(t, err)
	assert.Equal(t, "0", value)
}
