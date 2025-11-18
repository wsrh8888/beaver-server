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

func BatchDeleteFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchDeleteFriendsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if len(req.Ids) == 0 {
			response.Response(r, w, nil, errors.New("删除的好友关系ID列表不能为空"))
			return
		}

		// 验证ID格式
		for _, idStr := range req.Ids {
			_, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				response.Response(r, w, nil, errors.New("无效的好友关系ID: "+idStr))
				return
			}
		}

		l := logic.NewBatchDeleteFriendsLogic(r.Context(), svcCtx)
		resp, err := l.BatchDeleteFriends(&req)
		response.Response(r, w, resp, err)
	}
}
