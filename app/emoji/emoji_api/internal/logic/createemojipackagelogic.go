package logic

import (
	"context"

	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"beaver/app/emoji/emoji_models"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmojiPackageLogic {
	return &CreateEmojiPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEmojiPackageLogic) CreateEmojiPackage(req *types.CreateEmojiPackageReq) (*types.CreateEmojiPackageRes, error) {
	// 创建表情包
	emojiPackage := emoji_models.EmojiPackage{
		Title:       req.Title,
		CoverFile:   req.CoverFile,
		Description: req.Description,
		UserID:      req.UserID,
		Type:        "user", // 默认为用户自定义类型
		Status:      1,      // 默认为正常状态
	}

	// 保存到数据库
	err := l.svcCtx.DB.Create(&emojiPackage).Error
	if err != nil {
		logx.Error("创建表情包失败", err)
		return nil, status.Error(codes.Internal, "创建表情包失败")
	}

	return &types.CreateEmojiPackageRes{
		PackageID: emojiPackage.Id,
	}, nil
}
