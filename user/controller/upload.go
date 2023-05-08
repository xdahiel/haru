package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"haru/logs"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数解析失败！",
		})
		logs.Error("Failed resolve parameter: %v", err)
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "不支持该文件格式，目前只支持jpg、png、jpeg格式的文件！",
		})
		logs.Error("Invalid format")
		return
	}

	// 检查文件大小
	maxFileSize := int64(2 * 1024 * 1024) // 2 MB
	if file.Size > maxFileSize {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "文件过大！请上传2M以内大小的图片！",
		})
		logs.Error("File memory exceed")
		return
	}

	username := c.PostForm("username")
	if err := c.SaveUploadedFile(file, username+ext); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2002",
			"msg":  "文件上传失败！",
		})
		logs.Error("Failed save uploaded file: %v", err)
		return
	}

	getImageResults(username + ext)

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg":  "success",
	})
}

func getImageResults(filename string) []string {
	url := "http://localhost:5000/search"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open(filename)
	defer file.Close()
	part1, errFile1 := writer.CreateFormFile("image", filepath.Base(filename))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return nil
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Content-Type", "multipart/form-data")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(string(body))
	return nil
}