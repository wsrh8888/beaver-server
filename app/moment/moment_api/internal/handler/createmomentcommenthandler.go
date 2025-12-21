package handler

import (
	"beaver/app/moment/moment_api/internal/logic"
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// CreateMomentCommentHandler 发表评论
func CreateMomentCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateMomentCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewCreateMomentCommentLogic(r.Context(), svcCtx)
		resp, err := l.CreateMomentComment(&req)
		response.Response(r, w, resp, err)
	}
}
