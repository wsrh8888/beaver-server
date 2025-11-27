package logic

import (
	"context"

	"beaver/app/emoji/emoji_models"
	"beaver/app/emoji/emoji_rpc/internal/svc"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageEmojisListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmojiPackageEmojisListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageEmojisListByIdsLogic {
	return &GetEmojiPackageEmojisListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmojiPackageEmojisListByIdsLogic) GetEmojiPackageEmojisListByIds(in *emoji_rpc.GetEmojiPackageEmojisListByIdsReq) (*emoji_rpc.GetEmojiPackageEmojisListByIdsRes, error) {
	if len(in.Ids) == 0 {
		return &emoji_rpc.GetEmojiPackageEmojisListByIdsRes{PackageEmojis: []*emoji_rpc.EmojiPackageEmojiListById{}}, nil
	}

	var packageEmojis []emoji_models.EmojiPackageEmoji
	query := l.svcCtx.DB.Where("id IN (?)", in.Ids)

	// 时间戳过滤：只返回更新时间大于since的记录
	if in.Since > 0 {
		query = query.Where("updated_at > ?", in.Since)
	}

	err := query.Find(&packageEmojis).Error
	if err != nil {
		l.Errorf("查询表情包表情关联列表失败: ids=%v, since=%d, error=%v", in.Ids, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情包表情关联详情", len(packageEmojis))

	// 转换为响应格式
	var packageEmojiList []*emoji_rpc.EmojiPackageEmojiListById
	for _, pkgEmoji := range packageEmojis {
		packageEmojiList = append(packageEmojiList, &emoji_rpc.EmojiPackageEmojiListById{
			Id:        uint32(pkgEmoji.ID),
			PackageId: uint32(pkgEmoji.PackageID),
			EmojiId:   uint32(pkgEmoji.EmojiID),
			SortOrder: int32(pkgEmoji.SortOrder),
			Version:   pkgEmoji.Version,
			CreateAt:  pkgEmoji.CreatedAt.UnixMilli(),
			UpdateAt:  pkgEmoji.UpdatedAt.UnixMilli(),
		})
	}

	return &emoji_rpc.GetEmojiPackageEmojisListByIdsRes{PackageEmojis: packageEmojiList}, nil
}
