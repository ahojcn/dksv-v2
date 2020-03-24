package controllers

import (
	"dksv-v2/models"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"os/exec"
)

type HostController struct {
	beego.Controller
}

// 查看文件夹里的文件列表
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

// 上传一个文件
// 如果文件存在则覆盖
func (this *HostController) UploadFile() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type createFileForm struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	req := createFileForm{
		Name: this.GetString("name"),
		Path: this.GetString("path"),
	}

	logrus.Warnln("name:", req.Name)
	logrus.Warnln("path:", req.Path)

	f, _, err := this.GetFile("file")
	if err != nil {
		data.Msg = fmt.Sprintf("上传文件失败:%v", err)
		data.Status = -1
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	defer f.Close()

	// 判断 path 是否存在，不存在则创建
	_, err = os.Stat(req.Path)
	if os.IsNotExist(err) {
		// 不存在
		err = os.MkdirAll(req.Path, os.ModePerm)
		if err != nil {
			data.Msg = fmt.Sprintf("创建文件夹失败:%v", err)
			data.Status = -1
			this.Data["json"] = data
			this.ServeJSON()
			return
		}
	}

	err = this.SaveToFile("file", fmt.Sprintf("%s/%s", req.Path, req.Name))
	if err != nil {
		data.Msg = fmt.Sprintf("保存文件失败:%v", err)
		data.Status = -1
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Msg = "上传成功"
	this.Data["json"] = data
	this.ServeJSON()
}

// 寻找本机可用端口
func (this *HostController) UnusedPortList() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	ports := make([]int, 0)
	for i := 10000; i < 10100; i++ {
		if CheckPort(i) != true {
			ports = append(ports, i)
		}
	}

	data.Data = ports
	this.Data["json"] = data
	this.ServeJSON()
}

func CheckPort(port int) bool {
	checkStatement := fmt.Sprintf("lsof -i:%d", port)
	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if len(output) > 0 {
		return true
	}
	return false
}
