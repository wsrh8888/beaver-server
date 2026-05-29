package handler

import (
	"errors"
	"net/http"
	"unicode/utf8"

	"beaver/app/document/document_api/internal/logic"
	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/common/models/ctype"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateDocumentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateDocumentReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if req.UserID == "" {
			response.Response(r, w, nil, errors.New("用户ID不能为空"))
			return
		}
		if req.DocType != nil {
			if *req.DocType < ctype.CloudDocTypeFolder || *req.DocType > ctype.CloudDocTypeMind {
				response.Response(r, w, nil, errors.New("文档类型不合法"))
				return
			}
		}
		if req.Title != "" && utf8.RuneCountInString(req.Title) > 256 {
			response.Response(r, w, nil, errors.New("名称过长"))
			return
		}
		if req.DocType != nil && *req.DocType == ctype.CloudDocTypeFolder && req.FileKey != "" {
			response.Response(r, w, nil, errors.New("文件夹不能包含正文"))
			return
		}

		l := logic.NewCreateDocumentLogic(r.Context(), svcCtx)
		resp, err := l.CreateDocument(&req)
		response.Response(r, w, resp, err)
	}
}
