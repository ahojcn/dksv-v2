package controllers

import (
	"github.com/docker/docker/client"
)

type ApiConfig struct {
	Port string `json:"port"`
	Version string `json:"version"`
}

var AC ApiConfig

// 获取 moby cli
func getMobyCli() (*client.Client, error) {
	cli, err := client.NewClient("tcp://139.159.254.242:6060", "v1.40", nil, nil)
	//cli, err := client.NewClient(fmt.Sprintf("tcp://127.0.0.1:%s", AC.Port), AC.Version, nil, nil)

	if err != nil {
		return nil, err
	}
	return cli, err
}
