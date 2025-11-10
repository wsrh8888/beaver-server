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

	// 查询每个用户的最新点赞版本号（按动态发布者分组）
	for _, userId := range followeeIds {
		// 查询该用户作为动态发布者收到的点赞的最大版本号
		var maxVersion int64
		err := l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
			Where("moment_user_id = ? AND is_deleted = false", userId).
			Select("COALESCE(MAX(version), 0)").
			Scan(&maxVersion).Error

		if err != nil {
			l.Errorf("查询用户 %s 的点赞版本号失败: %v", userId, err)
			continue
		}

		// 如果版本号大于请求的since时间戳，则返回
		if maxVersion > in.Since {
			versions = append(versions, &moment.MomentLikeVersionItem{
				UserId:  userId,
				Version: maxVersion,
			})
		}
	}

	return &moment.GetMomentLikeVersionsRes{
		MomentLikeVersions: versions,
		ServerTimestamp:    time.Now().UnixMilli(),
	}, nil
}
