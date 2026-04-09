package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlockListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取黑名单列表
func NewBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockListLogic {
	return &BlockListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlockListLogic) BlockList(req *types.BlockListReq) (resp *types.BlockListRes, err error) {
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	var blocks []friend_models.FriendBlockModel
	var count int64
	query := l.svcCtx.DB.Model(&friend_models.FriendBlockModel{}).Where("user_id = ?", req.UserID)
	if err = query.Count(&count).Error; err != nil {
		return nil, errors.New("查询失败")
	}
	if err = query.Offset(offset).Limit(limit).Find(&blocks).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	list := make([]types.BlockUserInfo, 0, len(blocks))
	for _, b := range blocks {
		info := types.BlockUserInfo{
			BlockID: b.BlockID,
			UserID:  b.BlockedUserID,
		}
		userResp, uErr := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: b.BlockedUserID})
		if uErr == nil && userResp.UserInfo != nil {
			info.NickName = userResp.UserInfo.NickName
			info.Avatar = userResp.UserInfo.Avatar
		}
		list = append(list, info)
	}

	return &types.BlockListRes{List: list, Count: count}, nil
}
