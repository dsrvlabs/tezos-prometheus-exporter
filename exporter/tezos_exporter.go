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
	rpcClient     rpc.Client
	blockLevel    prometheus.Gauge
	peerCount     prometheus.Gauge
	fetchInterval time.Duration
}

func (e *tezosExporter) Collect() error {
	log.Println("Collect Tezos metric")

	go func() {
		err := e.getInfo()
		_ = err
		time.Sleep(e.fetchInterval)
	}()

	return nil
}

func (e *tezosExporter) getInfo() error {
	cli := e.rpcClient

	bootstrapStatus, err := cli.GetBootstrapStatus()
	if err != nil {
		log.Println(err)
		return err
	}

	_ = bootstrapStatus

	headBlock, err := cli.GetHeadBlock()
	if err != nil {
		log.Println(err)
		return err
	}

	peers, err := cli.GetPeers()
	if err != nil {
		log.Println(err)
		return err
	}

	runningPeers := make([]rpc.Peer, 0)
	for _, peer := range peers {
		if peer.State == rpc.PeerStateRunning {
			runningPeers = append(runningPeers, peer)
		}
	}

	e.blockLevel.Set(float64(headBlock.Header.Level))
	e.peerCount.Set(float64(len(runningPeers)))

	return err
}

func (e *tezosExporter) Stop() error {
	// TODO:
	return nil
}

func createTezosExporter(rpcEndpoint string, fetchInterval int) Exporter {
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

	return &tezosExporter{
		rpcClient:     rpc.NewClient(rpcEndpoint),
		blockLevel:    blockLevel,
		peerCount:     peerCount,
		fetchInterval: time.Duration(fetchInterval) * time.Second,
	}
}
