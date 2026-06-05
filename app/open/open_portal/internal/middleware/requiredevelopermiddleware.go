package middleware

import (
	"context"
	"net/http"

	"beaver/app/open/open_rpc/open"
	"beaver/app/open/open_rpc/types/open_rpc"
)

type RequireDeveloperMiddleware struct {
	openRpc open.Open
}

func NewRequireDeveloperMiddleware(openRpc open.Open) *RequireDeveloperMiddleware {
	return &RequireDeveloperMiddleware{openRpc: openRpc}
}

func (m *RequireDeveloperMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Beaver-User-Id")
		if userID == "" {
			http.Error(w, `{"code":401,"msg":"未登录"}`, http.StatusUnauthorized)
			return
		}

		res, err := m.openRpc.GetDeveloperByUserID(r.Context(), &open_rpc.GetDeveloperByUserIDReq{UserId: userID})
		if err != nil || !res.Found || res.Developer == nil || res.Developer.Status != 1 {
			http.Error(w, `{"code":403,"msg":"您还不是认证开发者,请先申请开发者资质"}`, http.StatusForbidden)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), "developerId", res.Developer.Id)))
	}
}
