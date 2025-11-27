package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageEmojisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackageEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageEmojisLogic {
	return &GetEmojiPackageEmojisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackageEmojisLogic) GetEmojiPackageEmojis(in *emoji_rpc.GetEmojiPackageEmojisReq) (*emoji_rpc.GetEmojiPackageEmojisRes, error) {
	// 查询用户相关的表情包表情关联（通过表情包权限控制）
	// 先获取用户有权限的表情包ID列表
	var packageIds []uint
	subQuery := l.svcCtx.DB.Model(&emoji_models.EmojiPackage{}).
		Where("(type = ? OR user_id = ?) AND status = ?", "official", in.UserId, 1).
		Select("id")

	var packageEmojis []emoji_models.EmojiPackageEmoji
	query := l.svcCtx.DB.Where("package_id IN (?)", subQuery)

	// 增量同步：只返回版本号大于since的记录
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&packageEmojis).Error
	if err != nil {
		l.Errorf("查询表情包表情关联版本信息失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个表情包表情关联版本信息", in.UserId, len(packageEmojis))

	// 转换为响应格式
	var packageEmojiVersions []*emoji_rpc.GetEmojiPackageEmojisRes_EmojiPackageEmojiVersion
	for _, pkgEmoji := range packageEmojis {
		packageEmojiVersions = append(packageEmojiVersions, &emoji_rpc.GetEmojiPackageEmojisRes_EmojiPackageEmojiVersion{
			Id:      uint32(pkgEmoji.ID),
			Version: pkgEmoji.Version,
		})
	}

	return &emoji_rpc.GetEmojiPackageEmojisRes{PackageEmojiVersions: packageEmojiVersions}, nil
}
