package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncEmojiPackagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的表情包版本
func NewGetSyncEmojiPackagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojiPackagesLogic {
	return &GetSyncEmojiPackagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncEmojiPackagesLogic) GetSyncEmojiPackages(req *types.GetSyncEmojiPackagesReq) (resp *types.GetSyncEmojiPackagesRes, err error) {
	// 调用Emoji RPC获取表情包版本信息
	emojiResp, err := l.svcCtx.EmojiRpc.GetEmojiPackages(l.ctx, &emoji_rpc.GetEmojiPackagesReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取表情包版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个表情包版本信息", len(emojiResp.PackageVersions))

	// 转换为响应格式，确保返回空数组而不是null
	emojiPackageVersions := make([]types.EmojiPackageVersionItem, 0)
	if emojiResp.PackageVersions != nil {
		for _, pkg := range emojiResp.PackageVersions {
			emojiPackageVersions = append(emojiPackageVersions, types.EmojiPackageVersionItem{
				Id:      pkg.Id,
				Version: pkg.Version,
			})
		}
	}

	return &types.GetSyncEmojiPackagesRes{
		EmojiPackageVersions: emojiPackageVersions,
		ServerTimestamp:      time.Now().UnixMilli(),
	}, nil
}
