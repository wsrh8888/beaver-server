package track_public

import (
	"context"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"

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

func (l *LogEventsLogic) LogEvents(req *types.LogEventsReq) (*types.LogEventsRes, error) {
	logs := make([]platform_models.TrackLogger, 0, len(req.Logs))
	for _, logData := range req.Logs {
		logs = append(logs, platform_models.TrackLogger{
			Level:     logData.Level,
			Data:      datatypes.JSON(logData.Data),
			BucketID:  logData.BucketID,
			Timestamp: logData.Timestamp,
		})
	}

	go func(logsToSave []platform_models.TrackLogger) {
		err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			return tx.CreateInBatches(logsToSave, 100).Error
		})
		if err != nil {
			l.Logger.Errorf("save logs failed: %v", err)
		}
	}(logs)

	return &types.LogEventsRes{}, nil
}
