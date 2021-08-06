package exporter

import (
	"fmt"
	"log"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// TODO: From config?
	fetchInterval = time.Second * 1
)

func init() {
}

type systemExporter struct {
	cpuUsage    prometheus.Gauge
	memoryUsage prometheus.Gauge
	// diskUsage   prometheus.Gauge
}

func (e *systemExporter) Collect() error {
	log.Println("Starting exporter")

	go func() {
		for {
			cpuStat, err := cpu.Get()
			if err != nil {
				log.Println("Err get CPU Stat")
				continue
			}

			cpuLoad := float64(cpuStat.System+cpuStat.User) / float64(cpuStat.Total) * 100
			e.cpuUsage.Set(cpuLoad)

			memStat, err := memory.Get()
			if err != nil {
				log.Println("Err get MEM Stat")
				continue
			}

			memUsage := float64(memStat.Used) / float64(memStat.Total)
			e.memoryUsage.Set(memUsage)

			fmt.Printf("CPU: %+v, MEM: %+v\n", cpuLoad, memUsage)

			time.Sleep(fetchInterval)
		}
	}()

	return nil
}

func (e *systemExporter) Stop() error {
	return nil
}

func createSystemExporter() *systemExporter {
	log.Println("createSystemExporter")

	cpuUsage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "System CPU Usage in percentage",
	})

	memoryUsage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "System Memory Usage in percentage",
	})

	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)

	return &systemExporter{
		cpuUsage:    cpuUsage,
		memoryUsage: memoryUsage,
	}
}
