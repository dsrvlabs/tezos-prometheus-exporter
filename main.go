package main

import (
	"flag"
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
	config         cfg.Config
)

func init() {
	var configFilename string

	flag.StringVar(&configFilename, "config", "config.json", "-config=config.json")
	flag.Parse()

	var err error

	cfgLoader := cfg.NewLoader(configFilename)
	config, err = cfgLoader.Load()
	if err != nil {
		log.Panic(err)
	}

	metricExporter = exporter.NewExporter(config)
	err = metricExporter.Collect()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	e := echo.New()

	e.GET("/health", health)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	if err := e.Start(fmt.Sprintf(":%d", config.ExporterPort)); err != nil {
		panic(err)
	}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
