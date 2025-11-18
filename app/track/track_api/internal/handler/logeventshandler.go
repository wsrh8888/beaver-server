package handler

import (
	"beaver/app/track/track_api/internal/logic"
	"beaver/app/track/track_api/internal/svc"
	"beaver/app/track/track_api/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func logEventsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogEventsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数验证
		if len(req.Logs) == 0 {
			response.Response(r, w, nil, errors.New("logs cannot be empty"))
			return
		}

		// 验证每个日志的必填字段
		for i, log := range req.Logs {
			if log.Level == "" {
				logx.Errorf("Log at index %d missing level", i)
				response.Response(r, w, nil, errors.New("level is required"))
				return
			}
			// 验证日志级别
			if log.Level != "debug" && log.Level != "info" && log.Level != "warn" && log.Level != "error" {
				logx.Errorf("Log at index %d has invalid level: %s", i, log.Level)
				response.Response(r, w, nil, errors.New("level must be one of: debug, info, warn, error"))
				return
			}
			if log.BucketID == "" {
				logx.Errorf("Log at index %d missing bucketId", i)
				response.Response(r, w, nil, errors.New("bucketId is required"))
				return
			}
			if log.Data == "" {
				logx.Errorf("Log at index %d missing data", i)
				response.Response(r, w, nil, errors.New("data is required"))
				return
			}
			if log.Timestamp <= 0 {
				logx.Errorf("Log at index %d has invalid timestamp: %d", i, log.Timestamp)
				response.Response(r, w, nil, errors.New("timestamp must be greater than 0"))
				return
			}
		}

		l := logic.NewLogEventsLogic(r.Context(), svcCtx)
		resp, err := l.LogEvents(&req)
		response.Response(r, w, resp, err)
	}
}
