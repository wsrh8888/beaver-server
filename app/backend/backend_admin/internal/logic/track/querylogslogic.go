package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/track/track_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询日志
func NewQueryLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLogsLogic {
	return &QueryLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryLogsLogic) QueryLogs(req *types.QueryLogsReq) (resp *types.QueryLogsRes, err error) {
	// 构建查询条件
	db := l.svcCtx.DB.Model(&track_models.TrackLogger{}).Preload("BucketModel")

	// Bucket Id 筛选
	db = db.Where("bucket_id = ?", req.BucketID)

	// 日志级别筛选
	if req.Level != "" {
		db = db.Where("level = ?", req.Level)
	}

	// 用户ID筛选
	if req.UserFilter != "" {
		db = db.Where("user_id = ?", req.UserFilter)
	}

	// 关键词搜索（在日志数据中搜索）
	if req.Keyword != "" {
		db = db.Where("data LIKE ?", "%"+req.Keyword+"%")
	}

	// 时间范围筛选 - 直接使用时间戳
	if req.StartTime > 0 {
		db = db.Where("timestamp >= ?", req.StartTime)
	}

	if req.EndTime > 0 {
		db = db.Where("timestamp <= ?", req.EndTime)
	}

	// 查询总数
	var total int64
	if err = db.Count(&total).Error; err != nil {
		logx.Errorf("查询日志总数失败: %v", err)
		return nil, err
	}

	// 分页查询
	var logs []track_models.TrackLogger
	offset := (req.Page - 1) * req.PageSize
	if err = db.Offset(offset).Limit(req.PageSize).Order("timestamp DESC").Find(&logs).Error; err != nil {
		logx.Errorf("查询日志列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	logEntries := make([]types.QueryLogsResItem, 0, len(logs))
	for _, log := range logs {
		bucketName := ""
		if log.BucketModel != nil {
			bucketName = log.BucketModel.Name
		}

		// 将 JSON 数据转换为字符串
		dataStr := ""
		if log.Data != nil {
			dataBytes, _ := log.Data.MarshalJSON()
			dataStr = string(dataBytes)
		}

		logEntries = append(logEntries, types.QueryLogsResItem{
			Id:         uint(log.Id),
			Level:      log.Level,
			Data:       dataStr,
			BucketID:   log.BucketID,
			BucketName: bucketName,
			Timestamp:  log.Timestamp, // 直接返回时间戳字符串
			CreatedAt:  time.Time(log.CreatedAt).Format(time.RFC3339),
		})
	}

	resp = &types.QueryLogsRes{
		Total: total,
		Logs:  logEntries,
	}

	return
}
