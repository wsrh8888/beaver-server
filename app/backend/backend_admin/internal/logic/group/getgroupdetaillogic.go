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

type GetGroupDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群组详情
func NewGetGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupDetailLogic {
	return &GetGroupDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupDetailLogic) GetGroupDetail(req *types.GetGroupDetailReq) (resp *types.GetGroupDetailRes, err error) {
	var group group_models.GroupModel

	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("群组不存在: %d", req.Id)
			return nil, errors.New("群组不存在")
		}
		logx.Errorf("查询群组详情失败: %v", err)
		return nil, errors.New("查询群组详情失败")
	}

	return &types.GetGroupDetailRes{
		Id:        group.Id,
		GroupId:   group.GroupID,
		Type:      int(group.Type),
		Title:     group.Title,
		FileName:  group.Avatar,
		CreatorId: group.CreatorID,
		Notice:    group.Notice,
		Status:    int(group.Status),
		CreatedAt: group.CreatedAt.String(),
		UpdatedAt: group.UpdatedAt.String(),
	}, nil
}
