package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncEmojiCollectsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户表情收藏的版本信息
func NewGetSyncEmojiCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojiCollectsLogic {
	return &GetSyncEmojiCollectsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncEmojiCollectsLogic) GetSyncEmojiCollects(req *types.GetSyncEmojiCollectsReq) (resp *types.GetSyncEmojiCollectsRes, err error) {
	// 初始化结果变量
	emojiCollectVersions := make([]types.EmojiCollectVersionItem, 0)
	emojiPackageCollectVersions := make([]types.EmojiPackageCollectVersionItem, 0)
	emojiPackageVersions := make([]types.EmojiPackageVersionItem, 0)
	emojiPackageContentVersions := make([]types.EmojiPackageContentVersionItem, 0)

	// 1. 获取用户收藏的表情版本信息
	emojiCollectResp, err := l.svcCtx.EmojiRpc.GetUserEmojiCollects(l.ctx, &emoji_rpc.GetUserEmojiCollectsReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取用户收藏表情版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	// 转换用户收藏的表情数据（使用收藏记录ID）
	if emojiCollectResp.EmojiCollectVersions != nil {
		for _, item := range emojiCollectResp.EmojiCollectVersions {
			emojiCollectVersions = append(emojiCollectVersions, types.EmojiCollectVersionItem{
				EmojiCollectId: item.EmojiCollectId,
				Version:        item.Version,
			})
		}
	}

	// 2. 获取用户收藏的表情包版本信息
	emojiPackageCollectResp, err := l.svcCtx.EmojiRpc.GetUserEmojiPackageCollects(l.ctx, &emoji_rpc.GetUserEmojiPackageCollectsReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取用户收藏表情包版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	// 转换用户收藏的表情包数据（使用收藏记录ID）
	if emojiPackageCollectResp.EmojiPackageCollectVersions != nil {
		for _, item := range emojiPackageCollectResp.EmojiPackageCollectVersions {
			emojiPackageCollectVersions = append(emojiPackageCollectVersions, types.EmojiPackageCollectVersionItem{
				PackageCollectId: item.PackageCollectId,
				Version:          item.Version,
			})
		}
	}

	// 3. 获取表情包基础数据版本信息
	emojiPackagesResp, err := l.svcCtx.EmojiRpc.GetEmojiPackages(l.ctx, &emoji_rpc.GetEmojiPackagesReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取表情包基础数据版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	// 转换表情包数据
	if emojiPackagesResp.EmojiPackageVersions != nil {
		for _, item := range emojiPackagesResp.EmojiPackageVersions {
			emojiPackageVersions = append(emojiPackageVersions, types.EmojiPackageVersionItem{
				PackageId: item.PackageId,
				Version:   item.Version,
			})
		}
	}

	// 4. 获取表情包内容版本信息
	emojiPackageContentsResp, err := l.svcCtx.EmojiRpc.GetEmojiPackageContents(l.ctx, &emoji_rpc.GetEmojiPackageContentsReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取表情包内容版本信息失败: userId=%s, since=%d, error=%v", req.UserID, req.Since, err)
		return nil, err
	}

	// 转换表情包内容数据
	if emojiPackageContentsResp.EmojiPackageContentVersions != nil {
		for _, item := range emojiPackageContentsResp.EmojiPackageContentVersions {
			emojiPackageContentVersions = append(emojiPackageContentVersions, types.EmojiPackageContentVersionItem{
				PackageId: item.PackageId,
				Version:   item.Version,
			})
		}
	}

	l.Infof("用户 %s 表情收藏同步完成: 收藏表情=%d, 收藏表情包=%d, 表情包=%d, 表情包内容=%d",
		req.UserID, len(emojiCollectVersions), len(emojiPackageCollectVersions),
		len(emojiPackageVersions), len(emojiPackageContentVersions))

	return &types.GetSyncEmojiCollectsRes{
		EmojiCollectVersions:        emojiCollectVersions,
		EmojiPackageCollectVersions: emojiPackageCollectVersions,
		EmojiPackageVersions:        emojiPackageVersions,
		EmojiPackageContentVersions: emojiPackageContentVersions,
		ServerTimestamp:             time.Now().UnixMilli(),
	}, nil
}
