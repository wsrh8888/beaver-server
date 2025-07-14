package logic

import (
	"context"
	"errors"

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
	// 参数验证
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 查询好友验证列表
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
		} else if fv.RevUserID == req.UserID {
			// 我是接收方
			info.UserID = fv.SendUserID
			info.Nickname = fv.SendUserModel.NickName
			info.Avatar = fv.SendUserModel.Avatar
			info.Flag = "receive"
			info.Status = fv.RevStatus
		} else {
			// 这种情况理论上不应该发生，跳过
			continue
		}

		list = append(list, info)
	}

	l.Logger.Infof("获取好友验证列表成功: userID=%s, count=%d", req.UserID, len(list))
	return &types.ValidListRes{
		Count: count,
		List:  list,
	}, nil
}
