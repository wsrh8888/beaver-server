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


type MuteAllGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

// 全员禁言/解禁
func NewMuteAllGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteAllGroupLogic {
	return &MuteAllGroupLogic{
		ctx:    ctx,
		logger: logger.New("mute_all_group"),
		svcCtx: svcCtx,
	}
}

func (l *MuteAllGroupLogic) MuteAllGroup(req *types.MuteAllGroupReq) (resp *types.MuteAllGroupRes, err error) {
	// 验证操作者权限（群主或管理员）
	var operator group_models.GroupMemberModel
	if err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = 1", req.GroupID, req.UserID).
		First(&operator).Error; err != nil {
		return nil, errors.New("您不是该群成员")
	}
	if operator.Role != 1 && operator.Role != 2 {
		return nil, errors.New("权限不足，仅群主或管理员可操作全员禁言")
	}

	// 获取新版本号
	nextVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", req.GroupID)
	if nextVersion == -1 {
		return nil, errors.New("系统错误")
	}

	updates := map[string]any{
		"is_mute_all": req.IsMuteAll,
		"version":     nextVersion,
	}
	if req.IsMuteAll {
		now := time.Now()
		updates["mute_all_at"] = now
	} else {
		updates["mute_all_at"] = nil
	}

	if err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
		Where("group_id = ?", req.GroupID).
		Updates(updates).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("更新全员禁言失败: groupID=%s err=%v", req.GroupID, err)
		return nil, errors.New("操作失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "群全员禁言状态更新成功",
		Data: map[string]interface{}{
			"groupId":   req.GroupID,
			"userId":    req.UserID,
			"isMuteAll": req.IsMuteAll,
		},
	})

	return &types.MuteAllGroupRes{Version: nextVersion}, nil
}
