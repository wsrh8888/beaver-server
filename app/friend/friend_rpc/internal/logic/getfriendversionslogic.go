package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVersionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendVersionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVersionsLogic {
	return &GetFriendVersionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendVersionsLogic) GetFriendVersions(in *friend_rpc.GetFriendVersionsReq) (*friend_rpc.GetFriendVersionsRes, error) {
	// 查询用户相关的所有好友关系（作为发送者或接收者，且未删除）
	var friends []friend_models.FriendModel
	query := l.svcCtx.DB.Where("(send_user_id = ? OR rev_user_id = ?) AND is_deleted = ?",
		in.UserId, in.UserId, false)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&friends).Error
	if err != nil {
		l.Errorf("查询用户好友版本信息失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个好友版本信息", in.UserId, len(friends))

	// 转换为响应格式
	var friendVersions []*friend_rpc.GetFriendVersionsRes_FriendVersion
	for _, friend := range friends {
		friendVersions = append(friendVersions, &friend_rpc.GetFriendVersionsRes_FriendVersion{
			Id:      friend.UUID, // 使用数据库记录的UUID作为唯一标识符
			Version: friend.Version,
		})
	}

	return &friend_rpc.GetFriendVersionsRes{FriendVersions: friendVersions}, nil
}
