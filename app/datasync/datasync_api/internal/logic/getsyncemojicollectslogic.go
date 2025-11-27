package logic

import (
	"context"
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncEmojiCollectsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的收藏表情版本
func NewGetSyncEmojiCollectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncEmojiCollectsLogic {
	return &GetSyncEmojiCollectsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncEmojiCollectsLogic) GetSyncEmojiCollects(req *types.GetSyncEmojiCollectsReq) (resp *types.GetSyncEmojiCollectsRes, err error) {
	// 简化架构：datasync直接查询表情数据库，避免RPC调用
	// 注意：这会造成服务间耦合，但对于表情这种简单模块是合理的

	// 1. 获取用户收藏的表情数据（直接查询数据库）
	var emojiCollects []map[string]interface{}
	query := `
		SELECT uuid, version FROM emoji_collect_emoji
		WHERE user_id = ? AND updated_at > FROM_UNIXTIME(?)
	`
	err = l.svcCtx.DB.Raw(query, req.UserID, req.Since/1000).Scan(&emojiCollects).Error
	if err != nil {
		l.Errorf("查询用户收藏表情失败: userId=%s, error=%v", req.UserID, err)
		return nil, err
	}

	// 2. 获取用户收藏的表情包数据
	var emojiPackageCollects []map[string]interface{}
	query = `
		SELECT uuid, version FROM emoji_package_collect
		WHERE user_id = ? AND updated_at > FROM_UNIXTIME(?)
	`
	err = l.svcCtx.DB.Raw(query, req.UserID, req.Since/1000).Scan(&emojiPackageCollects).Error
	if err != nil {
		l.Errorf("查询用户收藏表情包失败: userId=%s, error=%v", req.UserID, err)
		return nil, err
	}

	// 检查表情包数据的版本变化（所有表情包，不仅仅是用户收藏的）
	var emojiPackageVersions []types.EmojiPackageVersionItem
	packageQuery := `SELECT id, version FROM emoji_package WHERE version > ?`
	var packages []struct {
		Id      uint32
		Version int64
	}
	l.svcCtx.DB.Raw(packageQuery, req.Since).Scan(&packages)

	for _, pkg := range packages {
		emojiPackageVersions = append(emojiPackageVersions, types.EmojiPackageVersionItem{
			Id:      pkg.Id,
			Version: pkg.Version,
		})
	}

	// 检查表情数据的版本变化（所有表情，不仅仅是用户收藏的）
	var emojiVersions []types.EmojiVersionItem
	emojiQuery := `SELECT id, version FROM emoji WHERE version > ?`
	var emojis []struct {
		Id      uint32
		Version int64
	}
	l.svcCtx.DB.Raw(emojiQuery, req.Since).Scan(&emojis)

	for _, emoji := range emojis {
		emojiVersions = append(emojiVersions, types.EmojiVersionItem{
			Id:      emoji.Id,
			Version: emoji.Version,
		})
	}

	// 检查表情包内容的版本变化（所有表情包内容）
	var emojiPackageContentVersions []types.EmojiPackageContentVersionItem
	contentQuery := `SELECT package_id, version FROM emoji_package_emoji WHERE version > ?`
	var contents []struct {
		PackageId uint32
		Version   int64
	}
	l.svcCtx.DB.Raw(contentQuery, req.Since).Scan(&contents)

	for _, content := range contents {
		emojiPackageContentVersions = append(emojiPackageContentVersions, types.EmojiPackageContentVersionItem{
			PackageId: content.PackageId,
			Version:   content.Version,
		})
	}

	l.Infof("用户表情同步：收藏表情=%d, 收藏表情包=%d",
		len(emojiCollects), len(emojiPackageCollects))

	// 转换为版本摘要格式
	emojiCollectVersions := make([]types.EmojiCollectVersionItem, 0)
	for _, item := range emojiCollects {
		emojiCollectVersions = append(emojiCollectVersions, types.EmojiCollectVersionItem{
			Id:      item["uuid"].(string),
			Version: item["version"].(int64),
		})
	}

	emojiPackageCollectVersions := make([]types.EmojiPackageCollectVersionItem, 0)
	for _, item := range emojiPackageCollects {
		emojiPackageCollectVersions = append(emojiPackageCollectVersions, types.EmojiPackageCollectVersionItem{
			Id:      item["uuid"].(string),
			Version: item["version"].(int64),
		})
	}

	return &types.GetSyncEmojiCollectsRes{
		EmojiCollectVersions:        emojiCollectVersions,
		EmojiPackageCollectVersions: emojiPackageCollectVersions,
		EmojiVersions:               emojiVersions,
		EmojiPackageVersions:        emojiPackageVersions,
		EmojiPackageContentVersions: emojiPackageContentVersions,
		ServerTimestamp:             time.Now().UnixMilli(),
	}, nil
}
