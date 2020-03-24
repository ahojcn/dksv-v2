package controllers

import (
	"dksv-v2/models"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestNetworkController_Create(t *testing.T) {
	f, err := os.Open(models.DefaultNetworkPath + "bridxxx")
	fmt.Println(f, err)
}

func TestHostController_UnusedPortList(t *testing.T) {
	fmt.Println(CheckPort(8080))
}

func CheckPort(port int) bool {
	checkStatement := fmt.Sprintf("lsof -i:%d ", port)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if len(output) > 0 {
		return true
	}
	return false
}
