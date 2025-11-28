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

type GetMomentVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMomentVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentVersionsLogic {
	return &GetMomentVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取动态版本摘要
func (l *GetMomentVersionsLogic) GetMomentVersions(in *moment.GetMomentVersionsReq) (*moment.GetMomentVersionsRes, error) {
	// 获取用户的好友列表（包括自己）- 微信朋友圈只有好友能看到动态
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

	var versions []*moment.MomentVersionItem

	// 查询好友（包括自己）在指定时间之后有更新的动态
	var moments []moment_models.MomentModel
	err := l.svcCtx.DB.Model(&moment_models.MomentModel{}).
		Where("user_id IN (?) AND is_deleted = false AND updated_at > ?", followeeIds, time.UnixMilli(in.Since)).
		Find(&moments).Error

	if err != nil {
		l.Errorf("查询动态版本失败: %v", err)
		return &moment.GetMomentVersionsRes{
			MomentVersions:  []*moment.MomentVersionItem{},
			ServerTimestamp: time.Now().UnixMilli(),
		}, nil
	}

	// 构建版本摘要
	for _, moment := range moments {
		versions = append(versions, &moment.MomentVersionItem{
			Uuid:    moment.UUID,
			Version: moment.Version,
		})
	}

	return &moment.GetMomentVersionsRes{
		MomentVersions:  versions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
