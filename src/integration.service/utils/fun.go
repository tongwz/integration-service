package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"integration.service/http/rsp"
	"integration.service/pkg/logging"
	"integration.service/pkg/setting"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

/**
 * @note: 文件同步
 * @auth: tongWz
 * @date: 2022年5月23日13:21:54
**/
func SyncToImageServer(fileFullPath string, upkey string) (fileUrl string, fileFormat string, err error) {
	// upkey默认为"file"
	if upkey == "" {
		upkey = "file"
	}

	// 获取文件格式
	fileFormat = path.Ext(fileFullPath)

	// 打开表单
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile(upkey, filepath.Base(fileFullPath))
	if err != nil {
		return "", fileFormat, err
	}

	// 打开文件
	fh, err := os.Open(fileFullPath)
	if err != nil {
		return "", fileFormat, err
	}

	// io copy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "", fileFormat, err
	}
	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	// 同步文件到image-server-api服务器
	uploadUrl := setting.Cfg.Section("image").Key("hz_image").MustString("https://image-server-api.zq332.com") +
		setting.Cfg.Section("image").Key("attachment_upload_api").MustString("/api/attachment/upload")
	resp, err := http.Post(uploadUrl, contentType, bodyBuf)
	if err != nil {
		logging.Error("图片上传失败", err, resp)
		return "", fileFormat, err
	}

	//获取响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fileFormat, err
	}

	// 解析响应
	imageServerResp := &rsp.ImageServerResp{}
	err = json.Unmarshal(respBody, imageServerResp)
	if err != nil {
		logging.Error("图片上传返回失败", err, respBody)
		return "", "", err
	}
	if imageServerResp.Code != 0 {
		return "", fileFormat, errors.Errorf("同步文件到image-server-api服务器失败，%s,path:%s", imageServerResp.Msg, fileFullPath)
	}

	logging.Info(fmt.Sprintf("图 上 传 结 果：%+v, %s \n", imageServerResp.Msg, time.Now().Format("15:04:05")))

	// 返回文件url、格式
	fileUrl = setting.Cfg.Section("image").Key("attachment_show_api").MustString("/api/attachment/show") +
		"?id=" + fmt.Sprintf("%v", imageServerResp.Data.(map[string]interface{})["id"])

	defer func() {
		_ = resp.Body.Close()
		//关闭文件
		_ = fh.Close()
	}()

	return fileUrl, fileFormat, nil
}

/**
 * @note: 下载远程链接的文件
 * @auth: tongWz
 * @date: 2022年5月23日16:27:17
**/
func DownFileUrl(fileFullPath string, url string) error {
	rsp, err := http.Get(url)
	if err != nil {
		logging.Error("获取图片链接失败：", err)
		return err
	}
	defer rsp.Body.Close()
	out, err := os.Create(fileFullPath)
	if err != nil {
		logging.Error("文件句柄创建失败：", err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rsp.Body)
	if err != nil {
		logging.Error("句柄内容复制失败：", err)
		return err
	}
	return nil
}
