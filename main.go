package main

import (
	"fmt"
	"log"
	"net/http"

	cfg "dsrvlabs/tezos-prometheus-exporter/config"
	"dsrvlabs/tezos-prometheus-exporter/exporter"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricExporter exporter.Exporter
)

func init() {
	config := cfg.GetConfig()

	metricExporter = exporter.NewExporter(config)
	err := metricExporter.Collect()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	e := echo.New()

	e.GET("/health", health)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	config := cfg.GetConfig()

	if err := e.Start(fmt.Sprintf(":%d", config.ServicePort)); err != nil {
		panic(err)
	}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
