package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

const (
	defaultRPCEndpoint    = "http://localhost:8732"
	defaultServicePort    = 9489
	defaultUpdateInterval = 1
	defaultMountPath      = "/"
)

var (
	defaultDaemons = []string{"tezos-node"}
)

// Config maintain service parameters
type Config struct {
	RPCEndpoint           string   `json:"rpc-addr"`
	ExporterPort          int      `json:"exporter-port"`
	UpdateIntervalSeconds int      `json:"refresh-interval"`
	DataDir               string   `json:"data-dir"`
	Daemons               []string `json:"daemons"`
}

// Loader privides interfaces to load exporter config.
type Loader interface {
	Load() (Config, error)
}

type loader struct {
	config   Config
	filename string
}

func (l *loader) Load() (Config, error) {
	log.Println("Load configs")

	loadConfig, err := l.loadConfigFile(l.filename)
	if err != nil {
		log.Println("Can't find config file. Set default config")

		loadConfig = &Config{}
		loadConfig.RPCEndpoint = defaultRPCEndpoint
		loadConfig.ExporterPort = defaultServicePort
		loadConfig.UpdateIntervalSeconds = defaultUpdateInterval
		loadConfig.DataDir = defaultMountPath
		loadConfig.Daemons = defaultDaemons
	}

	l.config = *loadConfig

	return l.config, nil
}

func (l *loader) loadConfigFile(filename string) (*Config, error) {
	log.Println("Load config file", filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	rawContent, err := io.ReadAll(f)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	exporterConfig := Config{}

	err = json.Unmarshal(rawContent, &exporterConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &exporterConfig, nil
}

// NewLoader creates loader to manager config.
func NewLoader(filename string) Loader {
	return &loader{filename: filename}
}
