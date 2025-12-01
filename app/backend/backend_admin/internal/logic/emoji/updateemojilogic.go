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

type UpdateEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新表情图片
func NewUpdateEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmojiLogic {
	return &UpdateEmojiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmojiLogic) UpdateEmoji(req *types.UpdateEmojiReq) (resp *types.UpdateEmojiRes, err error) {
	// 检查表情是否存在
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.Where("uuid = ?", req.UUID).First(&emoji).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情不存在: %s", req.UUID)
			return nil, errors.New("表情不存在")
		}
		logx.Errorf("查询表情失败: %v", err)
		return nil, errors.New("查询表情失败")
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.FileKey != nil {
		updateData["file_name"] = *req.FileKey
	}
	if req.Title != nil {
		// 检查同创建者下的表情名称是否重复
		if *req.Title != emoji.Title {
			var count int64
			err = l.svcCtx.DB.Model(&emoji_models.Emoji{}).
				Where("title = ? AND uuid != ?", *req.Title, req.UUID).
				Count(&count).Error
			if err != nil {
				logx.Errorf("检查表情名称失败: %v", err)
				return nil, errors.New("检查表情名称失败")
			}
			if count > 0 {
				return nil, errors.New("该创建者已存在同名表情")
			}
		}
		updateData["title"] = *req.Title
	}

	// 更新表情信息
	if len(updateData) > 0 {
		err = l.svcCtx.DB.Model(&emoji).Updates(updateData).Error
		if err != nil {
			logx.Errorf("更新表情信息失败: %v", err)
			return nil, errors.New("更新表情信息失败")
		}
	}

	return &types.UpdateEmojiRes{}, nil
}
