package logic

import (
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentsBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取动态数据（用于数据同步）
func NewGetMomentsBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentsBatchLogic {
	return &GetMomentsBatchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentsBatchLogic) GetMomentsBatch(req *types.GetMomentsBatchReq) (resp *types.GetMomentsBatchRes, err error) {
	if len(req.UserIds) == 0 {
		return &types.GetMomentsBatchRes{
			Moments:     []types.MomentBatchItem{},
			HasMore:     false,
			NextVersion: req.EndVersion,
		}, nil
	}

	var moments []moment_models.MomentModel
	query := l.svcCtx.DB.Where("user_id IN (?)", req.UserIds).
		Where("version > ? AND version <= ?", req.StartVersion, req.EndVersion).
		Order("version ASC").
		Limit(req.Limit + 1) // 多查一条来判断是否有更多数据

	err = query.Find(&moments).Error
	if err != nil {
		l.Errorf("查询动态数据失败: %v", err)
		return nil, err
	}

	// 检查是否有更多数据
	hasMore := len(moments) > req.Limit
	if hasMore {
		moments = moments[:req.Limit] // 移除多查的那条
	}

	// 转换为响应格式
	var momentItems []types.MomentBatchItem
	for _, moment := range moments {
		filesJson := ""
		if moment.Files != nil {
			if val, err := moment.Files.Value(); err == nil {
				if str, ok := val.(string); ok {
					filesJson = str
				}
			}
		}

		momentItems = append(momentItems, types.MomentBatchItem{
			UUID:      moment.UUID,
			UserID:    moment.UserID,
			Content:   moment.Content,
			Files:     filesJson,
			Version:   moment.Version,
			CreateAt:  time.Time(moment.CreatedAt).UnixMilli(),
			UpdateAt:  time.Time(moment.UpdatedAt).UnixMilli(),
			IsDeleted: moment.IsDeleted,
		})
	}

	// 计算下次同步的起始版本号
	nextVersion := req.EndVersion
	if len(moments) > 0 {
		nextVersion = moments[len(moments)-1].Version
	}

	return &types.GetMomentsBatchRes{
		Moments:     momentItems,
		HasMore:     hasMore,
		NextVersion: nextVersion,
	}, nil
}
