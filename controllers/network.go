package controllers

import (
	"context"
	"dksv-v2/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

type NetworkController struct {
	beego.Controller
}

// 创建网络
func (this *NetworkController) Create() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type networkCreateForm struct {
		Name string `json:"name"`
	}
	req := networkCreateForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	nw, err := cli.NetworkCreate(context.Background(), req.Name, types.NetworkCreate{
		CheckDuplicate: false,
		Driver:         "bridge",
		Scope:          "",
		EnableIPv6:     false,
		IPAM: &network.IPAM{
			Driver:  "",
			Options: nil,
			Config:  nil,
		},
		Internal:   false,
		Attachable: false,
		Ingress:    false,
		ConfigOnly: false,
		ConfigFrom: nil,
		Options:    nil,
		Labels:     nil,
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("创建网络失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Data = nw

	this.Data["json"] = data
	this.ServeJSON()
}

// 查看网络
func (this *NetworkController) List() {
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
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filters.Args{},
	})

	data.Data = networks

	this.Data["json"] = data
	this.ServeJSON()
}

// 删除网络
func (this *NetworkController) Remove() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type networkRemoveForm struct {
		Name string `json:"name"`
	}
	req := networkRemoveForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	err = cli.NetworkRemove(context.Background(), req.Name)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("删除网络失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	this.Data["json"] = data
	this.ServeJSON()
}
