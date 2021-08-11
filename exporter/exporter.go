package exporter

import (
	"log"

	cfg "dsrvlabs/tezos-prometheus-exporter/config"
)

// Exporter provides exporting features.
type Exporter interface {
	Collect() error
	Stop() error
}

type nodeExporter struct {
	systemExporter Exporter
	nodeExporter   Exporter
}

func (e *nodeExporter) Collect() error {
	err := e.systemExporter.Collect()
	if err != nil {
		log.Println(err)
		return err
	}

	err = e.nodeExporter.Collect()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (e *nodeExporter) Stop() error {
	// TODO:
	return nil
}

// NewExporter create new exporter instances.
func NewExporter(config cfg.Config) Exporter {
	// TODO: exporters should be configurable.
	return &nodeExporter{
		systemExporter: createSystemExporter(config.MountPath, config.UpdateIntervalSeconds),
		nodeExporter:   createTezosExporter(config.RPCEndpoint, config.UpdateIntervalSeconds),
	}
}
