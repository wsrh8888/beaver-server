package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/file"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BatchDeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchDeleteFileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if len(req.Ids) == 0 {
			response.Response(r, w, nil, errors.New("删除的文件ID列表不能为空"))
			return
		}

		l := logic.NewBatchDeleteFileLogic(r.Context(), svcCtx)
		resp, err := l.BatchDeleteFile(&req)
		response.Response(r, w, resp, err)
	}
}
