package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	utils "beaver/utils/rand"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupCreateLogic) GroupCreate(req *types.GroupCreateReq) (resp *types.GroupCreateRes, err error) {

	var groupModel = group_models.GroupModel{
		Creator:  req.UserID,
		UUID:     utils.GenerateUUId(),
		Abstract: "本群创建于" + time.Now().Format("2006-01-02") + "，欢迎大家加入",
		Size:     50,
	}
	var groupUserList = []string{string(req.UserID)}
	if len(req.UserIdList) == 0 {
		return nil, errors.New("请选择用户")
	}

	for _, u := range req.UserIdList {
		groupUserList = append(groupUserList, u)
	}

	groupModel.Title = req.Name

	err = l.svcCtx.DB.Create(&groupModel).Error
	if err != nil {
		logx.Errorf("创建群失败: %v", err)
		return nil, errors.New("创建群失败")
	}

	var members []group_models.GroupMemberModel
	for i, u := range groupUserList {

		memberMode := group_models.GroupMemberModel{
			GroupID: groupModel.UUID,
			UserID:  u,
			Role:    3,
		}
		if i == 0 {
			memberMode.Role = 1
		}
		members = append(members, memberMode)
	}

	err = l.svcCtx.DB.Create(&members).Error
	if err != nil {
		logx.Errorf("创建群成员失败: %v", err)
		return nil, errors.New("创建群成员失败")
	}
	return &types.GroupCreateRes{
		GroupID: groupModel.UUID,
	}, nil
}
