package exporter

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
)

var (
	errMetricNotFound = errors.New("Specified metric cannot be found")
)

func TestCPU(t *testing.T) {
	testExporter := createSystemExporter().(*systemExporter)

	assert.NotNil(t, testExporter)

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

	cpu, err := findMetric("cpu_usage", rr.Body)
	memory, err := findMetric("memory_usage", rr.Body)

	assert.Nil(t, err)
	assert.Equal(t, "10.1234", cpu)
	assert.Equal(t, "23.12345", memory)
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
