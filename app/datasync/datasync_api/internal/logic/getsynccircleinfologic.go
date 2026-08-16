package logic

import (
	"beaver/app/circle/circle_rpc/types/circle_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncCircleInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的圈子信息版本
func NewGetSyncCircleInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncCircleInfoLogic {
	return &GetSyncCircleInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncCircleInfoLogic) GetSyncCircleInfo(req *types.GetSyncCircleInfoReq) (resp *types.GetSyncCircleInfoRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	serverTimestamp := time.Now().UnixMilli()

	circleResp, err := l.svcCtx.CircleRpc.GetCircleVersions(l.ctx, &circle_rpc.GetCircleVersionsReq{
		UserId:  userId,
		Version: req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的圈子资料失败: %v", err)
		return nil, err
	}

	circleVersions := make([]types.CircleInfoVersionItem, 0)
	if circleResp.List != nil {
		for _, item := range circleResp.List {
			circleVersions = append(circleVersions, types.CircleInfoVersionItem{
				CircleID: item.CircleId,
				Version:  item.Version,
			})
		}
	}

	return &types.GetSyncCircleInfoRes{
		CircleVersions:  circleVersions,
		ServerTimestamp: serverTimestamp,
	}, nil
}
