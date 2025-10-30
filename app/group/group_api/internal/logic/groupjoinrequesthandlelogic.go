package logic

import (
	"context"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupJoinRequestHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 处理群组申请
func NewGroupJoinRequestHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestHandleLogic {
	return &GroupJoinRequestHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupJoinRequestHandleLogic) GroupJoinRequestHandle(req *types.GroupJoinRequestHandleReq) (resp *types.GroupJoinRequestHandleRes, err error) {
	// 查询申请记录
	var request group_models.GroupJoinRequestModel
	err = l.svcCtx.DB.Where("id = ?", req.RequestID).First(&request).Error
	if err != nil {
		l.Errorf("查询群组申请记录失败: %v", err)
		return nil, err
	}

	// 检查申请状态
	if request.Status != 0 {
		l.Errorf("申请已被处理，申请ID: %d, 当前状态: %d", req.RequestID, request.Status)
		return nil, err
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新申请状态
	now := time.Now()
	err = tx.Model(&request).Updates(map[string]interface{}{
		"status":     req.Status,
		"handled_by": req.UserID,
		"handled_at": &now,
	}).Error
	if err != nil {
		tx.Rollback()
		l.Errorf("更新申请状态失败: %v", err)
		return nil, err
	}

	// 如果同意申请，添加群成员
	if req.Status == 1 {
		// 检查用户是否已经是群成员
		var existingMember group_models.GroupMemberModel
		err = tx.Where("group_id = ? AND user_id = ?", request.GroupID, request.ApplicantUserID).First(&existingMember).Error
		if err == nil {
			// 用户已经是群成员，更新状态为正常
			err = tx.Model(&existingMember).Update("status", 1).Error
			if err != nil {
				tx.Rollback()
				l.Errorf("更新群成员状态失败: %v", err)
				return nil, err
			}
		} else {
			// 添加新群成员
			member := group_models.GroupMemberModel{
				GroupID:  request.GroupID,
				UserID:   request.ApplicantUserID,
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
			GroupID:    request.GroupID,
			UserID:     request.ApplicantUserID,
			ChangeType: "join",
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

	resp = &types.GroupJoinRequestHandleRes{}

	statusText := "拒绝"
	if req.Status == 1 {
		statusText = "同意"
	}

	l.Infof("处理群组申请完成，申请ID: %d, 处理结果: %s, 处理者: %s", req.RequestID, statusText, req.UserID)
	return resp, nil
}
