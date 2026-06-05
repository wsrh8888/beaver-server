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

type ListFriendVerifiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFriendVerifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendVerifiesLogic {
	return &ListFriendVerifiesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListFriendVerifiesLogic) ListFriendVerifies(in *friend_rpc.ListFriendVerifiesReq) (*friend_rpc.ListFriendVerifiesRes, error) {
	if in.VerifyId != "" {
		v, err := findFriendVerify(l.svcCtx.DB, in.VerifyId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &friend_rpc.ListFriendVerifiesRes{}, nil
		}
		if err != nil {
			return nil, err
		}
		return &friend_rpc.ListFriendVerifiesRes{
			Total: 1,
			List:  []*friend_rpc.FriendVerifyItem{toFriendVerifyItem(*v)},
		}, nil
	}

	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&friend_models.FriendVerifyModel{})
	if in.SendUserId != "" {
		db = db.Where("send_user_id = ?", in.SendUserId)
	}
	if in.RevUserId != "" {
		db = db.Where("rev_user_id = ?", in.RevUserId)
	}
	if in.SendStatus > 0 {
		db = db.Where("send_status = ?", in.SendStatus)
	}
	if in.RevStatus > 0 {
		db = db.Where("rev_status = ?", in.RevStatus)
	}
	if in.StartTime != "" {
		db = db.Where("created_at >= ?", in.StartTime)
	}
	if in.EndTime != "" {
		db = db.Where("created_at <= ?", in.EndTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计好友验证失败: %v", err)
		return nil, err
	}

	var list []friend_models.FriendVerifyModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询好友验证列表失败: %v", err)
		return nil, err
	}

	items := make([]*friend_rpc.FriendVerifyItem, 0, len(list))
	for _, v := range list {
		items = append(items, toFriendVerifyItem(v))
	}
	return &friend_rpc.ListFriendVerifiesRes{Total: total, List: items}, nil
}
