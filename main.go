package main

import (
	"log"
	"net/http"

	"dsrvlabs/tezos-prometheus-exporter/exporter"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricExporter exporter.Exporter
)

func init() {
	metricExporter = exporter.NewExporter()
	err := metricExporter.Collect()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	// TODO: Get RPC address.

	e := echo.New()

	e.GET("/health", health)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// TODO: port should be configurable.
	if err := e.Start(":9489"); err != nil {
		panic(err)
	}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
