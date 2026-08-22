package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveCircleMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveCircleMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCircleMembersLogic {
	return &RemoveCircleMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveCircleMembersLogic) RemoveCircleMembers(req *types.RemoveCircleMembersReq) (resp *types.RemoveCircleMembersRes, err error) {
	if len(req.UserIds) == 0 {
		return nil, fmt.Errorf("请选择要移除的成员")
	}

	var operator circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&operator).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if operator.Role > 2 {
		return nil, fmt.Errorf("仅圈主和管理员可移除成员")
	}

	var members []circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id IN ?", req.CircleID, req.UserIds).Find(&members).Error; err != nil {
		return nil, fmt.Errorf("查询成员失败")
	}

	for _, member := range members {
		if member.UserID == req.UserID {
			return nil, fmt.Errorf("不能移除自己")
		}
		if operator.Role == 1 {
			if member.Role == 1 {
				return nil, fmt.Errorf("不能移除圈主")
			}
		} else if operator.Role == 2 {
			if member.Role != 3 {
				return nil, fmt.Errorf("管理员只能移除普通成员")
			}
		}
	}

	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id IN ?", req.CircleID, req.UserIds).
		Delete(&circle_models.CircleMemberModel{}).Error; err != nil {
		return nil, fmt.Errorf("移除成员失败: %v", err)
	}

	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Update("version", circleVersion)

	return &types.RemoveCircleMembersRes{}, nil
}
