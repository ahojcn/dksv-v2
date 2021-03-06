package controllers

import (
	"bufio"
	"context"
	"dksv-v2/models"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/docker"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type ContainerController struct {
	beego.Controller
}

// 创建容器
// status -> created
func (this *ContainerController) Create() {
	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	type v struct {
		HostVolume      string `json:"host_volume"`
		ContainerVolume string `json:"container_volume"`
	}
	type containerCreateForm struct {
		ImageName          string   `json:"image_name"`
		WorkingDir         string   `json:"working_dir"`          // 工作路径
		ContainerName      string   `json:"container_name"`       // 不起名字就可以填 ""
		ContainerPortProto string   `json:"container_port_proto"` // 容器端口协议 tcp / udp
		ContainerPort      string   `json:"container_port"`       // 端口 80
		HostPort           string   `json:"host_port"`            // 主机端口
		CPUShares          int64    `json:"cpu_shares"`           // CPUShares 默认 1024
		Memory             int64    `json:"memory"`               // 内存限制 bytes
		Cmd                []string `json:"cmd"`                  // run 时候执行的命令
		Volumes            []v      `json:"volumes"`
	}
	req := containerCreateForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	logrus.Infoln("request body:", this.Ctx.Input.RequestBody)
	logrus.Infoln("Cmd: ", req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	exports := make(nat.PortSet, 3)
	//port, _ := nat.NewPort("tcp", "80")
	port, _ := nat.NewPort(req.ContainerPortProto, req.ContainerPort)
	exports[port] = struct{}{}

	//c := make([]string, 0)
	//c = append(c, req.Cmd)
	config := &container.Config{
		ExposedPorts: exports,
		Cmd:          req.Cmd,
		Image:        req.ImageName,
	}

	portBind := nat.PortBinding{HostPort: req.HostPort}
	portMap := make(nat.PortMap, 0)
	tmp := make([]nat.PortBinding, 0, 1)
	tmp = append(tmp, portBind)
	portMap[port] = tmp

	mnt := make([]mount.Mount, 0)
	for index := range req.Volumes {
		mnt = append(mnt, mount.Mount{
			Type:   mount.TypeBind, /// 注意这里的类型
			Source: req.Volumes[index].HostVolume,
			Target: req.Volumes[index].ContainerVolume,
		})
		logrus.Errorln([]byte(req.Volumes[index].HostVolume))
	}
	hostConfig := &container.HostConfig{
		PortBindings: portMap,
		Mounts:       mnt,
	}
	hostConfig.Resources = container.Resources{
		CPUShares: req.CPUShares, // CPU共享(相对于其他容器的相对重量)
		Memory:    req.Memory,    // Memory limit (in bytes)
	}

	containerInfo, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, req.ContainerName)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("创建容器失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Data = containerInfo
	cli.Info(context.Background())

	this.Data["json"] = data
	this.ServeJSON()
}

// 运行一个已经存在的容器
func (this *ContainerController) Start() {
	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	type containerStartForm struct {
		ContainerName string `json:"container_name"`
	}
	req := containerStartForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	err = cli.ContainerStart(context.Background(), req.ContainerName, types.ContainerStartOptions{
		CheckpointID:  "",
		CheckpointDir: "",
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("启动容器失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 检查容器信息
// 有关容器的底层信息。
func (this *ContainerController) Inspect() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	// 解析参数
	containerName := this.GetString("container_name")

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	info, err := cli.ContainerInspect(context.Background(), containerName)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("查询信息出错:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Data = info
	this.Data["json"] = data
	this.ServeJSON()
}

// 获取容器 stat 信息
func (this *ContainerController) Stat() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	containerName := this.GetString("container_name")

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	stats, err := cli.ContainerStats(context.Background(), containerName, false)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("未知错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	b, err := ioutil.ReadAll(stats.Body)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("转换错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	type MyJsonName struct {
		BlkioStats struct{} `json:"blkio_stats"`
		CPUStats   struct {
			CPUUsage struct {
				PercpuUsage       []int64 `json:"percpu_usage"`
				TotalUsage        int64   `json:"total_usage"`
				UsageInKernelmode int64   `json:"usage_in_kernelmode"`
				UsageInUsermode   int64   `json:"usage_in_usermode"`
			} `json:"cpu_usage"`
			OnlineCpus     int64 `json:"online_cpus"`
			SystemCPUUsage int64 `json:"system_cpu_usage"`
			ThrottlingData struct {
				Periods          int64 `json:"periods"`
				ThrottledPeriods int64 `json:"throttled_periods"`
				ThrottledTime    int64 `json:"throttled_time"`
			} `json:"throttling_data"`
		} `json:"cpu_stats"`
		MemoryStats struct {
			Failcnt  int64 `json:"failcnt"`
			Limit    int64 `json:"limit"`
			MaxUsage int64 `json:"max_usage"`
			Stats    struct {
				ActiveAnon              int64 `json:"active_anon"`
				ActiveFile              int64 `json:"active_file"`
				Cache                   int64 `json:"cache"`
				HierarchicalMemoryLimit int64 `json:"hierarchical_memory_limit"`
				InactiveAnon            int64 `json:"inactive_anon"`
				InactiveFile            int64 `json:"inactive_file"`
				MappedFile              int64 `json:"mapped_file"`
				Pgfault                 int64 `json:"pgfault"`
				Pgmajfault              int64 `json:"pgmajfault"`
				Pgpgin                  int64 `json:"pgpgin"`
				Pgpgout                 int64 `json:"pgpgout"`
				Rss                     int64 `json:"rss"`
				RssHuge                 int64 `json:"rss_huge"`
				TotalActiveAnon         int64 `json:"total_active_anon"`
				TotalActiveFile         int64 `json:"total_active_file"`
				TotalCache              int64 `json:"total_cache"`
				TotalInactiveAnon       int64 `json:"total_inactive_anon"`
				TotalInactiveFile       int64 `json:"total_inactive_file"`
				TotalMappedFile         int64 `json:"total_mapped_file"`
				TotalPgfault            int64 `json:"total_pgfault"`
				TotalPgmajfault         int64 `json:"total_pgmajfault"`
				TotalPgpgin             int64 `json:"total_pgpgin"`
				TotalPgpgout            int64 `json:"total_pgpgout"`
				TotalRss                int64 `json:"total_rss"`
				TotalRssHuge            int64 `json:"total_rss_huge"`
				TotalUnevictable        int64 `json:"total_unevictable"`
				TotalWriteback          int64 `json:"total_writeback"`
				Unevictable             int64 `json:"unevictable"`
				Writeback               int64 `json:"writeback"`
			} `json:"stats"`
			Usage int64 `json:"usage"`
		} `json:"memory_stats"`
		Networks struct {
			Eth0 struct {
				RxBytes   int64 `json:"rx_bytes"`
				RxDropped int64 `json:"rx_dropped"`
				RxErrors  int64 `json:"rx_errors"`
				RxPackets int64 `json:"rx_packets"`
				TxBytes   int64 `json:"tx_bytes"`
				TxDropped int64 `json:"tx_dropped"`
				TxErrors  int64 `json:"tx_errors"`
				TxPackets int64 `json:"tx_packets"`
			} `json:"eth0"`
			Eth5 struct {
				RxBytes   int64 `json:"rx_bytes"`
				RxDropped int64 `json:"rx_dropped"`
				RxErrors  int64 `json:"rx_errors"`
				RxPackets int64 `json:"rx_packets"`
				TxBytes   int64 `json:"tx_bytes"`
				TxDropped int64 `json:"tx_dropped"`
				TxErrors  int64 `json:"tx_errors"`
				TxPackets int64 `json:"tx_packets"`
			} `json:"eth5"`
		} `json:"networks"`
		PidsStats struct {
			Current int64 `json:"current"`
		} `json:"pids_stats"`
		PrecpuStats struct {
			CPUUsage struct {
				PercpuUsage       []int64 `json:"percpu_usage"`
				TotalUsage        int64   `json:"total_usage"`
				UsageInKernelmode int64   `json:"usage_in_kernelmode"`
				UsageInUsermode   int64   `json:"usage_in_usermode"`
			} `json:"cpu_usage"`
			OnlineCpus     int64 `json:"online_cpus"`
			SystemCPUUsage int64 `json:"system_cpu_usage"`
			ThrottlingData struct {
				Periods          int64 `json:"periods"`
				ThrottledPeriods int64 `json:"throttled_periods"`
				ThrottledTime    int64 `json:"throttled_time"`
			} `json:"throttling_data"`
		} `json:"precpu_stats"`
		Read string `json:"read"`
	}
	d := MyJsonName{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("转换错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Data = d
	this.Data["json"] = data
	this.ServeJSON()
}

// 停止容器
func (this *ContainerController) Stop() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type containerStopForm struct {
		ContainerName string `json:"container_name"`
		TimeOut       int64  `json:"time_out"`
	}
	req := containerStopForm{}
	//this.Ctx.Input.CopyBody()
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	logrus.Errorln("stop json.unmarshal err:", err)
	logrus.Infoln("stop() req:", req)

	logrus.Warnln("get string:", this.GetString("container_name"))

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	t := time.Duration(req.TimeOut)
	err = cli.ContainerStop(context.Background(), req.ContainerName, &t)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("停止容器失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Msg = "停止容器成功"
	this.Data["json"] = data
	this.ServeJSON()
}

// 删除容器
func (this *ContainerController) Remove() {
	logrus.Infoln("Remove()...")
	data := &models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type containerRemoveForm struct {
		ContainerName string `json:"container_name"`
		Force         bool   `json:"force"` // -f
	}
	req := containerRemoveForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误: %v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	err = cli.ContainerRemove(context.Background(), req.ContainerName, types.ContainerRemoveOptions{
		RemoveVolumes: false,
		RemoveLinks:   false,
		Force:         req.Force,
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("删除容器失败: %v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Msg = "删除容器成功"
	this.Data["json"] = data
	this.ServeJSON()
}

// 列出容器
func (this *ContainerController) List() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	type myContainers struct {
		Container types.Container       `json:"container"`
		MemInfo   *docker.CgroupMemStat `json:"mem_info"`
		CpuInfo   *cpu.TimesStat        `json:"cpu_info"`
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Quiet:   false,
		Size:    false,
		All:     true,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filters.Args{},
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	d := make([]myContainers, len(containers))
	for i := range containers {
		d[i].Container = containers[i]
		d[i].CpuInfo, _ = docker.CgroupCPUDocker(containers[i].ID)
		d[i].MemInfo, _ = docker.CgroupMemDocker(containers[i].ID)
	}

	data.Data = d

	this.Data["json"] = data
	this.ServeJSON()
}

// 查看容器的日志
func (this *ContainerController) Logs() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	containerName := this.GetString("container_name")

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	logs, err := cli.ContainerLogs(context.Background(), containerName, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      "",
		Until:      "",
		Timestamps: true,
		Follow:     false,
		Tail:       "",
		Details:    true,
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("查看日志失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	containerLogs := make([]string, 1)
	r := bufio.NewReader(logs)
	for {
		//buf, e := r.ReadBytes('\n')
		buf, e := r.ReadString('\n')
		if e != nil && len(buf) == 0 {
			break
		}
		containerLogs = append(containerLogs, string(buf))
	}

	data.Data = containerLogs

	this.Data["json"] = data
	this.ServeJSON()
}

// 将本机的容器打包成压缩包
func (this *ContainerController) Export() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type containerExportForm struct {
		ContainerName string `json:"container_name"`
		ImageName     string `json:"image_name"`
	}
	req := containerExportForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	err := os.MkdirAll(models.RootUrl, os.ModePerm)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("导出容器错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	cmd := exec.Command("docker", "export", req.ContainerName, "-o", fmt.Sprintf("%s/%s", models.RootUrl, req.ImageName))
	err = cmd.Run()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("导出容器错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Msg = "导出成功"
	this.Data["json"] = data
	this.ServeJSON()
}

func (this *ContainerController) Commit() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type containerCommitForm struct {
		ContainerName string `json:"container_name"`
		Ref string `json:"ref"`
	}
	req := containerCommitForm{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &req)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("解析参数错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("未知错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	logrus.Warningln("ContainerName:", req.ContainerName, "Ref:", req.Ref)
	id, err := cli.ContainerCommit(context.Background(), req.ContainerName, types.ContainerCommitOptions{
		Reference: req.Ref,
		Comment:   "",
		Author:    "",
		Changes:   nil,
		Pause:     false,
		Config:    nil,
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("未知错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Data = id
	this.Data["json"] = data
	this.ServeJSON()
}
