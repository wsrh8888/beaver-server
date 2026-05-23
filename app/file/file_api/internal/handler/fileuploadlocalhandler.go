package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"beaver/app/file/file_api/internal/handler/common"
	"beaver/app/file/file_api/internal/logic"
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"
)

func FileUploadLocalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
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
		defer file.Close()

		// 使用公共工具函数验证和处理文件
		fileReq, err := common.ValidateAndProcessFile(file, fileHead, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFileUploadLocalLogic(r.Context(), svcCtx)
		resp, _ := l.FileUploadLocal(&req)

		// 检查文件是否已经存在于数据库中
		existingFile, err := common.CheckFileExists(fileReq.FileMd5, svcCtx)
		if err == nil {
			resp.FileKey = existingFile.FileKey
			resp.OriginalName = existingFile.OriginalName

			// 生成完整URL（使用项目域名）
			if svcCtx.Config.Domain != "" {
				resp.FileURL = fmt.Sprintf("%s/api/file/preview/%s", svcCtx.Config.Domain, existingFile.FileKey)
			} else {
				// 如果未配置域名，返回相对路径
				resp.FileURL = fmt.Sprintf("/api/file/preview/%s", existingFile.FileKey)
			}

			// 转换FileInfo为API响应格式
			if existingFile.FileInfo != nil {
				resp.FileInfo = common.ConvertFileInfoToAPI(existingFile.FileInfo)
			}

			response.Response(r, w, resp, nil)
			return
		}

		// 创建本地存储目录
		uploadDir := svcCtx.Config.Local.UploadDir
		projectName := svcCtx.Config.Local.ProjectName

		// 生成本地文件路径（如果配置了项目名称，则添加项目目录前缀）
		localFilePath := common.GenerateFilePath(uploadDir, projectName, fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

		// 保存文件到本地
		if err := common.SaveFileToLocal(localFilePath, fileReq.ByteData); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 初始化文件信息
		var fileInfo *file_models.FileInfo

		// 手动解析FormData中的fileInfo字段
		if fileInfoStr := r.FormValue("fileInfo"); fileInfoStr != "" {
			var apiFileInfo types.FileInfo
			if err := json.Unmarshal([]byte(fileInfoStr), &apiFileInfo); err == nil {
				fileInfo = common.ConvertAPIFileInfoToModel(&apiFileInfo)
			}
		}

		// 生成相对路径用于数据库存储（如果配置了项目名称则包含项目目录，不包含uploadDir）
		relativePath := common.GenerateRelativePath(projectName, fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

		// 创建文件记录
		newFileModel, err := common.CreateFileRecord(fileReq, relativePath, file_models.LocalSource, svcCtx)
		if err != nil {
			// 如果数据库保存失败，删除已保存的文件
			// TODO: 可以添加删除文件的公共函数
			response.Response(r, w, nil, err)
			return
		}

		// 更新文件信息
		if fileInfo != nil {
			newFileModel.FileInfo = fileInfo
			svcCtx.DB.Save(newFileModel)
		}

		resp.FileKey = newFileModel.FileKey
		resp.OriginalName = newFileModel.OriginalName

		// 生成完整URL（使用项目域名）
		if svcCtx.Config.Domain != "" {
			resp.FileURL = fmt.Sprintf("%s/api/file/preview/%s", svcCtx.Config.Domain, newFileModel.FileKey)
		} else {
			// 如果未配置域名，返回相对路径
			resp.FileURL = fmt.Sprintf("/api/file/preview/%s", newFileModel.FileKey)
		}

		// 转换FileInfo为API响应格式
		if fileInfo != nil {
			resp.FileInfo = common.ConvertFileInfoToAPI(fileInfo)
		}

		logx.Infof("本地文件上传成功: %s", newFileModel.FileKey)
		response.Response(r, w, resp, nil)
	}
}
