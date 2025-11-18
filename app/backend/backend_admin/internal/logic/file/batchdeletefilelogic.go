package logic

import (
	"context"
	"fmt"
	"os"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除文件
func NewBatchDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFileLogic {
	return &BatchDeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteFileLogic) BatchDeleteFile(req *types.BatchDeleteFileReq) (resp *types.BatchDeleteFileRes, err error) {
	// 先查询要删除的文件
	var files []file_models.FileModel
	err = l.svcCtx.DB.Where("id IN ?", req.Ids).Find(&files).Error
	if err != nil {
		logx.Errorf("查询要删除的文件失败: %v", err)
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("没有找到要删除的文件")
	}

	// 批量删除数据库记录
	err = l.svcCtx.DB.Where("id IN ?", req.Ids).Delete(&file_models.FileModel{}).Error
	if err != nil {
		logx.Errorf("批量删除文件记录失败: %v", err)
		return nil, err
	}

	// 删除物理文件
	var deletedCount int
	var failedPaths []string

	for _, file := range files {
		if file.Path != "" {
			if err := os.Remove(file.Path); err != nil {
				logx.Errorf("删除物理文件失败: %v, 路径: %s", err, file.Path)
				failedPaths = append(failedPaths, file.Path)
			} else {
				deletedCount++
			}
		}
	}

	logx.Infof("批量删除完成, 数据库删除: %d 条记录, 物理文件删除成功: %d 个, 失败: %d 个",
		len(files), deletedCount, len(failedPaths))

	if len(failedPaths) > 0 {
		logx.Errorf("以下物理文件删除失败: %v", failedPaths)
	}

	return &types.BatchDeleteFileRes{}, nil
}
