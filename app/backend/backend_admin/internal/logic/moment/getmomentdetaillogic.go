package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetMomentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态详情
func NewGetMomentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentDetailLogic {
	return &GetMomentDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentDetailLogic) GetMomentDetail(req *types.GetMomentDetailReq) (resp *types.GetMomentDetailRes, err error) {
	var moment moment_models.MomentModel

	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&moment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("动态不存在: %d", req.Id)
			return nil, errors.New("动态不存在")
		}
		logx.Errorf("查询动态详情失败: %v", err)
		return nil, errors.New("查询动态详情失败")
	}

	// 转换文件信息
	var files []types.FileInfo
	if moment.Files != nil {
		for _, file := range *moment.Files {
			files = append(files, types.FileInfo{
				FileName: file.FileName,
			})
		}
	}

	return &types.GetMomentDetailRes{
		MomentInfo: types.MomentInfo{
			Id:        moment.Id,
			UserId:    moment.UserID,
			Content:   moment.Content,
			Files:     files,
			IsDeleted: moment.IsDeleted,
			CreatedAt: moment.CreatedAt.String(),
			UpdatedAt: moment.UpdatedAt.String(),
		},
	}, nil
}
