package handler

import (
	logic "beaver/app/relation/relation_api/internal/logic/friend_machine"
	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FriendIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFriendIdsLogic(r.Context(), svcCtx)
		resp, err := l.FriendIds(&req)
		response.Response(r, w, resp, err)
	}
}
