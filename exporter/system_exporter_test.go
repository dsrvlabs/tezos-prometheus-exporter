package exporter

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"

	cfg "dsrvlabs/tezos-prometheus-exporter/config"
)

var (
	errMetricNotFound = errors.New("Specified metric cannot be found")
	testExporter      *systemExporter
)

func init() {
	// config := cfg.GetConfig()
	loader := cfg.NewLoader("")
	config, _ := loader.Load()

	testExporter = createSystemExporter(config.DataDir, config.UpdateIntervalSeconds).(*systemExporter)
}

func TestCPU(t *testing.T) {
	// Set dummy
	testExporter.cpuUsage.Set(10.1234)
	testExporter.memoryUsage.Set(23.12345)

	// Test
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rr)

	h := echo.WrapHandler(promhttp.Handler())
	err := h(c)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)

	rawBody, err := io.ReadAll(rr.Body)
	assert.Nil(t, err)

	cpu, err := findMetric("system_cpu_usage", strings.NewReader(string(rawBody)))
	memory, err := findMetric("system_memory_usage", strings.NewReader(string(rawBody)))

	assert.Nil(t, err)
	assert.Equal(t, "10.1234", cpu)
	assert.Equal(t, "23.12345", memory)
}

func TestGetSystemMetric(t *testing.T) {
	load, idle, err := testExporter.getCPUUsage(1 * time.Second)

	assert.Nil(t, err)

	// Intentionally increase allowed value.
	// For more detail resolution, cached value should be included.
	assert.InDelta(t, 100, load+idle, 5.0)

	used, free, err := testExporter.getMemUsage()

	assert.Nil(t, err)

	// Intentionally increase allowed value.
	// For more detail resolution, cached value should be included.
	assert.InDelta(t, 100, used+free, 5.0)
}

func TestDiskMetric(t *testing.T) {
	free, total, err := testExporter.getDiskUsage("/")

	assert.Nil(t, err)
	assert.True(t, free < total)
}

func findMetric(name string, content io.Reader) (string, error) {
	contentReader := bufio.NewReader(content)

	for {
		line, _, err := contentReader.ReadLine()
		if err != nil {
			break
		}

		foundIndex := strings.Index(string(line), name)
		if foundIndex == 0 {
			tokens := strings.Split(string(line), " ")
			return tokens[1], nil
		}
	}

	return "", errMetricNotFound
}
