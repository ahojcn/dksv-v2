package controllers

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/net"
	//"os"
	//"runtime"
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/shirou/gopsutil/process"
)

func TestGetSysInfo(t *testing.T) {
	fmt.Println(mem.VirtualMemory()) // 虚拟内存
	fmt.Println(mem.SwapMemory())    // 交换内存

	fmt.Println(cpu.Info()) // cpu info

	fmt.Println(disk.Partitions(true)) // 分区返回磁盘分区。
	fmt.Println(disk.Usage("/"))       // 磁盘使用量

	fmt.Println(host.Info())                // 操作系统等等信息
	fmt.Println(host.SensorsTemperatures()) // 传感器温度信息
	fmt.Println(host.Users())               // 用户信息

	fmt.Println(load.Avg()) // 负载平均统计

	fmt.Println(net.IOCounters(true))   // 网络IO统计
	fmt.Println(net.ProtoCounters(nil)) // 协议统计
	fmt.Println(net.Connections("all")) // 连接状态统计

	fmt.Println(process.Pids()) //进程统计
}
