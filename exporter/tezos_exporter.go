package exporter

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"dsrvlabs/tezos-prometheus-exporter/rpc"
)

func init() {
}

type tezosExporter struct {
	rpcClient  rpc.Client
	blockLevel prometheus.Gauge
	peerCount  prometheus.Gauge
}

func (e *tezosExporter) Collect() error {
	log.Println("Collect Tezos metric")

	go func() {
		err := e.getInfo()
		_ = err
		time.Sleep(fetchInterval)
	}()

	return nil
}

func (e *tezosExporter) getInfo() error {
	rpc := e.rpcClient

	bootstrapStatus, err := rpc.GetBootstrapStatus()
	if err != nil {
		log.Println(err)
		return err
	}

	headBlock, err := rpc.GetHeadBlock()
	if err != nil {
		log.Println(err)
		return err
	}

	peers, err := rpc.GetPeers()
	if err != nil {
		log.Println(err)
		return err
	}

	// TODO: Filter running.

	_ = bootstrapStatus

	// TODO: Write bootstrap status
	e.blockLevel.Set(float64(headBlock.Header.Level))
	e.peerCount.Set(float64(len(peers)))

	return err
}

func (e *tezosExporter) Stop() error {
	// TODO:
	return nil
}

func createTezosExporter() Exporter {

	blockLevel := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "block_level",
		Help: "Current block level of Tezos network",
	})

	peerCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "peer_count",
		Help: "Count of connected peers on the node.",
	})

	prometheus.MustRegister(blockLevel)
	prometheus.MustRegister(peerCount)

	// TODO: Hosts?? Fetch from config file?
	return &tezosExporter{
		rpcClient:  rpc.NewClient("http://localhost:8732"),
		blockLevel: blockLevel,
		peerCount:  peerCount,
	}
}
