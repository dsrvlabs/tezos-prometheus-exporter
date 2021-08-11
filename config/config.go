package config

import (
	"flag"
	"log"
	"os"
	"sync"
)

const (
	defaultRPCEndpoint    = "http://localhost:8732"
	defaultServicePort    = 9489
	defaultUpdateInterval = 1
	defaultMountPath      = "/"
)

var (
	once   sync.Once
	config Config
)

func init() {
}

// Config maintain service parameters
type Config struct {
	RPCEndpoint           string `json:"rpc_endpoint"`
	ServicePort           int    `json:"service_port"`
	UpdateIntervalSeconds int    `json:"update_interval_seconds"`
	MountPath             string `json:"mount_path"`
}

// GetConfig returns ...
func GetConfig() Config {
	log.Println("Load configs")

	once.Do(func() {
		log.Println("Check service configs")
		config = parseConfig()
	})

	return config
}

func parseConfig() Config {
	cmdLine := flag.NewFlagSet("test", flag.ContinueOnError)

	var (
		rpcURL        string
		servicePort   int
		fetchInterval int
		mountPath     string
	)

	cmdLine.StringVar(&rpcURL, "rpc-endpoint", defaultRPCEndpoint, "RPC Endpoint")
	cmdLine.IntVar(&servicePort, "prometheus-port", defaultServicePort, "Default service port")
	cmdLine.IntVar(&fetchInterval, "fetch-interval", defaultUpdateInterval, "Update interval(s) ")
	cmdLine.StringVar(&mountPath, "mount-path", defaultMountPath, "Mount path to check free space")

	cmdLine.Parse(os.Args[1:])

	return Config{
		RPCEndpoint:           rpcURL,
		ServicePort:           servicePort,
		UpdateIntervalSeconds: fetchInterval,
		MountPath:             mountPath,
	}
}
