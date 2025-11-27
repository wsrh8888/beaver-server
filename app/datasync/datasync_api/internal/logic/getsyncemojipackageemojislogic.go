package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncEmojiPackageEmojisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的表情包表情关联版本
func NewGetSyncEmojiPackageEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojiPackageEmojisLogic {
	return &GetSyncEmojiPackageEmojisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncEmojiPackageEmojisLogic) GetSyncEmojiPackageEmojis(req *types.GetSyncEmojiPackageEmojisReq) (resp *types.GetSyncEmojiPackageEmojisRes, err error) {
	// 调用Emoji RPC获取表情包表情关联版本信息
	emojiResp, err := l.svcCtx.EmojiRpc.GetEmojiPackageEmojis(l.ctx, &emoji_rpc.GetEmojiPackageEmojisReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取表情包表情关联版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情包表情关联版本信息", len(emojiResp.PackageEmojiVersions))

	// 转换为响应格式，确保返回空数组而不是null
	emojiPackageEmojiVersions := make([]types.EmojiPackageEmojiVersionItem, 0)
	if emojiResp.PackageEmojiVersions != nil {
		for _, pkgEmoji := range emojiResp.PackageEmojiVersions {
			emojiPackageEmojiVersions = append(emojiPackageEmojiVersions, types.EmojiPackageEmojiVersionItem{
				Id:      pkgEmoji.Id,
				Version: pkgEmoji.Version,
			})
		}
	}

	return &types.GetSyncEmojiPackageEmojisRes{
		EmojiPackageEmojiVersions: emojiPackageEmojiVersions,
		ServerTimestamp:           time.Now().UnixMilli(),
	}, nil
}
