package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 保存文件信息到数据库
func NewSaveFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFileLogic {
	return &SaveFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveFileLogic) SaveFile(req *types.SaveFileReq) (resp *types.SaveFileRes, err error) {
	l.Logger.Infof("开始保存文件信息: %s, 大小: %d, 类型: %s", req.OriginalName, req.Size, req.Type)

	// 检查文件是否已经存在于数据库中（通过MD5）
	var existingFile file_models.FileModel
	err = l.svcCtx.DB.Take(&existingFile, "md5 = ?", req.Md5).Error
	if err == nil {
		l.Logger.Infof("文件已存在，返回现有文件ID: %s", existingFile.FileKey)
		return &types.SaveFileRes{
			FileKey: existingFile.FileKey,
		}, nil
	}

	// 从文件名中提取后缀
	suffix := "jpg" // 默认后缀
	if strings.Contains(req.OriginalName, ".") {
		parts := strings.Split(req.OriginalName, ".")
		if len(parts) > 1 {
			suffix = strings.ToLower(parts[len(parts)-1])
		}
	}

	// 使用MD5作为文件名，这样相同内容的文件会有相同的FileKey，实现缓存复用
	fileKey := req.Md5 + "." + suffix
	l.Logger.Infof("使用MD5生成文件ID: %s", fileKey)

	// 创建新的文件记录
	newFileModel := &file_models.FileModel{
		FileKey:      fileKey,
		OriginalName: req.OriginalName,
		Size:         req.Size,
		Path:         req.Path,
		Md5:          req.Md5,
		Type:         req.Type,
	}

	// 保存到数据库
	err = l.svcCtx.DB.Create(newFileModel).Error
	if err != nil {
		l.Logger.Errorf("保存文件信息到数据库失败: %v", err)
		return nil, errors.New("保存文件信息失败")
	}

	l.Logger.Infof("文件信息保存成功: %s", fileKey)

	return &types.SaveFileRes{
		FileKey: fileKey,
	}, nil
}
