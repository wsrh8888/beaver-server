package logic

import (
	"context"
	"os"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BatchDeleteFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFilesLogic {
	return &BatchDeleteFilesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BatchDeleteFilesLogic) BatchDeleteFiles(in *file_rpc.BatchDeleteFilesReq) (*file_rpc.BatchDeleteFilesRes, error) {
	if len(in.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "ids不能为空")
	}

	var files []file_models.FileModel
	if err := l.svcCtx.DB.Where("id IN ?", in.Ids).Find(&files).Error; err != nil {
		l.Errorf("查询文件失败: %v", err)
		return nil, err
	}
	if len(files) == 0 {
		return nil, status.Error(codes.NotFound, "没有找到要删除的文件")
	}

	if err := l.svcCtx.DB.Where("id IN ?", in.Ids).Delete(&file_models.FileModel{}).Error; err != nil {
		l.Errorf("批量删除文件记录失败: %v", err)
		return nil, status.Error(codes.Internal, "批量删除失败")
	}

	paths := make([]string, 0, len(files))
	for _, file := range files {
		if file.Path == "" {
			continue
		}
		paths = append(paths, file.Path)
		if err := os.Remove(file.Path); err != nil {
			l.Errorf("删除物理文件失败 path=%s: %v", file.Path, err)
		}
	}

	return &file_rpc.BatchDeleteFilesRes{Paths: paths}, nil
}
