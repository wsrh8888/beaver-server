package handler

import (
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/common/response"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PreviewReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("err:", err)
			response.Response(r, w, nil, err)
			return
		}

		fileDetail, err := svcCtx.FileRpc.GetFileDetail(r.Context(), &file_rpc.GetFileDetailReq{
			FileKey: req.FileName,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				response.Response(r, w, nil, errors.New("图片不存在"))
				return
			}
			response.Response(r, w, nil, err)
			return
		}

		mac := qbox.NewMac(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK)
		deadline := time.Now().Add(time.Duration(svcCtx.Config.Qiniu.ExpireTime) * time.Second).Unix()
		privateAccessURL := storage.MakePrivateURL(mac, svcCtx.Config.Qiniu.Domain, fileDetail.Path, deadline)

		fmt.Println("privateAccessURL:", privateAccessURL)
		http.Redirect(w, r, privateAccessURL, http.StatusFound)
	}
}
