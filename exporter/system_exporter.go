package exporter

import (
	"fmt"
	"log"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	// TODO: From config?
	fetchInterval = time.Second * 1
)

func init() {
}

type systemExporter struct {
	cpuUsage      prometheus.Gauge
	memoryUsage   prometheus.Gauge
	diskFreeSpace prometheus.Gauge
}

func (e *systemExporter) Collect() error {
	log.Println("Collect system metric")

	go func() {
		for {
			load, _, err := e.getCPUUsage(3 * time.Second)
			if err != nil {
				log.Println(err)
				continue
			}

			e.cpuUsage.Set(load)

			used, _, err := e.getMemUsage()
			if err != nil {
				log.Println(err)
				continue
			}

			e.memoryUsage.Set(used)

			// TODO: Path should be configurable
			free, _, err := e.getDiskUsage("/mnt/tezos")
			if err != nil {
				log.Println(err)
				continue
			}

			e.diskFreeSpace.Set(float64(free))

			fmt.Printf("CPU: %+v, MEM: %+v Disk Free: %+v\n", load, used, free)

			time.Sleep(fetchInterval)
		}
	}()

	return nil
}

// getCPUUsage function returns load and idle percentages of CPU time.
// This function get interval parameter to make delay intentionally for measuring CPU counter.
func (e *systemExporter) getCPUUsage(measureInterval time.Duration) (float64, float64, error) {
	cpuBefore, err := cpu.Get()
	if err != nil {
		return 0, 0, err
	}

	time.Sleep(measureInterval)

	cpuAfter, err := cpu.Get()
	if err != nil {
		return 0, 0, err
	}

	userTotal := cpuAfter.User - cpuBefore.User
	systemTotal := cpuAfter.System - cpuBefore.System
	idleTotal := cpuAfter.Idle - cpuBefore.Idle

	total := cpuAfter.Total - cpuBefore.Total

	cpuLoad := float64(userTotal+systemTotal) / float64(total) * 100
	idleLoad := float64(idleTotal) / float64(total) * 100

	return cpuLoad, idleLoad, nil
}

// getMemUsage function returns use / free memory space.
func (e *systemExporter) getMemUsage() (float64, float64, error) {
	memStat, err := memory.Get()
	if err != nil {
		return 0, 0, err
	}

	usingSpace := float64(memStat.Used) / float64(memStat.Total) * 100
	freeSpace := float64(memStat.Cached+memStat.Free) / float64(memStat.Total) * 100

	return usingSpace, freeSpace, nil
}

// getDiskUsage function returns free MB / tatal MB of disk space.
func (e *systemExporter) getDiskUsage(path string) (free, total uint64, err error) {
	var stat unix.Statfs_t

	err = unix.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}

	const MB = 1024 * 1024

	free = uint64(stat.Bsize) * stat.Bfree / MB
	total = uint64(stat.Bsize) * stat.Blocks / MB

	return free, total, nil
}

func (e *systemExporter) Stop() error {
	// TODO:
	return nil
}

func createSystemExporter() Exporter {
	log.Println("createSystemExporter")

	cpuUsage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "System CPU Usage in percentage",
	})

	memoryUsage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "memory_usage",
		Help: "System Memory Usage in percentage",
	})

	diskFreeSpace := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_free",
		Help: "Free disk space of selected path in megabyte",
	})

	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsage)
	prometheus.MustRegister(diskFreeSpace)

	return &systemExporter{
		cpuUsage:      cpuUsage,
		memoryUsage:   memoryUsage,
		diskFreeSpace: diskFreeSpace,
	}
}
