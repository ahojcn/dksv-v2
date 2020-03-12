package main

import (
	"bufio"
	"dksv-v2/controllers"
	_ "dksv-v2/routers"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/"] = "/"
	}

	logrus.Infoln("连接控制中心...")
	file, err := os.Open("ip.txt")
	if err != nil {
		logrus.Panicf("读取 ip 文件错误:%v!", err)
	}
	line, _, _ := bufio.NewReader(file).ReadLine()
	path := fmt.Sprintf("http://tim.natapp1.cc/host/init-notify?ip=%s", string(line))
	client := http.Client{Timeout: 5 * time.Second}
	_, err = client.Get(path)
	if err != nil {
		logrus.Panicf("连接控制中心失败:%v!", err)
	}
	logrus.Infoln("连接控制中心成功!")

	//////////////////////////////////////////////////////////////////////
	// 获取 api 的 port 和 version
	logrus.Infoln("获取 api 版本...")
	file, err = os.Open("apiversion.txt")
	if err != nil {
		logrus.Panicf("获取版本失败:%v!", err)
	}
	version, _, _ := bufio.NewReader(file).ReadLine()
	if string(version) == "" {
		logrus.Panicf("获取版本失败!")
	}
	controllers.AC.Version = string(version)
	logrus.Infoln("版本:", string(version))

	logrus.Infoln("获取 api 端口...")
	file, err = os.Open("port.txt")
	if err != nil {
		logrus.Panicf("获取端口失败:%v!", err)
	}
	port, _, _ := bufio.NewReader(file).ReadLine()
	if string(port) == "" {
		logrus.Panicf("获取端口失败!")
	}
	controllers.AC.Version = string(port)
	logrus.Infoln("端口:", string(port))
	//////////////////////////////////////////////////////////////////////

	logrus.Infoln("my-docker api 状态检测...")
	cmd := exec.Command("curl", "127.0.0.1:6060/info")
	err = cmd.Run()
	if err != nil {
		logrus.Panicf("my-docker api 存活检测失败:%v!", err)
	}
	logrus.Infoln("my-docker api 存活!")

	beego.Run()
}
