package handler

import (
	"beaver/app/friend/friend_api/internal/logic"
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func friendInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数验证
		if req.FriendID == "" {
			response.Response(r, w, nil, errors.New("好友ID不能为空"))
			return
		}

		// 不能查询自己的信息
		if req.UserID == req.FriendID {
			response.Response(r, w, nil, errors.New("不能查询自己的信息"))
			return
		}

		l := logic.NewFriendInfoLogic(r.Context(), svcCtx)
		resp, err := l.FriendInfo(&req)
		response.Response(r, w, resp, err)
	}
}
