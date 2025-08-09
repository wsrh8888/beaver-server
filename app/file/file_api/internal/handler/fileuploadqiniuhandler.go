package handler

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"beaver/app/file/file_api/internal/handler/common"
	"beaver/app/file/file_api/internal/logic"
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"
)

// getQiniuFileInfo 从七牛云获取文件信息
func getQiniuFileInfo(fileName string, svcCtx *svc.ServiceContext) *file_models.FileInfo {
	// 创建七牛云管理对象
	mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
	bucketManager := storage.NewBucketManager(mac, nil)

	// 获取文件信息
	fileInfo, err := bucketManager.Stat(svcCtx.Config.Qiniu.Bucket, fileName)
	if err != nil {
		logx.Errorf("获取七牛云文件信息失败: %v", err)
		return nil
	}

	logx.Infof("七牛云文件信息: MimeType=%s, Size=%d, Hash=%s", fileInfo.MimeType, fileInfo.Fsize, fileInfo.Hash)

	// 解析文件类型
	fileType := common.GetFileTypeFromMimeType(fileInfo.MimeType)

	result := &file_models.FileInfo{
		Type: file_models.FileType(fileType),
	}

	// 根据文件类型处理
	switch fileType {
	case "image":
		// 对于图片，尝试从七牛云获取尺寸信息
		if width, height := getImageSizeFromQiniu(fileName, svcCtx); width > 0 && height > 0 {
			result.ImageFile = &file_models.ImageFile{
				Width:  width,
				Height: height,
			}
		}
	case "video":
		// 对于视频，尝试从七牛云获取视频信息
		if width, height, duration := getVideoInfoFromQiniu(fileName, svcCtx); width > 0 || height > 0 {
			result.VideoFile = &file_models.VideoFile{
				Width:    width,
				Height:   height,
				Duration: duration,
			}
		}
	case "audio":
		// 对于音频，尝试从七牛云获取音频信息
		if duration := getAudioInfoFromQiniu(fileName, svcCtx); duration > 0 {
			result.AudioFile = &file_models.AudioFile{
				Duration: duration,
			}
		}
	}

	return result
}

// getImageSizeFromQiniu 从七牛云获取图片尺寸
func getImageSizeFromQiniu(fileName string, svcCtx *svc.ServiceContext) (width, height int) {
	// 使用七牛云SDK生成带签名的URL
	mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
	deadline := time.Now().Add(time.Duration(svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()

	// 生成带签名的URL，包含imageInfo查询参数
	url := storage.MakePrivateURL(mac, svcCtx.Config.Qiniu.Domain, fileName+"?imageInfo", deadline)

	// 创建带超时的HTTP客户端，跳过TLS证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	logx.Infof("正在获取图片信息，URL: %s", url)
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorf("获取图片信息失败: %v", err)
		return 0, 0
	}
	defer resp.Body.Close()

	// 读取响应内容用于调试
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("读取响应内容失败: %v", err)
		return 0, 0
	}

	logx.Infof("图片信息API响应状态码: %d, 内容: %s", resp.StatusCode, string(body))

	var result struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logx.Errorf("解析图片信息失败: %v, 响应内容: %s", err, string(body))
		return 0, 0
	}

	logx.Infof("解析到的图片尺寸: width=%d, height=%d", result.Width, result.Height)
	return result.Width, result.Height

}

// getVideoInfoFromQiniu 从七牛云获取视频信息
func getVideoInfoFromQiniu(fileName string, svcCtx *svc.ServiceContext) (width, height, duration int) {
	// 使用七牛云SDK生成带签名的URL
	mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
	deadline := time.Now().Add(time.Duration(svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()

	// 生成带签名的URL，包含avinfo查询参数
	url := storage.MakePrivateURL(mac, svcCtx.Config.Qiniu.Domain, fileName+"?avinfo", deadline)

	// 创建带超时的HTTP客户端，跳过TLS证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	logx.Infof("正在获取视频信息，URL: %s", url)
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorf("获取视频信息失败: %v", err)
		return 0, 0, 0
	}
	defer resp.Body.Close()

	// 读取响应内容用于调试
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("读取视频信息响应失败: %v", err)
		return 0, 0, 0
	}

	logx.Infof("视频信息API响应状态码: %d, 内容: %s", resp.StatusCode, string(body))

	var result struct {
		Streams []struct {
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration string `json:"duration"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logx.Errorf("解析视频信息失败: %v, 响应内容: %s", err, string(body))
		return 0, 0, 0
	}

	if len(result.Streams) > 0 {
		stream := result.Streams[0]
		// 解析时长字符串为秒数
		duration := common.ParseDuration(stream.Duration)
		logx.Infof("解析到的视频信息: width=%d, height=%d, duration=%d", stream.Width, stream.Height, duration)
		return stream.Width, stream.Height, duration
	}

	logx.Errorf("视频信息中没有找到streams数据")
	return 0, 0, 0
}

// getAudioInfoFromQiniu 从七牛云获取音频信息
func getAudioInfoFromQiniu(fileName string, svcCtx *svc.ServiceContext) (duration int) {
	// 使用七牛云SDK生成带签名的URL
	mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
	deadline := time.Now().Add(time.Duration(svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()

	// 生成带签名的URL，包含avinfo查询参数
	url := storage.MakePrivateURL(mac, svcCtx.Config.Qiniu.Domain, fileName+"?avinfo", deadline)

	// 创建带超时的HTTP客户端，跳过TLS证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	logx.Infof("正在获取音频信息，URL: %s", url)
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorf("获取音频信息失败: %v", err)
		return 0
	}
	defer resp.Body.Close()

	// 读取响应内容用于调试
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("读取音频信息响应失败: %v", err)
		return 0
	}

	logx.Infof("音频信息API响应状态码: %d, 内容: %s", resp.StatusCode, string(body))

	var result struct {
		Streams []struct {
			Duration string `json:"duration"`
		} `json:"streams"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logx.Errorf("解析音频信息失败: %v, 响应内容: %s", err, string(body))
		return 0
	}

	if len(result.Streams) > 0 {
		duration := common.ParseDuration(result.Streams[0].Duration)
		logx.Infof("解析到的音频时长: %d秒", duration)
		return duration
	}

	logx.Errorf("音频信息中没有找到streams数据")
	return 0
}

func FileUploadQiniuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		file, fileHead, err := r.FormFile("file")
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		// 使用公共工具函数验证和处理文件
		fileReq, err := common.ValidateAndProcessFile(file, fileHead, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFileUploadQiniuLogic(r.Context(), svcCtx)
		resp, _ := l.FileUploadQiniu(&req)

		// 检查文件是否已经存在于数据库中
		existingFile, err := common.CheckFileExists(fileReq.FileMd5, svcCtx)
		if err == nil {
			resp.FileName = existingFile.FileName
			resp.OriginalName = existingFile.OriginalName

			// 如果文件已存在但FileInfo为空，同步从七牛云获取
			if existingFile.FileInfo == nil {
				fileInfo := getQiniuFileInfo(existingFile.Path, svcCtx)
				if fileInfo != nil {
					existingFile.FileInfo = fileInfo
					svcCtx.DB.Save(existingFile)
					logx.Infof("已存在文件元数据获取成功: %s", existingFile.FileName)
				}
			}

			// 转换FileInfo为API响应格式
			if existingFile.FileInfo != nil {
				resp.FileInfo = common.ConvertFileInfoToAPI(existingFile.FileInfo)
			}

			response.Response(r, w, resp, nil)
			return
		}

		// 根据文件类型创建目录结构，并生成七牛云文件路径
		fileMd5Name := fileReq.FileMd5 + "." + fileReq.Suffix
		qiniuFilePath := fmt.Sprintf("%s/%s", fileReq.FileType, fileMd5Name)

		// 上传文件到七牛云
		qiniuURL, err := uploadToQiniu(qiniuFilePath, fileReq.ByteData, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 创建文件记录
		newFileModel, err := common.CreateFileRecord(fileReq, qiniuURL, file_models.QiniuSource, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 同步从七牛云获取文件详细信息
		fileInfo := getQiniuFileInfo(qiniuFilePath, svcCtx)
		if fileInfo != nil {
			// 更新数据库中的FileInfo
			newFileModel.FileInfo = fileInfo
			svcCtx.DB.Save(newFileModel)
			logx.Infof("文件元数据获取成功: %s", newFileModel.FileName)
		}

		resp.FileName = newFileModel.FileName
		resp.OriginalName = newFileModel.OriginalName

		// 转换FileInfo为API响应格式
		if fileInfo != nil {
			resp.FileInfo = common.ConvertFileInfoToAPI(fileInfo)
		}

		response.Response(r, w, resp, nil)
	}
}

func uploadToQiniu(filePath string, fileData []byte, config *svc.ServiceContext) (string, error) {
	// 设置认证信息
	mac := credentials.NewCredentials(config.Config.Qiniu.AK, config.Config.Qiniu.SK)

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})

	reader := bytes.NewReader(fileData)
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: config.Config.Qiniu.Bucket,
		FileName:   filePath,
		ObjectName: &filePath,
	}, nil)

	if err != nil {
		return "", fmt.Errorf("failed to upload file to Qiniu: %v", err)
	}

	return filePath, nil
}
