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

type DeleteEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除表情图片
func NewDeleteEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmojiLogic {
	return &DeleteEmojiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteEmojiLogic) DeleteEmoji(req *types.DeleteEmojiReq) (resp *types.DeleteEmojiRes, err error) {
	// 检查表情是否存在
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.Where("emoji_id = ?", req.EmojiId).First(&emoji).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("表情不存在: %s", req.EmojiId)
			return nil, errors.New("表情不存在")
		}
		logx.Errorf("查询表情失败: %v", err)
		return nil, errors.New("查询表情失败")
	}

	// 使用逻辑删除
	err = l.svcCtx.DB.Delete(&emoji).Error
	if err != nil {
		logx.Errorf("删除表情失败: %v", err)
		return nil, errors.New("删除表情失败")
	}

	return &types.DeleteEmojiRes{}, nil
}
