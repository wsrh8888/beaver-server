package handler

import (
	logic "beaver/app/relation/relation_api/internal/logic/group_machine"
	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGroupMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUpdateGroupMemberLogic(r.Context(), svcCtx)
		resp, err := l.UpdateGroupMember(&req)
		response.Response(r, w, resp, err)
	}
}
