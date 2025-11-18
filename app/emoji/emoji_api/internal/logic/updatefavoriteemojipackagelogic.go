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

type UpdateFavoriteEmojiPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFavoriteEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFavoriteEmojiPackageLogic {
	return &UpdateFavoriteEmojiPackageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFavoriteEmojiPackageLogic) UpdateFavoriteEmojiPackage(req *types.UpdateFavoriteEmojiPackageReq) (*types.UpdateFavoriteEmojiPackageRes, error) {
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

	// 3. 检查是否已收藏
	var collectRecord emoji_models.EmojiPackageCollect
	err = l.svcCtx.DB.Where("user_id = ? AND package_id = ?", req.UserID, req.PackageID).
		First(&collectRecord).Error

	// 4. 根据操作类型处理
	if req.Type == "favorite" {
		// 收藏
		if err == nil {
			return nil, status.Error(codes.AlreadyExists, "已经收藏过了")
		}
		collectRecord = emoji_models.EmojiPackageCollect{
			UserID:    req.UserID,
			PackageID: req.PackageID,
		}
		err = l.svcCtx.DB.Create(&collectRecord).Error
		if err != nil {
			return nil, status.Error(codes.Internal, "收藏失败")
		}
	} else if req.Type == "unfavorite" {
		// 取消收藏
		if err != nil {
			return nil, status.Error(codes.NotFound, "未收藏过")
		}
		err = l.svcCtx.DB.Delete(&collectRecord).Error
		if err != nil {
			return nil, status.Error(codes.Internal, "取消收藏失败")
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}

	return &types.UpdateFavoriteEmojiPackageRes{}, nil
}
