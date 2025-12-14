package logic

import (
	"context"
	"errors"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmojiLogic {
	return &AddEmojiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddEmojiLogic) AddEmoji(req *types.AddEmojiReq) (resp *types.AddEmojiRes, err error) {
	// 先按 FileKey 查重，已有则复用，不重复落库
	var emoji emoji_models.Emoji
	err = l.svcCtx.DB.Where("file_key = ?", req.FileKey).First(&emoji).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 仅当不存在时创建新的 emoji
	if errors.Is(err, gorm.ErrRecordNotFound) {
		emojiVersion := l.svcCtx.VersionGen.GetNextVersion("emoji", "", "")
		if emojiVersion == -1 {
			logx.Error("生成表情版本号失败")
			return nil, errors.New("生成版本号失败")
		}

		emoji = emoji_models.Emoji{
			EmojiID: uuid.New().String(),
			FileKey: req.FileKey,
			Title:   req.Title,
			Version: emojiVersion,
			EmojiInfo: &emoji_models.EmojiInfo{
				Width:  req.EmojiInfo.Width,
				Height: req.EmojiInfo.Height,
			},
		}

		if err := l.svcCtx.DB.Create(&emoji).Error; err != nil {
			logx.Error("添加表情失败", err)
			return nil, err
		}
	}

	// 生成收藏版本号（按用户ID分区）
	collectVersion := l.svcCtx.VersionGen.GetNextVersion("emoji_collect", "user_id", req.UserID)
	if collectVersion == -1 {
		logx.Error("生成收藏版本号失败")
		return nil, errors.New("生成版本号失败")
	}

	// 添加表情并收藏
	favoriteEmoji := emoji_models.EmojiCollectEmoji{
		EmojiCollectID: uuid.New().String(),
		UserID:         req.UserID,
		EmojiID:        emoji.EmojiID,
		Version:        collectVersion,
		PackageID:      req.PackageID,
	}

	// 去重：同一用户对同一 emoji 已收藏则跳过创建
	var existFavorite emoji_models.EmojiCollectEmoji
	err = l.svcCtx.DB.Where("user_id = ? AND emoji_id = ?", req.UserID, emoji.EmojiID).First(&existFavorite).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := l.svcCtx.DB.Create(&favoriteEmoji).Error; err != nil {
			logx.Error("收藏表情失败", err)
			return nil, err
		}

		// 异步通过 WS 通知用户其他客户端
		go func(etcdAddr string, userId string, emojiCollectId string, collectVersion int64, emojiVersion int64) {
			// 构建表更新数据
			var tableUpdates []map[string]interface{}

			// 通知表情表更新（如果创建了新表情）
			if emojiVersion > 0 {
				emojiUpdates := map[string]interface{}{
					"table": "emoji",
					"data": []map[string]interface{}{
						{
							"version": emojiVersion,
							"emojiId": emoji.EmojiID,
						},
					},
				}
				tableUpdates = append(tableUpdates, emojiUpdates)
			}

			// 通知表情收藏表更新
			collectUpdates := map[string]interface{}{
				"table":  "emoji_collect",
				"userId": userId,
				"data": []map[string]interface{}{
					{
						"version":        collectVersion,
						"emojiCollectId": emojiCollectId,
					},
				},
			}
			tableUpdates = append(tableUpdates, collectUpdates)

			// 通知给自己（用户ID作为接收者，空字符串作为发送者表示系统操作）
			ajax.SendMessageToWs(etcdAddr, wsCommandConst.EMOJI, wsTypeConst.EmojiReceive, "", userId, map[string]interface{}{
				"tableUpdates": tableUpdates,
			}, "")
		}(l.svcCtx.Config.Etcd, req.UserID, favoriteEmoji.EmojiCollectID, collectVersion, emoji.Version)
	}

	return &types.AddEmojiRes{}, nil
}
