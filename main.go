package main

import (
	"net/http"

	"dsrvlabs/tezos-prometheus-exporter/exporter"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	tezosExporter exporter.Exporter
)

func init() {
	tezosExporter = exporter.NewExporter()
	tezosExporter.Collect()
}

func main() {
	e := echo.New()

	e.GET("/health", health)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
