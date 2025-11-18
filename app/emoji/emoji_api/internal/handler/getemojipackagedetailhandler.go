package handler

import (
	"beaver/app/emoji/emoji_api/internal/logic"
	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetEmojiPackageDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetEmojiPackageDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetEmojiPackageDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetEmojiPackageDetail(&req)
		response.Response(r, w, resp, err)
	}
}
