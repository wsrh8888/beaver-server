package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserEmojiPackageCollectsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserEmojiPackageCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserEmojiPackageCollectsLogic {
	return &GetUserEmojiPackageCollectsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserEmojiPackageCollectsLogic) GetUserEmojiPackageCollects(in *emoji_rpc.GetUserEmojiPackageCollectsReq) (*emoji_rpc.GetUserEmojiPackageCollectsRes, error) {
	var packageCollects []emoji_models.EmojiPackageCollect
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		query = query.Where("updated_at > ?", sinceTime)
	}

	err := query.Find(&packageCollects).Error
	if err != nil {
		l.Errorf("查询用户收藏表情包版本失败: userId=%s, since=%d, error=%v", in.UserId, in.Since, err)
		return nil, err
	}

	l.Infof("查询到用户 %s 的 %d 个收藏表情包版本变更", in.UserId, len(packageCollects))

	// 转换为版本摘要格式
	var packageCollectVersions []*emoji_rpc.EmojiPackageCollectVersionItem
	for _, pkgCollect := range packageCollects {
		packageCollectVersions = append(packageCollectVersions, &emoji_rpc.EmojiPackageCollectVersionItem{
			PackageCollectId: pkgCollect.PackageCollectID,
			Version:          pkgCollect.Version,
		})
	}

	return &emoji_rpc.GetUserEmojiPackageCollectsRes{
		EmojiPackageCollectVersions: packageCollectVersions,
		ServerTimestamp:             time.Now().UnixMilli(),
	}, nil
}
