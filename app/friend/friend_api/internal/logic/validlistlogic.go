package logic

import (
	"context"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidListLogic {
	return &ValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidListLogic) ValidList(req *types.ValidListReq) (resp *types.ValidListRes, err error) {

	fvs, count, _ := list_query.ListQuery(l.svcCtx.DB, friend_models.FriendVerifyModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where:   l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID),
		Preload: []string{"RevUserModel", "SendUserModel"},
	})

	var list []types.FriendValidInfo

	for _, fv := range fvs {
		info := types.FriendValidInfo{
			Message: fv.Message,
			ID:      fv.ID,
		}
		if fv.SendUserID == req.UserID {
			// 我是发起方
			info.UserID = fv.RevUserID
			info.Nickname = fv.RevUserModel.NickName
			info.Avatar = fv.RevUserModel.Avatar
			info.Flag = "send"
			info.Status = fv.RevStatus
		}
		if fv.RevUserID == req.UserID {
			// 我是接受方
			info.UserID = fv.SendUserID
			info.Nickname = fv.SendUserModel.NickName
			info.Avatar = fv.SendUserModel.Avatar
			info.Flag = "rev"
			info.Status = fv.RevStatus
		}

		list = append(list, info)
	}

	return &types.ValidListRes{
		Count: count,
		List:  list,
	}, nil
}
