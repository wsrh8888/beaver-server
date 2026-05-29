package handler

import (
	"errors"
	"net/http"

	"beaver/app/document/document_api/internal/logic"
	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDocumentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDocumentReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if req.UserID == "" {
			response.Response(r, w, nil, errors.New("用户ID不能为空"))
			return
		}
		if req.DocID == "" {
			response.Response(r, w, nil, errors.New("文档标识不能为空"))
			return
		}

		l := logic.NewGetDocumentLogic(r.Context(), svcCtx)
		resp, err := l.GetDocument(&req)
		response.Response(r, w, resp, err)
	}
}
