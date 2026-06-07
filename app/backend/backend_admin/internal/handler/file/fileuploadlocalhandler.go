package handler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	filecommon "beaver/app/backend/backend_admin/internal/handler/file/common"
	logic "beaver/app/backend/backend_admin/internal/logic/file"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"

	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadLocalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileUploadLocalReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 从表单获取文件
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		defer file.Close()

		// 获取fileInfo
		fileInfoStr := r.FormValue("fileInfo")

		// 使用公共工具函数验证和处理文件
		fileReq, err := filecommon.ValidateAndProcessFile(file, fileHeader, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 确定文件来源
		source := file_models.LocalSource
		if req.Source != "" && req.Source == "qiniu" {
			source = file_models.QiniuSource
		}

		l := logic.NewFileUploadLocalLogic(r.Context(), svcCtx)
		resp, err := l.FileUploadLocal(&req)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 检查文件是否已经存在于数据库中
		existingFile, err := checkFileExists(fileReq.FileMd5, svcCtx)
		if err == nil {
			resp.FileKey = existingFile.FileKey
			resp.OriginalName = existingFile.OriginalName
			if svcCtx.Config.Domain == "" {
				response.Response(r, w, nil, errors.New("未配置域名"))
				return
			}
			resp.FileURL = filecommon.BuildLocalFileURL(svcCtx.Config.Domain, existingFile.FileKey)
			response.Response(r, w, resp, nil)
			return
		}

		// 创建本地存储目录
		uploadDir := svcCtx.Config.Local.UploadDir
		if uploadDir == "" {
			uploadDir = "./uploads" // 默认目录
		}

		// 获取项目名称（如果为空则不加项目目录前缀）
		projectName := svcCtx.Config.Local.ProjectName

		// 生成本地文件路径（如果配置了项目名称，则添加项目目录前缀）
		localFilePath := generateFilePath(uploadDir, projectName, fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

		// 保存文件到本地
		if err := saveFileToLocal(localFilePath, fileReq.ByteData); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 生成相对路径用于数据库存储（如果配置了项目名称则包含项目目录，不包含uploadDir）
		relativePath := generateRelativePath(projectName, fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

		// 保存文件信息到数据库
		saveReq := &types.SaveFileReq{
			OriginalName: fileReq.OriginalName,
			Size:         fileReq.Size,
			Path:         relativePath,
			Md5:          fileReq.FileMd5,
			Type:         fileReq.FileType,
			Source:       string(source),
			FileInfo:     fileInfoStr,
		}

		saveLogic := logic.NewSaveFileLogic(r.Context(), svcCtx)
		saveResp, err := saveLogic.SaveFile(saveReq)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		resp.FileKey = saveResp.FileKey
		resp.OriginalName = fileReq.OriginalName
		resp.FileURL = filecommon.BuildLocalFileURL(svcCtx.Config.Domain, saveResp.FileKey)

		logx.Infof("本地文件上传成功: %s, url: %s", saveResp.FileKey, resp.FileURL)
		response.Response(r, w, resp, nil)
	}
}

// checkFileExists 检查文件是否已存在于数据库中
func checkFileExists(fileMd5 string, svcCtx *svc.ServiceContext) (*file_models.FileModel, error) {
	var fileModel file_models.FileModel
	err := svcCtx.DB.Take(&fileModel, "md5 = ?", fileMd5).Error
	if err != nil {
		return nil, err
	}
	return &fileModel, nil
}

// saveFileToLocal 保存文件到本地
func saveFileToLocal(filePath string, data []byte) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	logx.Infof("文件保存成功: %s", filePath)
	return nil
}

// generateFilePath 生成文件路径
// projectName: 项目名称，如果为空则不加项目目录前缀
func generateFilePath(uploadDir, projectName, fileType, fileMd5, suffix string) string {
	fileMd5Name := fileMd5 + "." + suffix
	if projectName != "" {
		return filepath.Join(uploadDir, projectName, fileType, fileMd5Name)
	}
	return filepath.Join(uploadDir, fileType, fileMd5Name)
}

// generateRelativePath 生成相对路径（不包含uploadDir，用于数据库存储）
// projectName: 项目名称，如果为空则不加项目目录前缀
func generateRelativePath(projectName, fileType, fileMd5, suffix string) string {
	fileMd5Name := fileMd5 + "." + suffix
	var path string
	if projectName != "" {
		path = filepath.Join(projectName, fileType, fileMd5Name)
	} else {
		path = filepath.Join(fileType, fileMd5Name)
	}
	return filepath.ToSlash(path)
}
