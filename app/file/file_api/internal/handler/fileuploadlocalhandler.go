package handler

import (
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
			resp.FileName = existingFile.FileName
			resp.OriginalName = existingFile.OriginalName

			// 转换FileInfo为API响应格式
			if existingFile.FileInfo != nil {
				resp.FileInfo = common.ConvertFileInfoToAPI(existingFile.FileInfo)
			}

			response.Response(r, w, resp, nil)
			return
		}

		// 创建本地存储目录
		uploadDir := svcCtx.Config.Local.UploadDir

		// 生成本地文件路径
		localFilePath := common.GenerateFilePath(uploadDir, fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

		// 保存文件到本地
		if err := common.SaveFileToLocal(localFilePath, fileReq.ByteData); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 获取文件信息
		fileInfo := common.GetLocalFileInfo(localFilePath, fileReq.FileType)

		// 生成相对路径用于数据库存储（不包含uploadDir）
		relativePath := common.GenerateRelativePath(fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

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

		resp.FileName = newFileModel.FileName
		resp.OriginalName = newFileModel.OriginalName

		// 转换FileInfo为API响应格式
		if fileInfo != nil {
			resp.FileInfo = common.ConvertFileInfoToAPI(fileInfo)
		}

		logx.Infof("本地文件上传成功: %s", newFileModel.FileName)
		response.Response(r, w, resp, nil)
	}
}
