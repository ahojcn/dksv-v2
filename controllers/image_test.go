package controllers

import (
	"bytes"
	"context"
	"dksv-v2/models"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
)

// 测试查看文件
func TestImageController_List(t *testing.T) {
	files, _ := ioutil.ReadDir(models.RootUrl)
	for index := range files {
		// fmt.Println(strings.HasSuffix(files[index].Name(), ".tar"), files[index].Name())
		fname := files[index].Name()
		if strings.HasSuffix(fname, ".tar") {
			fmt.Println(files[index].Sys())
			fmt.Println(files[index].Mode())    // 读写权限
			fmt.Println(files[index].ModTime()) // 修改时间
			fmt.Println(files[index].Size())    // 大小
		}
	}
}

// 测试从hub上,下载文件
func TestImageController_Pull(t *testing.T) {
	cli, _ := getMobyCliTest()
	f, _ := cli.ImagePull(context.Background(), "python", types.ImagePullOptions{
		All:           false,
		RegistryAuth:  "",
		PrivilegeFunc: nil,
		Platform:      "",
	})
	logrus.Infof("%T", f)

	for {
		p := make([]byte, 1024)
		n, err := f.Read(p)
		if n == 0 && err == io.EOF {
			logrus.Info("ok!")
			break
		} else if err != nil {
			// 报错
		}
	}
}

// 测试 push 镜像
func TestImageController_Push2(t *testing.T) {
	//cli, _ := getMobyCliTest()
	//
	//_ = cli.ImageTag(context.Background(), "", "")
	//f, _ := cli.ImagePush(context.Background(), "laughing_engelbart", types.ImagePushOptions{
	//	All:           true,
	//	RegistryAuth:  "",
	//	PrivilegeFunc: nil,
	//	Platform:      "",
	//})
	//logrus.Errorln(f, "push2")
	//
	//for {
	//	p := make([]byte, 1024)
	//	n, err := f.Read(p)
	//	if n == 0 && err == io.EOF {
	//		logrus.Info("ok!", "push2")
	//		break
	//	} else if err != nil && err != io.EOF {
	//		// 报错
	//		logrus.Errorln("error", err)
	//	}
	//}


	cli, _ := getMobyCliTest()
	f, err := cli.ImagePush(context.Background(), "139.159.254.242:5000/myngx:0.2", types.ImagePushOptions{
		All:           true,
		RegistryAuth:  "139.159.254.242:5000",
		PrivilegeFunc: nil,
		Platform:      "",
	})
	logrus.Warnln(err)

	for {
		buf := make([]byte, 1024)
		n, err := f.Read(buf)
		if n == 0 && err == io.EOF {
			logrus.Info("push ok!")
			break
		} else if err != nil && err != io.EOF {
			// 报错
			logrus.Errorln("error", err)
		}
	}
}

func TestImageController_Tag(t *testing.T) {
	cli, _ := getMobyCliTest()
	// docker tag nginx 139.159.254.242:5000/nginx
	err := cli.ImageTag(context.Background(), "nginx", "139.159.254.242:5000/myngx:0.2")
	logrus.Infoln(err)
}

// 获取 moby cli
func getMobyCliTest() (*client.Client, error) {
	cli, err := client.NewClient("tcp://139.159.254.242:6060", "v1.40", nil, nil)

	if err != nil {
		return nil, err
	}
	return cli, err
}

func TestImageController_Push(t *testing.T) {
	_, err := os.Open("/Users/ahojcn/plantumlsss.jar")
	logrus.Errorln(err)

	//resp, err := UploadFileTest("http://tim.natapp1.cc/images/upload", map[string]string{}, "file", "plantuml.jar", f)
	//logrus.Infoln(json.Unmarshal(resp, ), err)
}
func UploadFileTest(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
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

func TestImageController_Search(t *testing.T) {
	cli, _ := client.NewClient("tcp://139.159.254.242:6060", "v1.40", nil, nil)

	res, err := cli.ImageSearch(context.Background(), "python", types.ImageSearchOptions{
		RegistryAuth:  "",
		PrivilegeFunc: nil,
		Filters:       filters.Args{},
		Limit:         10,
	})

	logrus.Infoln("res", res)
	logrus.Errorln("err", err)

	//cli.ImagePull(context.Background(), "python")
}
