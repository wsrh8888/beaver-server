package feedback

import (
	"errors"
	"net/http"

	feedbacklogic "beaver/app/platform/platform_api/internal/logic/feedback"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SubmitFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubmitFeedbackReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if req.UserID == "" {
			response.Response(r, w, nil, errors.New("用户ID不能为空"))
			return
		}
		if req.Content == "" {
			response.Response(r, w, nil, errors.New("反馈内容不能为空"))
			return
		}
		if req.Type < 1 || req.Type > 4 {
			response.Response(r, w, nil, errors.New("反馈类型不合法"))
			return
		}

		l := feedbacklogic.NewSubmitFeedbackLogic(r.Context(), svcCtx)
		resp, err := l.SubmitFeedback(&req)
		response.Response(r, w, resp, err)
	}
}
