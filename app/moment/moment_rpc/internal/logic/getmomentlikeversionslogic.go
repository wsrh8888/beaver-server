package logic

import (
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/moment"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentLikeVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMomentLikeVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentLikeVersionsLogic {
	return &GetMomentLikeVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取动态点赞版本摘要
func (l *GetMomentLikeVersionsLogic) GetMomentLikeVersions(in *moment.GetMomentLikeVersionsReq) (*moment.GetMomentLikeVersionsRes, error) {
	// 获取用户的好友列表（包括自己）- 只有好友的动态点赞才能看到
	followeeIds := []string{in.UserId} // 先包含自己

	// 调用好友服务获取好友列表
	if l.svcCtx.FriendRpc != nil {
		friendReq := &friend_rpc.GetFriendIdsRequest{UserID: in.UserId}
		friendResp, err := l.svcCtx.FriendRpc.GetFriendIds(l.ctx, friendReq)
		if err == nil && len(friendResp.FriendIds) > 0 {
			followeeIds = append(followeeIds, friendResp.FriendIds...)
		} else {
			l.Errorf("获取好友列表失败: %v", err)
		}
	}

	var versions []*moment.MomentLikeVersionItem

	// 查询好友（包括自己）的动态在指定时间之后收到的点赞
	var likes []moment_models.MomentLikeModel
	err := l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
		Where("moment_user_id IN (?) AND is_deleted = false AND updated_at > ?", followeeIds, time.UnixMilli(in.Since)).
		Find(&likes).Error

	if err != nil {
		l.Errorf("查询点赞版本失败: %v", err)
		return &moment.GetMomentLikeVersionsRes{
			MomentLikeVersions: []*moment.MomentLikeVersionItem{},
			ServerTimestamp:    time.Now().UnixMilli(),
		}, nil
	}

	// 构建版本摘要
	for _, like := range likes {
		versions = append(versions, &moment.MomentLikeVersionItem{
			Uuid:    like.UUID,
			Version: like.Version,
		})
	}

	return &moment.GetMomentLikeVersionsRes{
		MomentLikeVersions: versions,
		ServerTimestamp:    time.Now().UnixMilli(),
	}, nil
}
