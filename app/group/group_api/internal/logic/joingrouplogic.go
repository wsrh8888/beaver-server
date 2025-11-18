package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinGroupLogic) JoinGroup(req *types.GroupJoinReq) (resp *types.GroupJoinRes, err error) {
	var memberVersion int64 = 0

	// 检查群组是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Where("group_id = ? AND status = ?", req.GroupID, 1).First(&group).Error
	if err != nil {
		l.Errorf("群组不存在或已解散，群组ID: %s", req.GroupID)
		return nil, err
	}

	// 检查用户是否已经是群成员
	var existingMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("group_id = ? AND user_id = ?", req.GroupID, req.UserID).First(&existingMember).Error
	if err == nil {
		// 用户已经是群成员
		if existingMember.Status == 1 {
			l.Errorf("用户已经是群成员，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
			return nil, err
		} else {
			// 用户之前被踢出，现在重新加入
			err = l.svcCtx.DB.Model(&existingMember).Updates(map[string]interface{}{
				"status":    1,
				"join_time": time.Now(),
			}).Error
			if err != nil {
				l.Errorf("更新群成员状态失败: %v", err)
				return nil, err
			}

			// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel
		}
	} else {
		// 检查群组加入方式
		if group.JoinType == 1 {
			// 需要申请，创建申请记录
			// 获取该群入群申请的版本号（按群独立递增）
			requestVersion := l.svcCtx.VersionGen.GetNextVersion("group_join_requests", "group_id", req.GroupID)
			if requestVersion == -1 {
				l.Errorf("获取入群申请版本号失败")
				return nil, errors.New("获取版本号失败")
			}

			joinRequest := group_models.GroupJoinRequestModel{
				GroupID:         req.GroupID,
				ApplicantUserID: req.UserID,
				Message:         req.Message,
				Status:          0, // 待审核
				Version:         requestVersion,
			}
			err = l.svcCtx.DB.Create(&joinRequest).Error
			if err != nil {
				l.Errorf("创建入群申请失败: %v", err)
				return nil, err
			}

			resp = &types.GroupJoinRes{
				Version: requestVersion,
			}
			l.Infof("用户申请加入群组，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
			return resp, nil
		} else {
			// 获取该群成员的版本号（按群独立递增）
			memberVersion = l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
			if memberVersion == -1 {
				l.Errorf("获取群成员版本号失败")
				return nil, errors.New("获取版本号失败")
			}

			// 直接加入
			member := group_models.GroupMemberModel{
				GroupID:  req.GroupID,
				UserID:   req.UserID,
				Role:     3, // 普通成员
				Status:   1, // 正常状态
				JoinTime: time.Now(),
				Version:  memberVersion,
			}
			err = l.svcCtx.DB.Create(&member).Error
			if err != nil {
				l.Errorf("添加群成员失败: %v", err)
				return nil, err
			}

			// 注意：群成员的版本号通过 GroupMemberModel 的 Version 字段管理，不需要更新 GroupModel

			// 记录群成员变更日志
			changeLog := group_models.GroupMemberChangeLogModel{
				GroupID:    req.GroupID,
				UserID:     req.UserID,
				ChangeType: "join",
				OperatedBy: req.UserID,
				ChangeTime: time.Now(),
			}
			err = l.svcCtx.DB.Create(&changeLog).Error
			if err != nil {
				l.Errorf("记录群成员变更日志失败: %v", err)
				return nil, err
			}
		}
	}

	// 确保memberVersion有值（在直接加入的情况下）
	if memberVersion == 0 {
		memberVersion = l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID)
		if memberVersion == -1 {
			l.Errorf("获取群成员版本号失败")
			return nil, errors.New("获取版本号失败")
		}
	}

	// 异步通知群成员新成员加入
	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()

		// 获取群成员列表，用于推送通知
		response, err := l.svcCtx.GroupRpc.GetGroupMembers(ctx, &group_rpc.GetGroupMembersReq{
			GroupID: req.GroupID,
		})
		if err != nil {
			l.Errorf("获取群成员列表失败: %v", err)
			return
		}

		// 推送给所有群成员 - 群成员变动通知
		for _, member := range response.Members {
			if member.UserID != req.UserID { // 不通知操作者自己
				ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.GROUP_OPERATION, wsTypeConst.GroupMemberReceive, req.UserID, member.UserID, map[string]interface{}{
					"table": "group_members",
					"data": []map[string]interface{}{
						{
							"version": memberVersion,
							"groupId": req.GroupID,
						},
					},
				}, "")
			}
		}
	}()

	resp = &types.GroupJoinRes{
		Version: memberVersion,
	}

	l.Infof("用户加入群组成功，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
	return resp, nil
}
