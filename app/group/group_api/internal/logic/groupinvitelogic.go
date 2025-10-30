package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInviteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 邀请用户加入群组
func NewGroupInviteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInviteLogic {
	return &GroupInviteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInviteLogic) GroupInvite(req *types.GroupInviteReq) (resp *types.GroupInviteRes, err error) {
	// 检查群组是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Where("group_id = ? AND status = ?", req.GroupID, 1).First(&group).Error
	if err != nil {
		l.Errorf("群组不存在或已解散，群组ID: %s", req.GroupID)
		return nil, err
	}

	// 检查邀请者权限（群主或管理员）
	var inviterMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		req.GroupID, req.UserID, 1).First(&inviterMember).Error
	if err != nil {
		l.Errorf("邀请者不是群成员，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
		return nil, err
	}

	// 检查邀请者角色（群主或管理员）
	if inviterMember.Role != 1 && inviterMember.Role != 2 {
		l.Errorf("邀请者权限不足，群组ID: %s, 用户ID: %s, 角色: %d", req.GroupID, req.UserID, inviterMember.Role)
		return nil, err
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	now := time.Now()

	// 处理每个被邀请的用户
	for _, userId := range req.UserIds {
		// 检查用户是否已经是群成员
		var existingMember group_models.GroupMemberModel
		err = tx.Where("group_id = ? AND user_id = ?", req.GroupID, userId).First(&existingMember).Error
		if err == nil {
			// 用户已经是群成员，更新状态为正常
			if existingMember.Status != 1 {
				err = tx.Model(&existingMember).Update("status", 1).Error
				if err != nil {
					tx.Rollback()
					l.Errorf("更新群成员状态失败: %v", err)
					return nil, err
				}
			}
		} else {
			// 添加新群成员
			member := group_models.GroupMemberModel{
				GroupID:  req.GroupID,
				UserID:   userId,
				Role:     3, // 普通成员
				Status:   1, // 正常状态
				JoinTime: now,
				Version:  time.Now().Unix(),
			}
			err = tx.Create(&member).Error
			if err != nil {
				tx.Rollback()
				l.Errorf("添加群成员失败: %v", err)
				return nil, err
			}
		}

		// 记录群成员变更日志
		changeLog := group_models.GroupMemberChangeLogModel{
			GroupID:    req.GroupID,
			UserID:     userId,
			ChangeType: "invite",
			OperatedBy: req.UserID,
			ChangeTime: now,
			Version:    time.Now().Unix(),
		}
		err = tx.Create(&changeLog).Error
		if err != nil {
			tx.Rollback()
			l.Errorf("记录群成员变更日志失败: %v", err)
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		l.Errorf("提交事务失败: %v", err)
		return nil, err
	}

	resp = &types.GroupInviteRes{}

	l.Infof("群组邀请完成，群组ID: %s, 邀请者: %s, 被邀请用户数: %d", req.GroupID, req.UserID, len(req.UserIds))
	return resp, nil
}
