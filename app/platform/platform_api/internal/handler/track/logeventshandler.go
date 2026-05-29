package track

import (
	"errors"
	"net/http"

	tracklogic "beaver/app/platform/platform_api/internal/logic/track"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func LogEventsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogEventsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if len(req.Logs) == 0 {
			response.Response(r, w, nil, errors.New("logs cannot be empty"))
			return
		}

		for i, log := range req.Logs {
			if log.Level == "" {
				logx.Errorf("log at index %d missing level", i)
				response.Response(r, w, nil, errors.New("level is required"))
				return
			}
			if log.Level != "debug" && log.Level != "info" && log.Level != "warn" && log.Level != "error" {
				logx.Errorf("log at index %d has invalid level: %s", i, log.Level)
				response.Response(r, w, nil, errors.New("level must be one of: debug, info, warn, error"))
				return
			}
			if log.BucketID == "" {
				logx.Errorf("log at index %d missing bucketId", i)
				response.Response(r, w, nil, errors.New("bucketId is required"))
				return
			}
			if log.Data == "" {
				logx.Errorf("log at index %d missing data", i)
				response.Response(r, w, nil, errors.New("data is required"))
				return
			}
			if log.Timestamp <= 0 {
				logx.Errorf("log at index %d has invalid timestamp: %d", i, log.Timestamp)
				response.Response(r, w, nil, errors.New("timestamp must be greater than 0"))
				return
			}
		}

		l := tracklogic.NewLogEventsLogic(r.Context(), svcCtx)
		resp, err := l.LogEvents(&req)
		response.Response(r, w, resp, err)
	}
}
