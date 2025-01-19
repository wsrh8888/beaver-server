package logic

import (
	"context"
	"fmt"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type Group_mineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroup_mineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Group_mineLogic {
	return &Group_mineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Group_mineLogic) Group_mine(req *types.GroupMineReq) (resp *types.GroupMineRes, err error) {
	// todo: add your logic here and delete this line
	var groupIdList []string
	l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).Where("user_id = ?", req.UserID).Select("group_id").Scan(&groupIdList)

	groups, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Preload: []string{"MemberList"},
		Where:   l.svcCtx.DB.Where("group_id in ?", groupIdList),
	})
	fmt.Println(groupIdList)
	resp = &types.GroupMineRes{}
	for _, model := range groups {

		resp.List = append(resp.List, types.GroupInfo{
			Title:          model.Title,
			Avatar:         model.Avatar,
			MemberCount:    len(model.MemberList),
			ConversationID: model.UUID,
		})
	}
	resp.Count = int(count)
	return
}
