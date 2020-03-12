package controllers

import (
	"context"
	"dksv-v2/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"io"
)

type ImageController struct {
	beego.Controller
}

// 从服务器拉取镜像
// docker pull image_name
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
	if err != nil || f == nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("拉取镜像错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	defer f.Close()

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

// 列出本机的 docker 镜像
// docker images
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

// 删除本机的 docker 镜像
// docker rmi image_name -f
// image_name:
func (this *ImageController) Remove() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type imageRemoveForm struct {
		ImageName string `json:"image_name"`
		Force     bool   `json:"force"`
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

// 给 image 打 tag
// docker tag image_name target_image_name
// image_name: nginx
// target_image_name: 139.159.254.242:5000/myngx:0.1
func (this *ImageController) Tag() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type imageTagForm struct {
		ImageName       string `json:"image_name"`
		TargetImageName string `json:"target_image_name"`
	}
	req := imageTagForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("没有此镜像文件信息:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}

	err = cli.ImageTag(context.Background(), req.ImageName, req.TargetImageName)
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("%s镜像tag失败:%v", req.ImageName, err)
	}

	data.Msg = "成功"
	this.Data["json"] = data
	this.ServeJSON()
}

// 推送本地镜像到仓库
// docker push image_name
// image_name: 139.159.254.242:5000/tomcat:latest
// path: 139.159.254.242:5000
func (this *ImageController) Push() {
	data := models.RESDATA{
		Status: 0,
		Msg:    "success",
		Data:   nil,
	}
	// 解析参数
	type imagePushForm struct {
		ImageName string `json:"image_name"`
		Path string `json:"path"`
	}
	req := imagePushForm{}
	_ = json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	cli, err := getMobyCli()
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("网络错误:%v", err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	f, err := cli.ImagePush(context.Background(), req.ImageName, types.ImagePushOptions{
		All:           false,
		RegistryAuth:  req.Path,
		PrivilegeFunc: nil,
		Platform:      "",
	})
	if err != nil {
		data.Status = -1
		data.Msg = fmt.Sprintf("没有此镜像%s信息:%v", req.ImageName, err)
		this.Data["json"] = data
		this.ServeJSON()
		return
	}
	for {
		buf := make([]byte, 1024)
		n, err := f.Read(buf)
		if n == 0 && err == io.EOF {
			data.Msg = "push成功"
			data.Status = 0
			data.Data = req
			break
		} else if err != nil && err != io.EOF {
			data.Msg = fmt.Sprintf("push失败:%v", err)
			data.Status = 0
			data.Data = req
			break
		}
	}

	this.Data["json"] = data
	this.ServeJSON()
}

// 从文件服务器下载镜像
//func (this *ImageController) Download() {
//	data := models.RESDATA{
//		Status: 0,
//		Msg:    "success",
//		Data:   nil,
//	}
//
//	// 解析参数
//	type imagePullForm struct {
//		ImageUrl  string `json:"image_url"`
//		ImageName string `json:"image_name"`
//	}
//	req := imagePullForm{}
//	json.Unmarshal(this.Ctx.Input.RequestBody, &req)
//
//	// 从 image_url 下载文件到本地镜像存储路径
//	fileURL := req.ImageUrl
//	filePath := models.RootUrl
//	// 要下载的文件并不是 .tar 结尾
//	//if !strings.HasSuffix(path.Base(fileURL), ".tar") {
//	//	data.Status = -1
//	//	data.Msg = "镜像文件格式错误"
//	//	this.Data["json"] = data
//	//	this.ServeJSON()
//	//	return
//	//}
//
//	res, err := http.Get(fileURL)
//	if err != nil || res.Status != "200" || strings.Index(res.Status, "200") != -1 {
//		data.Status = -1
//		logrus.Errorln(res.Status)
//		data.Msg = fmt.Sprintf("文件地址错误:%v %s %s", err, fileURL, res.Status)
//		this.Data["json"] = data
//		this.ServeJSON()
//		return
//	}
//	defer res.Body.Close()
//
//	// 获得 get 请求响应的reader对象
//	reader := bufio.NewReaderSize(res.Body, 32*1024)
//	file, err := os.Create(filePath + req.ImageName + ".tar")
//	if err != nil {
//		data.Status = -1
//		data.Msg = fmt.Sprintf("创建本地镜像文件错误:%v", err)
//		this.Data["json"] = data
//		this.ServeJSON()
//		return
//	}
//	// 获得文件的writer对象
//	writer := bufio.NewWriter(file)
//	io.Copy(writer, reader)
//}
