package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

// Client provides access interfaces.
type Client interface {
	GetHeadBlock() (*Block, error)
	GetPeers() ([]Peer, error)
	GetBootstrapStatus() (*BootstrapStatus, error)
}

type client struct {
	HostAddr   string
	httpClient *http.Client
}

func (c *client) GetHeadBlock() (*Block, error) {
	log.Println("GetHeadBlock")

	url := c.HostAddr + "/chains/main/blocks/head"

	log.Println("URL " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RPC Failed with code %d", resp.StatusCode)
	}

	rawBody, _ := io.ReadAll(resp.Body)

	headBlock := Block{}
	err = json.Unmarshal(rawBody, &headBlock)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &headBlock, nil
}

func (c *client) GetPeers() ([]Peer, error) {
	log.Println("GetPeers")

	url := c.HostAddr + "/network/peers"

	log.Println("URL " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RPC Failed with code %d", resp.StatusCode)
	}

	rawBody, _ := io.ReadAll(resp.Body)

	data := []interface{}{}
	err = json.Unmarshal(rawBody, &data)

	peers := make([]Peer, 0)

	for _, item := range data {
		values := item.([]interface{})

		peerID := values[0].(string)
		peer := Peer{ID: peerID}

		peerData := values[1].(map[string]interface{})
		mapstructure.Decode(peerData, &peer)

		peers = append(peers, peer)
	}

	return peers, nil
}

func (c *client) GetBootstrapStatus() (*BootstrapStatus, error) {
	log.Println("GetBootstrapStatus")

	url := c.HostAddr + "/chains/main/is_bootstrapped"

	log.Println("URL " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RPC Failed with code %d", resp.StatusCode)
	}

	rawBody, _ := io.ReadAll(resp.Body)

	status := BootstrapStatus{}

	err = json.Unmarshal(rawBody, &status)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &status, nil
}

// NewClient creates new instance of RPC client.
func NewClient(host string) Client {
	log.Println("Create New Client")
	return &client{
		HostAddr:   host,
		httpClient: &http.Client{},
	}
}
