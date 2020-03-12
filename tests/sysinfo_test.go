package test

import (
	"github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/cpu"
	"testing"
)

func TestGetSysInfo(t *testing.T) {
	infos, _ := cpu.Info()
	logrus.Infoln(infos)
}
