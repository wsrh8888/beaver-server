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
	"net/http"

	logic "beaver/app/open/open_portal/internal/logic/digital"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DetailDigitalKnowledgeSourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DetailDigitalKnowledgeSourceReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		l := logic.NewDetailDigitalKnowledgeSourceLogic(r.Context(), svcCtx)
		resp, err := l.DetailDigitalKnowledgeSource(&req)
		response.Response(r, w, resp, err)
	}
}
