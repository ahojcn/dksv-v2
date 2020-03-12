package main

import (
	_ "dksv-v2/routers"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os/exec"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/"] = "/"
	}

	go func() {
		logrus.Infoln("初始化成功，连接控制中心中...")
		ip, _ := ioutil.ReadFile("ip.txt")
		logrus.Infoln(string(ip))
		url := fmt.Sprintf("http://tim.natapp1.cc/host/init-notify?ip=%s", string(ip))
		cmd := exec.Command("curl", url)
		for i := 0; i < 100; i++ {
			err := cmd.Run()
			if err == nil {
				logrus.Infoln("连接控制中心成功!")
				break
			}
		}
	}()

	go func() {
		logrus.Infoln("mydocker状态检测中...")
		cmd:=exec.Command("curl", "127.0.0.1:6060/info")
		err := cmd.Run()
		logrus.Infoln(err)
		if err != nil {
			logrus.Errorln("mydocker存活检测失败!")
		} else {
			logrus.Infoln("mydocker存活!")
		}
	}()

	beego.Run()
}
