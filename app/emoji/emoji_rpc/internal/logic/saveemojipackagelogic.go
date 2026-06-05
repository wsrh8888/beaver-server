package logic

import (
	"context"
	"errors"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SaveEmojiPackageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveEmojiPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveEmojiPackageLogic {
	return &SaveEmojiPackageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SaveEmojiPackageLogic) SaveEmojiPackage(in *emoji_rpc.SaveEmojiPackageReq) (*emoji_rpc.SaveEmojiPackageRes, error) {
	if in.Delete != nil && *in.Delete {
		return l.deletePackage(in.PackageId)
	}
	if in.PackageId == "" {
		return l.createPackage(in)
	}
	return l.updatePackage(in)
}

func (l *SaveEmojiPackageLogic) createPackage(in *emoji_rpc.SaveEmojiPackageReq) (*emoji_rpc.SaveEmojiPackageRes, error) {
	var count int64
	if err := l.svcCtx.DB.Model(&emoji_models.EmojiPackage{}).
		Where("title = ? AND user_id = ?", in.Title, in.UserId).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("该用户已存在同名表情包")
	}

	version := l.svcCtx.VersionGen.GetNextVersion("emoji_package", "", "")
	if version == -1 {
		return nil, errors.New("获取版本号失败")
	}

	pkg := emoji_models.EmojiPackage{
		PackageID:   uuid.New().String(),
		Title:       in.Title,
		CoverFile:   in.CoverFile,
		UserID:      in.UserId,
		Description: in.Description,
		Type:        in.Type,
		Status:      1,
		Version:     version,
	}
	if err := l.svcCtx.DB.Create(&pkg).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.SaveEmojiPackageRes{PackageId: pkg.PackageID}, nil
}

func (l *SaveEmojiPackageLogic) updatePackage(in *emoji_rpc.SaveEmojiPackageReq) (*emoji_rpc.SaveEmojiPackageRes, error) {
	var pkg emoji_models.EmojiPackage
	if err := l.svcCtx.DB.Where("package_id = ?", in.PackageId).First(&pkg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表情包不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if in.PatchTitle != nil {
		if *in.PatchTitle != pkg.Title {
			var count int64
			if err := l.svcCtx.DB.Model(&emoji_models.EmojiPackage{}).
				Where("title = ? AND user_id = ? AND package_id != ?", *in.PatchTitle, pkg.UserID, in.PackageId).
				Count(&count).Error; err != nil {
				return nil, err
			}
			if count > 0 {
				return nil, errors.New("该用户已存在同名表情包")
			}
		}
		updates["title"] = *in.PatchTitle
	}
	if in.PatchCoverFile != nil {
		updates["cover_file"] = *in.PatchCoverFile
	}
	if in.PatchDescription != nil {
		updates["description"] = *in.PatchDescription
	}
	if in.PatchType != nil {
		updates["type"] = *in.PatchType
	}
	if in.PatchStatus != nil {
		updates["status"] = *in.PatchStatus
	}
	if len(updates) == 0 {
		return &emoji_rpc.SaveEmojiPackageRes{PackageId: in.PackageId}, nil
	}

	version := l.svcCtx.VersionGen.GetNextVersion("emoji_package", "", "")
	if version == -1 {
		return nil, errors.New("获取版本号失败")
	}
	updates["version"] = version

	if err := l.svcCtx.DB.Model(&pkg).Updates(updates).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.SaveEmojiPackageRes{PackageId: in.PackageId}, nil
}

func (l *SaveEmojiPackageLogic) deletePackage(packageID string) (*emoji_rpc.SaveEmojiPackageRes, error) {
	var pkg emoji_models.EmojiPackage
	if err := l.svcCtx.DB.Where("package_id = ?", packageID).First(&pkg).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("表情包不存在")
		}
		return nil, err
	}
	if err := l.svcCtx.DB.Delete(&pkg).Error; err != nil {
		return nil, err
	}
	return &emoji_rpc.SaveEmojiPackageRes{PackageId: packageID}, nil
}
