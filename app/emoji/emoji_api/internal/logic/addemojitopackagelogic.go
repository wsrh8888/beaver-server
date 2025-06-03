/*
 * @Author: renhao1 renhao1@100tal.com
 * @Date: 2025-03-27 21:05:02
 * @LastEditors: renhao1 renhao1@100tal.com
 * @LastEditTime: 2025-03-27 21:07:36
 * @FilePath: \beaver-server\app\emoji\emoji_api\internal\logic\addemojitopackagelogic.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AddEmojiToPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddEmojiToPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmojiToPackageLogic {
	return &AddEmojiToPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddEmojiToPackageLogic) AddEmojiToPackage(req *types.AddEmojiToPackageReq) (*types.AddEmojiToPackageRes, error) {
	// 1. 检查表情包是否存在
	var emojiPackage emoji_models.EmojiPackage
	err := l.svcCtx.DB.First(&emojiPackage, req.PackageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "表情包不存在")
		}
		return nil, status.Error(codes.Internal, "获取表情包失败")
	}

	// 2. 检查表情包状态
	if emojiPackage.Status != 1 {
		return nil, status.Error(codes.PermissionDenied, "表情包已禁用")
	}

	// 3. 检查是否是表情包作者
	if emojiPackage.UserID != req.UserID {
		return nil, status.Error(codes.PermissionDenied, "只有表情包作者可以添加表情")
	}

	// 4. 创建表情
	emoji := emoji_models.Emoji{
		FileId:   req.FileId,
		Title:    req.Title,
		AuthorID: req.UserID,
	}

	// 5. 保存表情到数据库
	err = l.svcCtx.DB.Create(&emoji).Error
	if err != nil {
		logx.Error("添加表情失败", err)
		return nil, status.Error(codes.Internal, "添加表情失败")
	}

	// 6. 创建表情包与表情的关联
	emojiPackageEmoji := emoji_models.EmojiPackageEmoji{
		PackageID: req.PackageID,
		EmojiID:   emoji.ID,
		SortOrder: 0, // 默认排序
	}

	err = l.svcCtx.DB.Create(&emojiPackageEmoji).Error
	if err != nil {
		logx.Error("添加表情到表情包关联失败", err)
		return nil, status.Error(codes.Internal, "添加表情到表情包关联失败")
	}

	return &types.AddEmojiToPackageRes{
		EmojiID: emoji.ID,
	}, nil
}
