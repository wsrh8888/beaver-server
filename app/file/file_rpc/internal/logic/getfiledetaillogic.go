package logic

import (
	"context"
	"errors"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过fileName查询文件详情
func (l *GetFileDetailLogic) GetFileDetail(in *file_rpc.GetFileDetailReq) (*file_rpc.GetFileDetailRes, error) {
	var file file_models.FileModel

	// 通过fileName查询文件信息
	err := l.svcCtx.DB.Take(&file, "file_key = ?", in.FileKey).Error
	if err != nil {
		logx.Errorf("查询文件失败: %s", err.Error())
		return nil, errors.New("文件不存在")
	}

	// 返回文件详情
	return &file_rpc.GetFileDetailRes{
		FileKey:      file.FileKey,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		Path:         file.Path,
		Md5:          file.Md5,
		Type:         file.Type,
		CreatedAt:    file.CreatedAt.String(),
		UpdatedAt:    file.UpdatedAt.String(),
	}, nil
}
