package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type NoticeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNoticeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NoticeUpdateLogic {
	return &NoticeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NoticeUpdateLogic) NoticeUpdate(req *types.NoticeUpdateReq) (resp *types.NoticeUpdateRes, err error) {
	var friend friend_models.FriendModel
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("不是好友关系")
	}
	if friend.SendUserID == req.UserID {
		if friend.SendUserNotice == req.Notice {
			return
		}
		// 我是发起方
		err = l.svcCtx.DB.Model(&friend).Update("send_user_notice", req.Notice).Error
	}
	if friend.RevUserID == req.UserID {
		if friend.RevUserNotice == req.Notice {
			return
		}
		// 我是接收方
		err = l.svcCtx.DB.Model(&friend).Update("rev_user_notice", req.Notice).Error
	}
	return
}
