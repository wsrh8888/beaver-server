package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/friend"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFriendDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.FriendID == "" {
			response.Response(r, w, nil, errors.New("好友关系ID不能为空"))
			return
		}

		// 验证ID格式
		_, err := strconv.ParseUint(req.FriendID, 10, 32)
		if err != nil {
			response.Response(r, w, nil, errors.New("无效的好友关系ID"))
			return
		}

		l := logic.NewGetFriendDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetFriendDetail(&req)
		response.Response(r, w, resp, err)
	}
}
