package controllers

import (
	"github.com/docker/docker/client"
)

// 获取 moby cli
func getMobyCli() (*client.Client, error) {
	cli, err := client.NewClient("tcp://139.159.254.242:6060", "v1.40", nil, nil)

	if err != nil {
		return nil, err
	}
	return cli, err
}
