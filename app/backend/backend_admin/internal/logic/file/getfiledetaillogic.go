package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetFileDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文件详情
func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileDetailLogic) GetFileDetail(req *types.GetFileDetailReq) (resp *types.GetFileDetailRes, err error) {
	var file file_models.FileModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("文件不存在, Id: %d", req.Id)
			return nil, errors.New("文件不存在")
		}
		logx.Errorf("查询文件详情失败: %v", err)
		return nil, err
	}

	return &types.GetFileDetailRes{
		FileInfo: types.FileInfo{
			Id:           file.Id,
			FileName:     file.FileName,
			OriginalName: file.OriginalName,
			Size:         file.Size,
			Path:         file.Path,
			Md5:          file.Md5,
			Type:         file.Type,
			CreatedAt:    time.Time(file.CreatedAt).Format(time.RFC3339),
			UpdatedAt:    time.Time(file.UpdatedAt).Format(time.RFC3339),
		},
	}, nil
}
