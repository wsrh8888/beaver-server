package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建表情包集合
func NewCreateEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmojiPackageLogic {
	return &CreateEmojiPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEmojiPackageLogic) CreateEmojiPackage(req *types.CreateEmojiPackageReq) (resp *types.CreateEmojiPackageRes, err error) {
	// 检查表情包名称是否已存在
	var count int64
	err = l.svcCtx.DB.Model(&emoji_models.EmojiPackage{}).Where("title = ? AND user_id = ?", req.Title, req.UserID).Count(&count).Error
	if err != nil {
		logx.Errorf("检查表情包名称失败: %v", err)
		return nil, errors.New("检查表情包名称失败")
	}
	if count > 0 {
		return nil, errors.New("该用户已存在同名表情包")
	}

	// 创建表情包
	pkg := emoji_models.EmojiPackage{
		PackageID:   uuid.New().String(),
		Title:       req.Title,
		UserID:      req.UserID,
		Description: req.Description,
		Type:        req.Type,
		Status:      1, // 默认状态为1
	}

	// 如果提供了封面文件，则设置
	if req.CoverFile != nil && *req.CoverFile != "" {
		pkg.CoverFile = *req.CoverFile
	}

	err = l.svcCtx.DB.Create(&pkg).Error
	if err != nil {
		logx.Errorf("创建表情包失败: %v", err)
		return nil, errors.New("创建表情包失败")
	}

	return &types.CreateEmojiPackageRes{
		PackageId: pkg.PackageID,
	}, nil
}
