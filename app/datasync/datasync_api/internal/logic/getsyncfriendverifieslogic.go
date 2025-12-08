package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncFriendVerifiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的好友验证版本
func NewGetSyncFriendVerifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncFriendVerifiesLogic {
	return &GetSyncFriendVerifiesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncFriendVerifiesLogic) GetSyncFriendVerifies(req *types.GetSyncFriendVerifiesReq) (resp *types.GetSyncFriendVerifiesRes, err error) {
	// 调用Friend RPC获取好友验证版本信息
	verifyResp, err := l.svcCtx.FriendRpc.GetFriendVerifyVersions(l.ctx, &friend_rpc.GetFriendVerifyVersionsReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取好友验证版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个好友验证版本信息", len(verifyResp.FriendVerifyVersions))

	// 转换为响应格式，确保返回空数组而不是null
	friendVerifyVersions := make([]types.FriendVerifyVersionItem, 0)
	if verifyResp.FriendVerifyVersions != nil {
		for _, verify := range verifyResp.FriendVerifyVersions {
			friendVerifyVersions = append(friendVerifyVersions, types.FriendVerifyVersionItem{
				VerifyId: verify.VerifyId,
				Version:  verify.Version,
			})
		}
	}

	return &types.GetSyncFriendVerifiesRes{
		FriendVerifyVersions: friendVerifyVersions,
		ServerTimestamp:      time.Now().UnixMilli(),
	}, nil
}
