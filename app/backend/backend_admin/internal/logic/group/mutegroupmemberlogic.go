package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type MuteGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 禁言群成员
func NewMuteGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MuteGroupMemberLogic {
	return &MuteGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MuteGroupMemberLogic) MuteGroupMember(req *types.MuteGroupMemberReq) (resp *types.MuteGroupMemberRes, err error) {
	// 检查群组成员是否存在
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("群组成员不存在: %d", req.Id)
			return nil, errors.New("群组成员不存在")
		}
		logx.Errorf("查询群组成员失败: %v", err)
		return nil, errors.New("查询群组成员失败")
	}

	// 计算禁言到期时间
	var endTime *time.Time
	if req.ProhibitionTime > 0 {
		expireTime := time.Now().Add(time.Duration(req.ProhibitionTime) * time.Minute)
		endTime = &expireTime
	}

	// 更新禁言状态
	err = l.svcCtx.DB.Model(&member).Updates(map[string]interface{}{
		"prohibition_time": &req.ProhibitionTime,
		"prohibition_end":  endTime,
	}).Error
	if err != nil {
		logx.Errorf("更新成员禁言状态失败: %v", err)
		return nil, errors.New("更新成员禁言状态失败")
	}

	return &types.MuteGroupMemberRes{}, nil
}
