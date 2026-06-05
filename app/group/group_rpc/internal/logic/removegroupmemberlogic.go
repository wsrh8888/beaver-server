package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RemoveGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupMemberLogic {
	return &RemoveGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveGroupMemberLogic) RemoveGroupMember(in *group_rpc.RemoveGroupMemberReq) (*group_rpc.RemoveGroupMemberRes, error) {
	var member group_models.GroupMemberModel
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = ?", in.GroupId, in.UserId, 1).
		First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("成员不在群内")
		}
		return nil, err
	}

	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", in.GroupId)
	if memberVersion == -1 {
		return nil, errors.New("获取群成员版本号失败")
	}

	memberStatus := int8(2)
	changeType := "leave"
	if in.Kick {
		memberStatus = 3
		changeType = "kick"
	}

	if err := l.svcCtx.DB.Model(&member).Updates(map[string]interface{}{
		"status":  memberStatus,
		"version": memberVersion,
	}).Error; err != nil {
		return nil, err
	}

	operatedBy := in.OperatedBy
	if operatedBy == "" {
		operatedBy = in.UserId
	}
	_ = l.svcCtx.DB.Create(&group_models.GroupMemberChangeLogModel{
		GroupID:    in.GroupId,
		UserID:     in.UserId,
		ChangeType: changeType,
		OperatedBy: operatedBy,
		ChangeTime: time.Now(),
		Version:    memberVersion,
	}).Error

	return &group_rpc.RemoveGroupMemberRes{MemberVersion: memberVersion}, nil
}
