package logic

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建表情图片
func NewCreateEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmojiLogic {
	return &CreateEmojiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEmojiLogic) CreateEmoji(req *types.CreateEmojiReq) (resp *types.CreateEmojiRes, err error) {
	// 检查表情名称是否已存在
	var count int64
	err = l.svcCtx.DB.Model(&emoji_models.Emoji{}).Where("title = ? AND author_id = ?", req.Title, req.AuthorID).Count(&count).Error
	if err != nil {
		logx.Errorf("检查表情名称失败: %v", err)
		return nil, errors.New("检查表情名称失败")
	}
	if count > 0 {
		return nil, errors.New("该创建者已存在同名表情")
	}

	// 创建表情
	emoji := emoji_models.Emoji{
		FileName: req.FileName,
		Title:    req.Title,
		AuthorID: req.AuthorID,
	}

	err = l.svcCtx.DB.Create(&emoji).Error
	if err != nil {
		logx.Errorf("创建表情失败: %v", err)
		return nil, errors.New("创建表情失败")
	}

	return &types.CreateEmojiRes{
		Id: strconv.Itoa(int(emoji.Id)),
	}, nil
}
