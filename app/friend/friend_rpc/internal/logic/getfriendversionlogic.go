package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVersionLogic {
	return &GetFriendVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendVersionLogic) GetFriendVersion(in *friend_rpc.GetFriendVersionReq) (*friend_rpc.GetFriendVersionRes, error) {
	var maxVersion int64
	err := l.svcCtx.DB.Model(&friend_models.FriendModel{}).
		Where("(send_user_id = ? OR rev_user_id = ?)", in.UserId, in.UserId).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error

	if err != nil {
		l.Errorf("获取最新好友版本号失败: %v", err)
		return nil, err
	}

	return &friend_rpc.GetFriendVersionRes{
		LatestVersion: maxVersion,
	}, nil
}
