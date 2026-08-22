package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleMembersLogic {
	return &GetCircleMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleMembersLogic) GetCircleMembers(req *types.GetCircleMembersReq) (resp *types.GetCircleMembersRes, err error) {
	var self circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&self).Error; err != nil {
		return nil, fmt.Errorf("无权限查看成员")
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	var total int64
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Model(&circle_models.CircleMemberModel{}).
		Where("circle_id = ?", req.CircleID).
		Count(&total)
	l.svcCtx.DB.Where("circle_id = ?", req.CircleID).
		Order("role ASC, id ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&members)

	if len(members) == 0 {
		return &types.GetCircleMembersRes{Count: total, List: []types.GetCircleMembersItem{}}, nil
	}

	userIDs := make([]string, 0, len(members))
	for _, m := range members {
		userIDs = append(userIDs, m.UserID)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	items := make([]types.GetCircleMembersItem, 0, len(members))
	for _, m := range members {
		item := types.GetCircleMembersItem{
			UserID: m.UserID,
			Role:   m.Role,
		}
		if userResp != nil {
			if info := userResp.UserInfo[m.UserID]; info != nil {
				item.UserName = info.NickName
				item.Avatar = info.Avatar
			}
		}
		items = append(items, item)
	}

	return &types.GetCircleMembersRes{Count: total, List: items}, nil
}
