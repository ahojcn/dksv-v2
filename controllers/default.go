package controllers

import (
	"dksv-v2/models"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"time"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Post() {
	beego.Info("固定路由的get类型的方法 ")
	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	d := make(map[string]interface{})

	// 解析参数
	type containerRemoveForm struct {
		ContainerName string `json:"container_name"`
	}
	req := containerRemoveForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	logrus.Infoln(req.ContainerName)
	containers, _ := docker.GetDockerStat()
	d["containers"] = containers

	cpusInfo := make([]*cpu.TimesStat, 0)
	memsInfo := make([]*docker.CgroupMemStat, 0)
	for index := range containers {
		cpuInfo, _ := docker.CgroupCPUDocker(containers[index].ContainerID)
		cpusInfo = append(cpusInfo, cpuInfo)

		memInfo, _ := docker.CgroupMemDocker(containers[index].ContainerID)
		memsInfo = append(memsInfo, memInfo)
	}
	d["cpusInfo"] = cpusInfo
	d["memsInfo"] = memsInfo

	data.Data = d
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *MainController) SysInfo() {
	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	d := make(map[string]interface{})

	type cpuInfo struct {
		LogicalCores  int       `json:"logical_cores"`  // 逻辑CPU个数
		PhysicalCores int       `json:"physical_cores"` // 物理CPU个数
		Percent       []float64 `json:"percent"`        // CPU 使用量
	}
	lc, _ := cpu.Counts(true)
	pc, _ := cpu.Counts(false)
	percent, _ := cpu.Percent(time.Second, false)
	d["cpu_info"] = cpuInfo{
		LogicalCores:  lc,
		PhysicalCores: pc,
		Percent:       percent,
	}

	// 磁盘使用量
	d["disk_info"], _ = disk.Usage("/")

	// host 信息
	d["host_info"], _ = host.Info()

	// 传感器信息
	d["sensors"], _ = host.SensorsTemperatures()

	// 负载信息
	loadInfo := make(map[string]interface{})
	loadInfo["load"], _ = load.Avg()
	loadInfo["misc"], _ = load.Misc()
	d["load_info"] = loadInfo

	// 内存
	memInfo := make(map[string]interface{})
	memInfo["swap_memory"], _ = mem.SwapMemory()
	memInfo["virtual_memory"], _ = mem.VirtualMemory()
	d["mem_info"] = memInfo

	// 网络
	d["io_counters"], _ = net.IOCounters(false) // 返回网卡收发数据总和
	proto := []string{"ip", "icmp", "icmpmsg", "tcp", "udp", "udplite"}
	d["proto_counters"], _ = net.ProtoCounters(proto)

	// 进程
	//pids, _ := process.Pids()
	//procs := make([]map[string]interface{}, 10)
	//for index := range pids {
	//	proc := process.Process{Pid:pids[index]}
	//}
	d["process"], _ = process.Pids()

	data.Data = d
	this.Data["json"] = data
	this.ServeJSON()
}
