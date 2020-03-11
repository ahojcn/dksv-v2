package test

import (
	"context"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
	"testing"
)

func TestMoby(t *testing.T) {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	cli, err := client.NewClient("tcp://139.159.254.242:6060", "v1.40", nil, nil)
	logrus.Infoln(cli)
	logrus.Errorln("1", err)

	listImage(cli)
}

func listImage(cli *client.Client) {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	logrus.Errorln("2", err)

	for _, image := range images {
		logrus.Infof("%+v", image)
	}
}

func createContainer(cli *client.Client) {
	//cli.ContainerCreate(context.Background(), )
}
