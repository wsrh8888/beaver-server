package logic

import (
	"context"
	"errors"
	"os"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除文件
func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileLogic) DeleteFile(req *types.DeleteFileReq) (resp *types.DeleteFileRes, err error) {
	// 先查询文件是否存在
	var file file_models.FileModel
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&file).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("文件不存在, Id: %d", req.Id)
			return nil, errors.New("文件不存在")
		}
		logx.Errorf("查询文件失败: %v", err)
		return nil, err
	}

	// 删除数据库记录
	err = l.svcCtx.DB.Delete(&file).Error
	if err != nil {
		logx.Errorf("删除文件记录失败: %v", err)
		return nil, err
	}

	// 删除物理文件（如果存在）
	if file.Path != "" {
		if err := os.Remove(file.Path); err != nil {
			logx.Errorf("删除物理文件失败: %v, 路径: %s", err, file.Path)
			// 物理文件删除失败不影响数据库删除的结果，仅记录日志
		}
	}

	logx.Infof("文件删除成功, Id: %d, 文件名: %s", file.Id, file.OriginalName)
	return &types.DeleteFileRes{}, nil
}
