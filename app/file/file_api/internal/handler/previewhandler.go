package handler

import (
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"
	"errors"
	"fmt"
	"net/http"
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
			response.Response(r, w, nil, err)
			return
		}

		// l := logic.NewPreviewLogic(r.Context(), svcCtx)
		// resp, err := l.Preview(&req)
		// response.Response(r, w, resp, err)

		var fileModel file_models.FileModel
		err := svcCtx.DB.Take(&fileModel, "file_id = ?", req.FileID).Error
		if err != nil {
			response.Response(r, w, nil, errors.New("图片不存在"))
			return
		}

		// 构建七牛云文件URL
		filePath := fileModel.Path

		// 生成带有时效限制的签名URL
		mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
		deadline := time.Now().Add(time.Duration(svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()
		privateAccessURL := storage.MakePrivateURL(mac, svcCtx.Config.Qiniu.Domain, filePath, deadline)

		fmt.Println("privateAccessURL:", privateAccessURL)
		// 返回重定向到签名URL
		http.Redirect(w, r, privateAccessURL, http.StatusFound)
	}
}
