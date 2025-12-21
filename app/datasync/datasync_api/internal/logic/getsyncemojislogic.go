package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncEmojisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取表情基础数据版本信息
func NewGetSyncEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojisLogic {
	return &GetSyncEmojisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncEmojisLogic) GetSyncEmojis(req *types.GetSyncEmojisReq) (resp *types.GetSyncEmojisRes, err error) {
	// 调用Emoji RPC获取表情版本信息
	emojiResp, err := l.svcCtx.EmojiRpc.GetEmojis(l.ctx, &emoji_rpc.GetEmojisReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取表情版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情版本信息", len(emojiResp.EmojiVersions))

	// 转换为响应格式，确保返回空数组而不是null
	emojiVersions := make([]types.EmojiVersionItem, 0)
	if emojiResp.EmojiVersions != nil {
		for _, emoji := range emojiResp.EmojiVersions {
			emojiVersions = append(emojiVersions, types.EmojiVersionItem{
				EmojiId: emoji.EmojiId,
				Version: emoji.Version,
			})
		}
	}

	return &types.GetSyncEmojisRes{
		EmojiVersions:   emojiVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
