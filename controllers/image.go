package controllers

import (
	"bufio"
	"dksv-v2/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type ImageController struct {
	beego.Controller
}

// 从服务器拉取镜像
func (this *ImageController) Pull() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	// 解析参数
	type imagePullForm struct {
		ImageUrl  string `json:"image_url"`
		ImageName string `json:"image_name"`
	}
	req := imagePullForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	// 从 image_url 下载文件到本地镜像存储路径
	fileURL := req.ImageUrl
	filePath := models.RootUrl

	res, err := http.Get(fileURL)
	if err != nil || res.Status != "200 OK" {
		data.Status = -1
		data.Msg = fmt.Sprintf("文件地址错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	defer res.Body.Close()

	// 获得 get 请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(filePath + req.ImageName + ".tar")
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("创建本地镜像文件错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)

	this.Data["json"] = getImageInfoByName(req.ImageName)
	this.ServeJSON()
}

// 列出本机的镜像
func (this *ImageController) List() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}

	data.Data = *getAllImageInfo()

	this.Data["json"] = data
	this.ServeJSON()
}

// 删除本机的镜像
func (this *ImageController) Remove() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type imageRemoveForm struct {
		ImageName string `json:"image_name"`
	}
	req := imageRemoveForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	imageName := models.RootUrl + req.ImageName + ".tar"
	err := os.Remove(imageName)
	// 删除文件失败
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("删除镜像失败:%v", err)
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 将本机镜像推送到服务器
// todo
func (this *ImageController) Push() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// todo 将本机镜像推送到服务器
	this.Data["json"] = data
	this.ServeJSON()
}

func getImageInfoByName(imageName string) *models.ImageInfo {
	f, err := os.Open(models.RootUrl + imageName + ".tar")

	if err != nil {
		return nil
	}

	info, err := f.Stat()
	if err != nil {
		return nil
	}

	return &models.ImageInfo{
		Name:    info.Name(),
		Sys:     info.Sys(),
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}
}

// 获取本机所有镜像文件信息
func getAllImageInfo() *[]models.ImageInfo {
	images := make([]models.ImageInfo, 0)
	files, _ := ioutil.ReadDir(models.RootUrl)
	for index := range files {
		f := files[index]
		if strings.HasSuffix(f.Name(), ".tar") {
			images = append(images, *getImageInfo(f))
		}
	}

	return &images
}

// 根据镜像文件获取单个镜像文件的信息
func getImageInfo(f os.FileInfo) *models.ImageInfo {
	return &models.ImageInfo{
		Name:    f.Name(),
		Sys:     f.Sys(),
		ModTime: f.ModTime(),
		Size:    f.Size(),
	}
}
