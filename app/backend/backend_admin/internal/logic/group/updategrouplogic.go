package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新群组信息
func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateGroupLogic) UpdateGroup(req *types.UpdateGroupReq) (resp *types.UpdateGroupRes, err error) {
	// 检查群组是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("群组不存在: %d", req.Id)
			return nil, errors.New("群组不存在")
		}
		logx.Errorf("查询群组失败: %v", err)
		return nil, errors.New("查询群组失败")
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.FileName != "" {
		updateData["file_name"] = req.FileName
	}
	if req.Notice != "" {
		updateData["notice"] = req.Notice
	}
	if req.Status != 0 {
		updateData["status"] = req.Status
	}
	if req.Category != "" {
		updateData["category"] = req.Category
	}
	updateData["mute_all"] = req.MuteAll

	// 更新群组信息
	err = l.svcCtx.DB.Model(&group).Updates(updateData).Error
	if err != nil {
		logx.Errorf("更新群组信息失败: %v", err)
		return nil, errors.New("更新群组信息失败")
	}

	return &types.UpdateGroupRes{}, nil
}
