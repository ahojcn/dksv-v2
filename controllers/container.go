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
		ImageName          string `json:"image_name"`
		Volumes            []v    `json:"volumes"`
		WorkingDir         string `json:"working_dir"`          // 工作路径
		ContainerName      string `json:"container_name"`       // 不起名字就可以填 ""
		ContainerPortProto string `json:"container_port_proto"` // 容器端口协议 tcp / udp
		ContainerPort      string `json:"container_port"`       // 端口 80
		HostPort           string `json:"host_port"`            // 主机端口
		CPUShares          int64  `json:"cpu_shares"`           // CPUShares 默认 1024
		Memory             int64  `json:"memory"`               // 内存限制 bytes
	}
	req := containerCreateForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

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
	config := &container.Config{
		ExposedPorts: exports,
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
			Type:   mount.TypeVolume,
			Source: req.Volumes[index].HostVolume,
			Target: req.Volumes[index].ContainerVolume,
		})
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
	type containerCreateForm struct {
		ContainerName string `json:"container_name"`
	}
	req := containerCreateForm{}
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
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

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
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

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
	data.Data = containers

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
