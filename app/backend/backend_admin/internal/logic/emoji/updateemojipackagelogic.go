package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新表情包集合
func NewUpdateEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmojiPackageLogic {
	return &UpdateEmojiPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmojiPackageLogic) UpdateEmojiPackage(req *types.UpdateEmojiPackageReq) (resp *types.UpdateEmojiPackageRes, err error) {
	// 检查表情包是否存在
	var pkg emoji_models.EmojiPackage
	err = l.svcCtx.DB.Where("package_id = ?", req.PackageId).First(&pkg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情包不存在: %s", req.PackageId)
			return nil, errors.New("表情包不存在")
		}
		logx.Errorf("查询表情包失败: %v", err)
		return nil, errors.New("查询表情包失败")
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.Title != nil {
		// 检查同用户下的表情包名称是否重复
		if *req.Title != pkg.Title {
			var count int64
			err = l.svcCtx.DB.Model(&emoji_models.EmojiPackage{}).
				Where("title = ? AND user_id = ? AND package_id != ?", *req.Title, pkg.UserID, req.PackageId).
				Count(&count).Error
			if err != nil {
				logx.Errorf("检查表情包名称失败: %v", err)
				return nil, errors.New("检查表情包名称失败")
			}
			if count > 0 {
				return nil, errors.New("该用户已存在同名表情包")
			}
		}
		updateData["title"] = *req.Title
	}
	if req.CoverFile != nil {
		updateData["cover_file"] = *req.CoverFile
	}

	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.Type != nil {
		updateData["type"] = *req.Type
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	// 更新表情包信息
	if len(updateData) > 0 {
		err = l.svcCtx.DB.Model(&pkg).Updates(updateData).Error
		if err != nil {
			logx.Errorf("更新表情包信息失败: %v", err)
			return nil, errors.New("更新表情包信息失败")
		}
	}

	return &types.UpdateEmojiPackageRes{}, nil
}
