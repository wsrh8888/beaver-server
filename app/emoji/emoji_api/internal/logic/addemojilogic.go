package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
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
	// 创建表情
	emoji := emoji_models.Emoji{
		FileName: req.FileName,
		Title:    req.Title,
		AuthorID: req.UserID,
	}

	// 保存到数据库
	err = l.svcCtx.DB.Create(&emoji).Error
	if err != nil {
		logx.Error("添加表情失败", err)
		return nil, err
	}

	// 如果指定了表情包ID，则添加到表情包
	if req.PackageID != 0 {
		// 创建表情包与表情的关联
		emojiPackageEmoji := emoji_models.EmojiPackageEmoji{
			PackageID: req.PackageID,
			EmojiID:   emoji.Id,
			SortOrder: 0, // 默认排序
		}

		err = l.svcCtx.DB.Create(&emojiPackageEmoji).Error
		if err != nil {
			logx.Error("添加表情到表情包失败", err)
			return nil, err
		}
	}

	// 添加表情并收藏
	favoriteEmoji := emoji_models.EmojiCollectEmoji{
		UserID:  req.UserID,
		EmojiID: emoji.Id, // 使用新增表情的ID
	}

	err = l.svcCtx.DB.Create(&favoriteEmoji).Error
	if err != nil {
		logx.Error("收藏表情失败", err)
		return nil, err
	}

	return &types.AddEmojiRes{}, nil
}
