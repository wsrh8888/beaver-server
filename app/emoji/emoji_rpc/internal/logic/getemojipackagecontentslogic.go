package logic

import (
	"context"
	"time"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageContentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackageContentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageContentsLogic {
	return &GetEmojiPackageContentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackageContentsLogic) GetEmojiPackageContents(in *emoji_rpc.GetEmojiPackageContentsReq) (*emoji_rpc.GetEmojiPackageContentsRes, error) {
	var packageContents []emoji_models.EmojiPackageEmoji

	// 时间戳过滤：只返回更新时间大于since的记录
	query := l.svcCtx.DB
	if in.Since > 0 {
		sinceTime := time.UnixMilli(in.Since)
		query = query.Where("updated_at > ?", sinceTime)
	}

	err := query.Find(&packageContents).Error
	if err != nil {
		l.Errorf("查询表情包内容版本失败: since=%d, error=%v", in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情包内容版本变更", len(packageContents))

	// 转换为版本摘要格式
	var contentVersions []*emoji_rpc.EmojiPackageContentVersionItem
	for _, content := range packageContents {
		contentVersions = append(contentVersions, &emoji_rpc.EmojiPackageContentVersionItem{
			PackageId: content.PackageID,
			Version:   content.Version,
		})
	}

	return &emoji_rpc.GetEmojiPackageContentsRes{
		EmojiPackageContentVersions: contentVersions,
		ServerTimestamp:             time.Now().UnixMilli(),
	}, nil
}
