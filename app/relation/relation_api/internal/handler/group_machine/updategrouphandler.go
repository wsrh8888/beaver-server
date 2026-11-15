package handler

import (
	logic "beaver/app/relation/relation_api/internal/logic/group_machine"
	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUpdateGroupLogic(r.Context(), svcCtx)
		resp, err := l.UpdateGroup(&req)
		response.Response(r, w, resp, err)
	}
}
