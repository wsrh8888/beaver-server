package handler

import (
	"beaver/app/friend/friend_api/internal/logic"
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func friendVerifySyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendVerifySyncReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFriendVerifySyncLogic(r.Context(), svcCtx)
		resp, err := l.FriendVerifySync(&req)
		response.Response(r, w, resp, err)
	}
}
