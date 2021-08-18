package exporter

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"dsrvlabs/tezos-prometheus-exporter/rpc"
)

var (
	bootstrapMap = map[bool]float64{
		false: 0,
		true:  1,
	}

	syncMap = map[rpc.ChainStatus]float64{
		rpc.ChainStatusStuck:    0,
		rpc.ChainStatusSynced:   1,
		rpc.ChainStatusUnsynced: 2,
	}
)

type tezosExporter struct {
	rpcClient     rpc.Client
	blockLevel    prometheus.Gauge
	peerCount     prometheus.Gauge
	bootstrap     prometheus.Gauge
	sync          prometheus.Gauge
	fetchInterval time.Duration
}

func (e *tezosExporter) Collect() error {
	go func() {
		for {
			log.Println("Collect Tezos metric")

			err := e.getInfo()
			if err != nil {
				log.Println(err)
				time.Sleep(e.fetchInterval)
				continue
			}
			_ = err
			time.Sleep(e.fetchInterval)
		}
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

	e.bootstrap.Set(bootstrapMap[bootstrapStatus.IsBootstrapped])
	e.sync.Set(syncMap[bootstrapStatus.SyncState])

	headBlock, err := cli.GetHeadBlock()
	if err != nil {
		log.Println(err)
		return err
	}

	e.blockLevel.Set(float64(headBlock.Header.Level))

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

	e.peerCount.Set(float64(len(runningPeers)))

	return err
}

func (e *tezosExporter) Stop() error {
	// TODO:
	return nil
}

func createTezosExporter(rpcEndpoint string, fetchInterval int) Exporter {
	blockLevel := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tezos_block_level",
		Help: "Current block level of Tezos network",
	})

	peerCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tezos_peer_count",
		Help: "Count of connected peers on the node.",
	})

	bootstrap := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tezos_is_bootstrapped",
		Help: "Bootstrap status. 0: No, 1: Yes",
	})

	sync := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tezos_sync_status",
		Help: "Block sync status. 0: stuck, 1: synced, 2: unsynced",
	})

	prometheus.MustRegister(blockLevel)
	prometheus.MustRegister(peerCount)
	prometheus.MustRegister(bootstrap)
	prometheus.MustRegister(sync)

	return &tezosExporter{
		rpcClient:     rpc.NewClient(rpcEndpoint),
		blockLevel:    blockLevel,
		peerCount:     peerCount,
		bootstrap:     bootstrap,
		sync:          sync,
		fetchInterval: time.Duration(fetchInterval) * time.Second,
	}
}
