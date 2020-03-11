package controllers

import (
	"bytes"
	"context"
	"dksv-v2/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
	"io/ioutil"
	"mime/multipart"
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
		ImageName string `json:"image_name"`
	}
	req := imagePullForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	f, err := cli.ImagePull(context.Background(), req.ImageName, types.ImagePullOptions{
		All:           false,
		RegistryAuth:  "",
		PrivilegeFunc: nil,
		Platform:      "",
	})
	defer f.Close()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("拉取镜像错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	for {
		p := make([]byte, 1024)
		n, err := f.Read(p)
		if n == 0 && err == io.EOF {
			data.Msg = "拉取镜像成功"
			data.Status = 0
			break
		} else if err != nil && err != io.EOF {
			// 报错
			data.Msg = fmt.Sprintf("拉取镜像失败:%v", err)
			data.Status = -1
			break
		}
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 列出本机的镜像
func (this *ImageController) List() {
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

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{
		All:     true,
		Filters: filters.Args{},
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("获取镜像列表失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Data = images

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
		Force bool `json:"force"`
	}
	req := imageRemoveForm{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	_, err = cli.ImageRemove(context.Background(), req.ImageName, types.ImageRemoveOptions{
		Force:         req.Force,
		PruneChildren: false,
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("删除镜像失败:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	data.Msg = "删除成功"

	this.Data["json"] = data
	this.ServeJSON()
}

// 将本机镜像推送到服务器
func (this *ImageController) Push() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type imagePushForm struct {
		ImageName string `json:"image_name"`
		Url string `json:"url"`
	}
	req := imagePushForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	// 打开文件
	f, err := os.Open(fmt.Sprintf("%s/%s", models.RootUrl, req.ImageName))
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("没有此镜像文件信息:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	defer f.Close()

	resp, err := uploadFile(req.Url, map[string]string{}, "file", req.ImageName, f)
	data.Data = resp

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

// 上传 image.tar 到镜像管理服务器
func uploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	HttpClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}
