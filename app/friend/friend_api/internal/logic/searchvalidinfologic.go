package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SearchValidInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchValidInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchValidInfoLogic {
	return &SearchValidInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchValidInfoLogic) SearchValidInfo(req *types.SearchValidInfoReq) (resp *types.SearchValidInfoRes, err error) {
	var friendVerify friend_models.FriendVerifyModel
	// 操作状态，当前用户为接受方
	err = l.svcCtx.DB.Where(
		"(rev_user_id = ? and send_user_id = ?) or (rev_user_id = ? and send_user_id = ?)",
		req.UserId, req.FriendId, req.FriendId, req.UserId,
	).First(&friendVerify).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("好友验证不存在")
		}
		return nil, err
	}

	if friendVerify.RevStatus != 0 {
		return nil, errors.New("操作异常")
	}

	// 填充返回结果
	resp = &types.SearchValidInfoRes{
		ValidId: friendVerify.Id,
	}

	return resp, nil
}
