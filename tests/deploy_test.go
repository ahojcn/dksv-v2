package test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestDeploy001(t *testing.T) {
	ip, _ := ioutil.ReadFile("ip.txt")
	url := fmt.Sprintf("http://tim.natapp1.cc/host/init-notify?ip=%s", string(ip))
	fmt.Println("url:", url)
	cmd := exec.Command("curl", url)
	err := cmd.Run()
	fmt.Println("err:", err)
}
