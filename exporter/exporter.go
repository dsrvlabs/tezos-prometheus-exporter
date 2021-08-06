package exporter

import "log"

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
func NewExporter() Exporter {
	// TODO: exporters should be configurable.
	return &nodeExporter{
		systemExporter: createSystemExporter(),
		nodeExporter:   createTezosExporter(),
	}
}
