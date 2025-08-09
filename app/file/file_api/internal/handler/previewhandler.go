package handler

import (
	"beaver/app/file/file_api/internal/handler/common"
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PreviewReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("err:", err)
			http.Error(w, "请求参数错误", http.StatusBadRequest)
			return
		}

		// l := logic.NewPreviewLogic(r.Context(), svcCtx)
		// resp, err := l.Preview(&req)
		// response.Response(r, w, resp, err)

		var fileModel file_models.FileModel
		err := svcCtx.DB.Take(&fileModel, "file_name = ?", req.FileName).Error
		if err != nil {
			// 文件记录不存在，返回404
			http.NotFound(w, r)
			return
		}

		// 根据文件来源提供不同的预览方式
		switch fileModel.Source {
		case file_models.QiniuSource:
			// 七牛云文件预览
			filePath := fileModel.Path

			// 生成带有时效限制的签名URL
			mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
			deadline := time.Now().Add(time.Duration(svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()
			privateAccessURL := storage.MakePrivateURL(mac, svcCtx.Config.Qiniu.Domain, filePath, deadline)

			fmt.Println("七牛云文件预览URL:", privateAccessURL)
			// 返回重定向到签名URL
			http.Redirect(w, r, privateAccessURL, http.StatusFound)

		case file_models.LocalSource:
			// 本地文件预览
			localFilePath := fileModel.Path

			// 检查文件是否存在
			if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
				// 文件不存在，返回404
				http.NotFound(w, r)
				return
			}

			// 设置响应头
			w.Header().Set("Content-Type", common.GetContentType(fileModel.Type))
			w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileModel.OriginalName))

			// 直接返回文件内容
			http.ServeFile(w, r, localFilePath)

		default:
			// 不支持的文件来源，返回400错误
			http.Error(w, "不支持的文件来源", http.StatusBadRequest)
		}
	}
}
