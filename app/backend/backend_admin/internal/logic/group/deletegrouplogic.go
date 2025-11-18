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

type DeleteGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除群组
func NewDeleteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGroupLogic {
	return &DeleteGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteGroupLogic) DeleteGroup(req *types.DeleteGroupReq) (resp *types.DeleteGroupRes, err error) {
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

	// 执行逻辑删除
	now := time.Now()
	err = l.svcCtx.DB.Model(&group).Updates(map[string]interface{}{
		"deleted_at":    &now,
		"dissolve_time": &now,
		"status":        2, // 已解散状态
	}).Error
	if err != nil {
		logx.Errorf("删除群组失败: %v", err)
		return nil, errors.New("删除群组失败")
	}

	return &types.DeleteGroupRes{}, nil
}
