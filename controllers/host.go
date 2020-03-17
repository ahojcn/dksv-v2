package controllers

import (
	"dksv-v2/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
)

type HostController struct {
	beego.Controller
}

// 查看文件夹里的文件 / 查看文件详情
func (this *HostController) ListFiles() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	d := make(map[string]interface{})
	// 解析参数
	path := this.GetString("path")
	if path == "" {
		path = "/root"
	}
	pathInfo, err := os.Stat(path)
	if err != nil {
		data.Msg = fmt.Sprintf("未知错误:%v", err)
		data.Status = -1
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	d["name"] = pathInfo.Name()
	d["is_dir"] = pathInfo.IsDir()
	d["size"] = pathInfo.Size()

	children := make([]interface{}, 0)
	if pathInfo.IsDir() {
		fileInfos, err := ioutil.ReadDir(path)
		if err != nil {
			data.Msg = fmt.Sprintf("未知错误:%v", err)
			data.Status = -1
			this.Data["json"] = data
			this.ServeJSON()
			return
		}
		for i := range fileInfos {
			child := make(map[string]interface{})
			child["is_dir"] = fileInfos[i].IsDir()
			child["name"] = fileInfos[i].Name()
			child["size"] = fileInfos[i].Size()
			children = append(children, child)
		}
	}
	d["children"] = children

	data.Data = d
	this.Data["json"] = data
	this.ServeJSON()
}

// 创建一个文件
func (this *HostController) CreateFile() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type createFileForm struct {
		Name    string   `json:"name"`
		Path    string   `json:"path"`
		Content []string `json:"content"`
	}
	req := createFileForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	this.Data["json"] = data
	this.ServeJSON()
}
