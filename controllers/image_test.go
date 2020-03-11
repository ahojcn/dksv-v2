package controllers

import (
	"bufio"
	"bytes"
	"dksv-v2/models"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"
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
	fileURL := "http://127.0.0.1:8080/media/psf/Home/Downloads/Python-3.7.5.tar.xz"
	filePath := "./"

	fileName := path.Base(fileURL)
	fmt.Println(fileName)

	res, err := http.Get(fileURL)
	fmt.Println(res.Request)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()
	// 获得 get 请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(filePath + fileName)
	if err != nil {
		fmt.Println(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d\n", written)
}

func TestImageController_Push(t *testing.T) {

	//UploadFile("http://tim.natapp1.cc/images/upload", map[string]string{}, "", "", )
}
func UploadFile(url string, params map[string]string, nameField, fileName string, file io.Reader) ([]byte, error) {
	HttpClient := &http.Client{
		Timeout: 3 * time.Second,
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
