package handler

import (
	logic "beaver/app/relation/relation_api/internal/logic/friend_machine"
	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateFriendVerifiesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFriendVerifiesReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUpdateFriendVerifiesLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFriendVerifies(&req)
		response.Response(r, w, resp, err)
	}
}
