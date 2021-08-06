package exporter

// Exporter provides exporting features.
type Exporter interface {
	Collect() error
	Stop() error
}

// NewExporter create new exporter instances.
func NewExporter() Exporter {
	return createSystemExporter()
}
