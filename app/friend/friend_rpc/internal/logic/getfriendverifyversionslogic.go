package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVerifyVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendVerifyVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifyVersionsLogic {
	return &GetFriendVerifyVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendVerifyVersionsLogic) GetFriendVerifyVersions(in *friend_rpc.GetFriendVerifyVersionsReq) (*friend_rpc.GetFriendVerifyVersionsRes, error) {
	// 查询用户相关的所有好友验证记录（作为发送者或接收者）
	var friendVerifies []friend_models.FriendVerifyModel
	query := l.svcCtx.DB.Where("send_user_id = ? OR rev_user_id = ?", in.UserId, in.UserId)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&friendVerifies).Error
	if err != nil {
		l.Errorf("查询用户好友验证版本信息失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个好友验证版本信息", in.UserId, len(friendVerifies))

	// 转换为响应格式
	var friendVerifyVersions []*friend_rpc.GetFriendVerifyVersionsRes_FriendVerifyVersion
	for _, verify := range friendVerifies {
		friendVerifyVersions = append(friendVerifyVersions, &friend_rpc.GetFriendVerifyVersionsRes_FriendVerifyVersion{
			Uuid:    verify.UUID,
			Version: verify.Version,
		})
	}

	return &friend_rpc.GetFriendVerifyVersionsRes{FriendVerifyVersions: friendVerifyVersions}, nil
}
