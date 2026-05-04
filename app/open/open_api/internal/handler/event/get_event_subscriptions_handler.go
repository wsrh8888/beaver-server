// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package event

import (
	"net/http"

	"beaver/app/open/open_api/internal/logic/event"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取事件订阅列表
func GetEventSubscriptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetEventSubscriptionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := event.NewGetEventSubscriptionsLogic(r.Context(), svcCtx)
		resp, err := l.GetEventSubscriptions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
