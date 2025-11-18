package logic

import (
	"context"
	"time"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsListByIdsLogic {
	return &GetFriendsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsListByIdsLogic) GetFriendsListByIds(in *friend_rpc.GetFriendsListByIdsReq) (*friend_rpc.GetFriendsListByIdsRes, error) {
	if len(in.Ids) == 0 {
		l.Errorf("ID列表为空")
		return &friend_rpc.GetFriendsListByIdsRes{Friends: []*friend_rpc.FriendListById{}}, nil
	}

	// 查询指定ID列表中的好友信息
	var friends []friend_models.FriendModel
	query := l.svcCtx.DB.Where("uuid IN (?)", in.Ids)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&friends).Error
	if err != nil {
		l.Errorf("查询好友信息失败: ids=%v, since=%d, error=%v", in.Ids, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个好友信息", len(friends))

	// 转换为响应格式
	var friendsList []*friend_rpc.FriendListById
	for _, friend := range friends {
		friendsList = append(friendsList, &friend_rpc.FriendListById{
			Id:             friend.UUID,
			SendUserId:     friend.SendUserID,
			RevUserId:      friend.RevUserID,
			SendUserNotice: friend.SendUserNotice,
			RevUserNotice:  friend.RevUserNotice,
			Source:         friend.Source,
			IsDeleted:      friend.IsDeleted,
			Version:        friend.Version,
			CreateAt:       time.Time(friend.CreatedAt).UnixMilli(),
			UpdateAt:       time.Time(friend.UpdatedAt).UnixMilli(),
		})
	}

	return &friend_rpc.GetFriendsListByIdsRes{Friends: friendsList}, nil
}
