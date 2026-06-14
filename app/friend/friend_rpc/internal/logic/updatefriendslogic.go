package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

const (
	friendActionHardDelete int32 = 1 // 物理删除好友关系
	friendActionRestore    int32 = 2 // 恢复软删除的好友关系
)

type UpdateFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendsLogic {
	return &UpdateFriendsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateFriendsLogic) UpdateFriends(in *friend_rpc.UpdateFriendsReq) (*friend_rpc.UpdateFriendsRes, error) {
	switch in.Action {
	case friendActionHardDelete:
		return l.hardDelete(in.RelationIds)
	case friendActionRestore:
		return l.restore(in.RelationIds)
	default:
		return nil, errors.New("不支持的操作类型")
	}
}

func (l *UpdateFriendsLogic) hardDelete(relationIDs []string) (*friend_rpc.UpdateFriendsRes, error) {
	var ids []uint
	for _, rid := range relationIDs {
		f, err := findFriend(l.svcCtx.DB, rid)
		if err != nil {
			continue
		}
		ids = append(ids, f.Id)
	}
	if len(ids) == 0 {
		return &friend_rpc.UpdateFriendsRes{}, nil
	}
	if err := l.svcCtx.DB.Unscoped().Where("id IN ?", ids).Delete(&friend_models.FriendModel{}).Error; err != nil {
		l.Errorf("删除好友失败: %v", err)
		return nil, err
	}
	return &friend_rpc.UpdateFriendsRes{AffectedCount: int64(len(ids))}, nil
}

func (l *UpdateFriendsLogic) restore(relationIDs []string) (*friend_rpc.UpdateFriendsRes, error) {
	var affected int64
	for _, rid := range relationIDs {
		f, err := findFriend(l.svcCtx.DB, rid)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !f.IsDeleted {
			affected++
			continue
		}
		if err := l.svcCtx.DB.Model(f).Update("is_deleted", false).Error; err != nil {
			l.Errorf("恢复好友失败: %v", err)
			return nil, err
		}
		affected++
	}
	return &friend_rpc.UpdateFriendsRes{AffectedCount: affected}, nil
}
