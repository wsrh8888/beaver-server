package handler

import (
	"beaver/app/friend/friend_api/internal/logic"
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/common/response"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func addFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddFriendReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数验证
		if req.UserID == "" || req.FriendID == "" {
			response.Response(r, w, nil, errors.New("用户ID和好友ID不能为空"))
			return
		}

		// 验证来源字段
		if req.Source == "" {
			response.Response(r, w, nil, errors.New("来源字段不能为空"))
			return
		}

		// 验证来源值是否合法
		validSources := map[string]bool{
			"email":  true,
			"qrcode": true,
		}
		if !validSources[req.Source] {
			response.Response(r, w, nil, errors.New("无效的来源值，只支持email和qrcode"))
			return
		}

		// 不能添加自己为好友
		if req.UserID == req.FriendID {
			response.Response(r, w, nil, errors.New("不能添加自己为好友"))
			return
		}

		fmt.Println("2222222222222222222")
		fmt.Println(r)

		l := logic.NewAddFriendLogic(r.Context(), svcCtx)
		resp, err := l.AddFriend(&req)
		response.Response(r, w, resp, err, "发送成功")
	}
}
