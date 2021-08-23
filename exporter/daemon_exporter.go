package exporter

import (
	"log"
	"time"

	ps "dsrvlabs/tezos-prometheus-exporter/process"

	"github.com/prometheus/client_golang/prometheus"
)

type daemonExporter struct {
	gaugeVec      *prometheus.GaugeVec
	daemonNames   []string
	manager       ps.Process
	fetchInterval time.Duration
}

func (e *daemonExporter) Collect() error {
	go func() {
		for {
			log.Println("Check daemons")
			err := e.checkDaemons()
			if err != nil {
				log.Println(err)
				time.Sleep(e.fetchInterval)
				continue
			}

			time.Sleep(e.fetchInterval)
		}
	}()

	return nil
}

func (e *daemonExporter) checkDaemons() error {
	for _, name := range e.daemonNames {
		running, err := e.manager.IsRunning(name)
		if err != nil {
			return err
		}

		log.Println("check daemon", name, running)

		if running {
			e.gaugeVec.WithLabelValues(name).Set(1)
		} else {
			e.gaugeVec.WithLabelValues(name).Set(0)
		}
	}

	return nil
}

func (e *daemonExporter) Stop() error {
	return nil
}

func createDaemonExporter(names []string, fetchInterval int) Exporter {
	// TODO: Refactoring.

	gaugeVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "tezos",
		Name:      "daemon",
	}, []string{"name"})

	prometheus.MustRegister(gaugeVec)

	return &daemonExporter{
		gaugeVec:      gaugeVec,
		daemonNames:   names,
		manager:       ps.NewProcessManager(),
		fetchInterval: time.Duration(fetchInterval) * time.Second,
	}
}
