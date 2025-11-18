package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/user"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BatchDeleteUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchDeleteUsersReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if len(req.Ids) == 0 {
			response.Response(r, w, nil, errors.New("请选择要删除的用户"))
			return
		}

		l := logic.NewBatchDeleteUsersLogic(r.Context(), svcCtx)
		resp, err := l.BatchDeleteUsers(&req)
		response.Response(r, w, resp, err)
	}
}
