package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

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

			// 更新群组的成员版本号
			err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
				Where("group_id = ?", req.GroupID).
				Update("member_version", l.svcCtx.DB.Raw("member_version + 1")).Error
			if err != nil {
				l.Errorf("更新群组成员版本失败: %v", err)
			}
		}
	} else {
		// 检查群组加入方式
		if group.JoinType == 1 {
			// 需要申请，创建申请记录
			joinRequest := group_models.GroupJoinRequestModel{
				GroupID:         req.GroupID,
				ApplicantUserID: req.UserID,
				Message:         req.Message,
				Status:          0, // 待审核
			}
			err = l.svcCtx.DB.Create(&joinRequest).Error
			if err != nil {
				l.Errorf("创建入群申请失败: %v", err)
				return nil, err
			}

			resp = &types.GroupJoinRes{}
			l.Infof("用户申请加入群组，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
			return resp, nil
		} else {
			// 获取该群成员的版本号（按群独立递增）
			memberVersion := l.svcCtx.VersionGen.GetNextVersion("group_members", "group_id", req.GroupID, nil)
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

			// 更新群组的成员版本号
			err = l.svcCtx.DB.Model(&group_models.GroupModel{}).
				Where("group_id = ?", req.GroupID).
				Update("member_version", l.svcCtx.DB.Raw("member_version + 1")).Error
			if err != nil {
				l.Errorf("更新群组成员版本失败: %v", err)
			}

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

	resp = &types.GroupJoinRes{}

	l.Infof("用户加入群组成功，群组ID: %s, 用户ID: %s", req.GroupID, req.UserID)
	return resp, nil
}
