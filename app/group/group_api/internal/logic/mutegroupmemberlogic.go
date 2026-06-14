package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type MuteGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 禁言/解禁群成员
func NewMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupMemberLogic {
	return &MuteGroupMemberLogic{
		ctx:    ctx,
		logger: logger.New("mute_group_member"),
		svcCtx: svcCtx,
	}
}

func (l *MuteGroupMemberLogic) MuteGroupMember(req *types.MuteGroupMemberReq) (resp *types.MuteGroupMemberRes, err error) {
	// 验证操作者权限（群主或管理员）
	var operator group_models.GroupMemberModel
	if err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", req.GroupID, req.UserID).
		First(&operator).Error; err != nil {
		return nil, errors.New("您不是该群成员")
	}
	if operator.Role != 1 && operator.Role != 2 {
		return nil, errors.New("权限不足，仅群主或管理员可禁言成员")
	}

	// 查找目标成员
	var target group_models.GroupMemberModel
	if err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", req.GroupID, req.MemberID).
		First(&target).Error; err != nil {
		return nil, errors.New("目标用户不是群成员")
	}

	// 管理员不能禁言群主，普通管理员不能禁言其他管理员
	if target.Role == 1 {
		return nil, errors.New("不能禁言群主")
	}
	if operator.Role == 2 && target.Role == 2 {
		return nil, errors.New("管理员不能禁言其他管理员")
	}

	// 计算禁言截止时间
	var mutedUntil *time.Time
	if req.Duration > 0 {
		t := time.Now().Add(time.Duration(req.Duration) * time.Minute)
		mutedUntil = &t
	}

	// 获取新版本号
	nextVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
	if nextVersion == -1 {
		return nil, errors.New("系统错误")
	}

	if err = l.svcCtx.DB.Model(&target).Updates(map[string]any{
		"muted_until": mutedUntil,
		"version":     nextVersion,
	}).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("更新禁言状态失败: groupID=%s memberID=%s err=%v", req.GroupID, req.MemberID, err)
		return nil, errors.New("操作失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "群成员禁言状态更新成功",
		Data: map[string]interface{}{
			"groupId":  req.GroupID,
			"userId":   req.UserID,
			"memberId": req.MemberID,
			"duration": req.Duration,
		},
	})

	return &types.MuteGroupMemberRes{}, nil
}
