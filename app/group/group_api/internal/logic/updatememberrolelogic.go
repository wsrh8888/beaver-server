package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type UpdateMemberRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdateMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMemberRoleLogic {
	return &UpdateMemberRoleLogic{
		ctx:    ctx,
		logger: logger.New("update_member_role"),
		svcCtx: svcCtx,
	}
}

func (l *UpdateMemberRoleLogic) UpdateMemberRole(req *types.UpdateMemberRoleReq) (resp *types.UpdateMemberRoleRes, err error) {
	// 检查群组是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Where("group_id = ? AND status = ?", req.GroupID, 1).First(&group).Error
	if err != nil {
		logx.WithContext(l.ctx).Errorf("群组不存在或已解散，群组ID: %s", req.GroupID)
		return nil, err
	}

	// 检查操作者权限（群主）
	var operatorMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		req.GroupID, req.UserID, 1).First(&operatorMember).Error
	if err != nil {
		logx.WithContext(l.ctx).Errorf("操作者不是群成员，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
		return nil, err
	}

	// 检查操作者角色（只有群主可以修改角色）
	if operatorMember.Role != 1 {
		logx.WithContext(l.ctx).Errorf("只有群主可以修改成员角色，群组ID: %s, 用户ID: %s, 角色: %d", req.GroupID, req.UserID, operatorMember.Role)
		return nil, err
	}

	// 检查目标成员是否存在
	var targetMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ? AND status = ?",
		req.GroupID, req.MemberID, 1).First(&targetMember).Error
	if err != nil {
		logx.WithContext(l.ctx).Errorf("目标用户不是群成员，群组ID: %s, 目标用户ID: %s", req.GroupID, req.MemberID)
		return nil, err
	}

	// 不能修改自己的角色
	if req.UserID == req.MemberID {
		logx.WithContext(l.ctx).Errorf("不能修改自己的角色，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
		return nil, err
	}

	// 开始事务
	tx := l.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 获取该群成员的版本号（按群独立递增）
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
	if memberVersion == -1 {
		tx.Rollback()
		logx.WithContext(l.ctx).Errorf("获取群成员版本号失败")
		return nil, errors.New("获取版本号失败")
	}

	// 更新成员角色
	err = tx.Model(&targetMember).Update("role", req.Role).Error
	if err != nil {
		tx.Rollback()
		logx.WithContext(l.ctx).Errorf("更新成员角色失败: %v", err)
		return nil, err
	}

	// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel

	// 记录群成员变更日志
	changeLog := group_models.GroupMemberChangeLogModel{
		GroupID:    req.GroupID,
		UserID:     req.MemberID,
		ChangeType: "role_change",
		OperatedBy: req.UserID,
		ChangeTime: time.Now(),
	}
	err = tx.Create(&changeLog).Error
	if err != nil {
		tx.Rollback()
		logx.WithContext(l.ctx).Errorf("记录群成员变更日志失败: %v", err)
		return nil, err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		logx.WithContext(l.ctx).Errorf("提交事务失败: %v", err)
		return nil, err
	}

	// 异步通知群成员角色变更
	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()

		// 获取群成员列表，用于推送通知
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			logx.WithContext(l.ctx).Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 推送给所有群成员 - 群成员变动通知
		for _, member := range response.Members {
			if member.UserID != req.UserID { // 不通知操作者自己
				payload := map[string]interface{}{
					"command":  wsCommandConst.GROUP_OPERATION,
					"type":     wsTypeConst.GroupMemberReceive,
					"senderId": req.UserID,
					"targetId": member.UserID,
					"body": map[string]interface{}{
						"table": "group_members",
						"data": []map[string]interface{}{
							{
								"version": memberVersion,
								"groupId": req.GroupID,
							},
						},
					},
					"conversationId": "",
				}
				l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload)
			}
		}
	}()

	resp = &types.UpdateMemberRoleRes{
		Version: memberVersion,
	}

	logx.WithContext(l.ctx).Infof("更新群成员角色完成，群组ID: %s, 目标用户: %s, 新角色: %d, 操作者: %s", req.GroupID, req.MemberID, req.Role, req.UserID)
	l.logger.Info(model.LogMsg{
		Text: "群成员角色更新成功",
		Data: map[string]interface{}{
			"groupId":  req.GroupID,
			"userId":   req.UserID,
			"memberId": req.MemberID,
			"role":     req.Role,
		},
	})
	return resp, nil
}
