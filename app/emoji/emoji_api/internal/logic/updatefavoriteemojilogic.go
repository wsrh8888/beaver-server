package logic

import (
	"context"
	"errors"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFavoriteEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFavoriteEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteEmojiLogic {
	return &UpdateFavoriteEmojiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFavoriteEmojiLogic) UpdateFavoriteEmoji(req *types.UpdateFavoriteEmojiReq) (resp *types.UpdateFavoriteEmojiRes, err error) {
	// 查找表情
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.First(&emoji, req.EmojiID).Error
	if err != nil {
		logx.Error("表情不存在", err)
		return nil, errors.New("表情不存在")
	}

	// 检查是否已收藏
	var favoriteEmoji emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("user_id = ? AND emoji_id = ?", req.UserID, req.EmojiID).First(&favoriteEmoji).Error

	switch req.Type {
	case "favorite":
		if err == nil {
			// 已经收藏过了
			logx.Error("表情已收藏")
			return nil, errors.New("表情已收藏")
		}

		// 生成收藏版本号（按用户ID分区）
		collectVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_collect", "user_id", req.UserID)
		if collectVersion == -1 {
			logx.Error("生成收藏版本号失败")
			return nil, errors.New("生成版本号失败")
		}

		// 添加收藏
		newFavoriteEmoji := emoji_models.EmojiCollectEmoji{
			UserID:  req.UserID,
			EmojiID: req.EmojiID,
			Version: collectVersion,
		}
		err = l.svcCtx.DB.Create(&newFavoriteEmoji).Error
		if err != nil {
			logx.Error("收藏表情失败", err)
			return nil, errors.New("收藏表情失败")
		}
	case "unfavorite":
		if err != nil {
			// 没有收藏过
			logx.Error("表情未收藏")
			return nil, errors.New("表情未收藏")
		}

		// 软删除：设置IsDeleted为true并更新版本号（按用户ID分区）
		favoriteEmoji.IsDeleted = true
		favoriteEmoji.Version = l.svcCtx.VersionGen.GetNextVersion("emoji_collect", "user_id", req.UserID)
		if favoriteEmoji.Version == -1 {
			logx.Error("生成版本号失败")
			return nil, errors.New("生成版本号失败")
		}

		err = l.svcCtx.DB.Save(&favoriteEmoji).Error
		if err != nil {
			logx.Error("软删除收藏失败", err)
			return nil, errors.New("软删除收藏失败")
		}
	default:
		logx.Error("无效的操作类型")
		return nil, errors.New("无效的操作类型")
	}

	resp = &types.UpdateFavoriteEmojiRes{}
	return resp, nil
}
