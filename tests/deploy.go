package test

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("nohup", "dksv2.0.linux-amd64", "&")
	fmt.Println("start...")
	err := cmd.Start()
	fmt.Println(err)
}
