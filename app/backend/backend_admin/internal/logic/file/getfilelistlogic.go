package logic

import (
	"context"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文件列表
func NewGetFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileListLogic) GetFileList(req *types.GetFileListReq) (resp *types.GetFileListRes, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&file_models.FileModel{})

	// 文件类型筛选
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	// 文件名关键词搜索
	if req.Keywords != "" {
		query = query.Where("file_key LIKE ? OR original_name LIKE ?", "%"+req.Keywords+"%", "%"+req.Keywords+"%")
	}

	// 查询总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		logx.Errorf("查询文件总数失败: %v", err)
		return nil, err
	}

	// 查询列表
	var files []file_models.FileModel
	err = query.Order("created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&files).Error
	if err != nil {
		logx.Errorf("查询文件列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.GetFileListItem, len(files))
	for i, file := range files {
		list[i] = types.GetFileListItem{
			Id:           file.Id,
			FileName:     file.FileKey,
			OriginalName: file.OriginalName,
			Size:         file.Size,
			Path:         file.Path,
			Md5:          file.Md5,
			Type:         file.Type,
			CreatedAt:    time.Time(file.CreatedAt).Format(time.RFC3339),
			UpdatedAt:    time.Time(file.UpdatedAt).Format(time.RFC3339),
		}
	}

	return &types.GetFileListRes{
		List:  list,
		Total: total,
	}, nil
}
