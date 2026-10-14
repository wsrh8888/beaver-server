/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package handler

import (
	logic "beaver/app/open/open_api/internal/logic/oauth_secret"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTokenByCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTokenByCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if req.OAuthSecret == "" {
			req.OAuthSecret = r.Header.Get("OAuth-Secret")
		}
		if req.OAuthID == "" {
			req.OAuthID = r.Header.Get("OAuth-Id")
		}

		l := logic.NewGetTokenByCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetTokenByCode(&req)
		response.Response(r, w, resp, err)
	}
}
