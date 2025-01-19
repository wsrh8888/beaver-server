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
		FileUrl:  req.FileUrl,
		Title:    req.Title,
		AuthorID: req.UserID,
	}
	if req.PackageID == 0 {
		emoji.PackageID = nil
	}

	// 保存到数据库
	err = l.svcCtx.DB.Create(&emoji).Error
	if err != nil {
		logx.Error("添加表情失败", err)
		return nil, err
	}

	// 添加表情并收藏
	favoriteEmoji := emoji_models.EmojiCollectEmoji{
		UserID:  req.UserID,
		EmojiID: emoji.ID, // 使用新增表情的ID
	}

	err = l.svcCtx.DB.Create(&favoriteEmoji).Error
	if err != nil {
		logx.Error("收藏表情失败", err)
		return nil, err
	}

	return &types.AddEmojiRes{}, nil
}
