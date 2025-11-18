package logic

import (
	"context"

	"beaver/app/track/track_api/internal/svc"
	"beaver/app/track/track_api/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LogEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogEventsLogic {
	return &LogEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 简化后的 LogEvents 方法
func (l *LogEventsLogic) LogEvents(req *types.LogEventsReq) (resp *types.LogEventsRes, err error) {
	// 准备批量插入的日志数据
	logs := make([]track_models.TrackLogger, 0, len(req.Logs))

	for _, logData := range req.Logs {
		// 解析JSON数据
		jsonData := datatypes.JSON(logData.Data)

		// 创建日志记录
		trackLog := track_models.TrackLogger{
			Level:     logData.Level,
			Data:      jsonData,
			BucketID:  logData.BucketID,
			Timestamp: logData.Timestamp,
		}

		logs = append(logs, trackLog)
	}

	// 异步批量保存日志
	go func(logsToSave []track_models.TrackLogger) {
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			result := tx.CreateInBatches(logsToSave, 100)
			if result.Error != nil {
				l.Logger.Errorf("Failed to save logs: %v", result.Error)
				return result.Error
			}
			return nil
		})

		if err != nil {
			l.Logger.Errorf("Transaction failed when saving logs: %v", err)
		}
	}(logs)

	return &types.LogEventsRes{}, nil
}
