package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"beaver/app/call/call_api/internal/logic"
	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/common/response"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/webhook"
)

func LiveKitWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	authProvider := auth.NewSimpleKeyProvider(
		svcCtx.Config.LiveKit.ApiKey,
		svcCtx.Config.LiveKit.ApiSecret,
	)

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		if _, err := webhook.ReceiveWebhookEvent(r, authProvider); err != nil {
			response.Response(r, w, nil, errors.New("webhook签名验证失败"))
			return
		}

		req := types.LiveKitWebhookReq{Body: body}
		l := logic.NewLiveKitWebhookLogic(r.Context(), svcCtx)
		resp, err := l.LiveKitWebhook(&req)
		response.Response(r, w, resp, err)
	}
}
