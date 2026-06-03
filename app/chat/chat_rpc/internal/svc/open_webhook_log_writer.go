package svc

import (
	"context"

	"beaver/app/open/open_rpc/open"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/core/corewebhook"

	"github.com/zeromicro/go-zero/core/logx"
)

type openRpcWebhookLogWriter struct {
	openRpc open.Open
}

func newOpenRpcWebhookLogWriter(openRpc open.Open) corewebhook.LogWriter {
	return &openRpcWebhookLogWriter{openRpc: openRpc}
}

func (w *openRpcWebhookLogWriter) SaveWebhookLog(ctx context.Context, configID, appID, eventType string, success bool) {
	status := int32(0)
	if success {
		status = 1
	}
	_, err := w.openRpc.SaveWebhookLog(ctx, &open_rpc.SaveWebhookLogReq{
		ConfigId:  configID,
		AppId:     appID,
		EventType: eventType,
		Status:    status,
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("SaveWebhookLog RPC 失败: %v", err)
	}
}
