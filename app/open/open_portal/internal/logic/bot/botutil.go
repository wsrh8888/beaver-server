/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package bot

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"gorm.io/gorm"
)

func loadOwnedBot(db *gorm.DB, ownerID, botID string) (*open_models.OpenBot, error) {
	if botID == "" {
		return nil, errors.New("botId 不能为空")
	}
	var row open_models.OpenBot
	if err := db.Where("bot_id = ? AND owner_id = ?", botID, ownerID).First(&row).Error; err != nil {
		return nil, errors.New("机器人不存在或无权限")
	}
	return &row, nil
}

func buildBotWebhookURL(apiBase, token string) string {
	base := strings.TrimSuffix(apiBase, "/")
	return fmt.Sprintf("%s/api/open/bot_public/v1/send?token=%s", base, token)
}

func loadBotAdminIDs(db *gorm.DB, botID string) []string {
	return loadBotMemberIDs(db, botID, open_models.BotMemberRoleAdmin)
}

func loadBotViewerIDs(db *gorm.DB, botID string) []string {
	return loadBotMemberIDs(db, botID, open_models.BotMemberRoleViewer)
}

func loadBotMemberIDs(db *gorm.DB, botID, role string) []string {
	var rows []open_models.OpenBotMember
	_ = db.Where("bot_id = ? AND role = ?", botID, role).Find(&rows).Error
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.UserID)
	}
	return ids
}

// syncBotMembers diff 增删指定 role 的成员：只删移除的、只加新增的，未变更的不动。
// ownerID 不入成员表（Owner 天然管理权）；空 ID 跳过。
func syncBotMembers(tx *gorm.DB, botID, role string, wantIDs []string, ownerID string) error {
	want := make(map[string]struct{}, len(wantIDs))
	for _, uid := range wantIDs {
		if uid == "" || uid == ownerID {
			continue
		}
		want[uid] = struct{}{}
	}

	var existing []open_models.OpenBotMember
	if err := tx.Where("bot_id = ? AND role = ?", botID, role).Find(&existing).Error; err != nil {
		return err
	}
	have := make(map[string]struct{}, len(existing))
	for _, r := range existing {
		have[r.UserID] = struct{}{}
	}

	// 删除：有现无
	for uid := range have {
		if _, ok := want[uid]; !ok {
			if err := tx.Where("bot_id = ? AND user_id = ? AND role = ?", botID, uid, role).
				Delete(&open_models.OpenBotMember{}).Error; err != nil {
				return err
			}
		}
	}
	// 新增：现有无
	for uid := range want {
		if _, ok := have[uid]; !ok {
			if err := tx.Create(&open_models.OpenBotMember{BotID: botID, UserID: uid, Role: role}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func fillBotInfo(ctx context.Context, svcCtx *svc.ServiceContext, row *open_models.OpenBot) types.BotInfo {
	info := types.BotInfo{
		BotID:      row.BotID,
		Visibility: row.Visibility,
		WebhookURL: buildBotWebhookURL(svcCtx.Config.Domain, "***"),
		AdminIDs:   loadBotAdminIDs(svcCtx.DB, row.BotID),
		ViewerIDs:  loadBotViewerIDs(svcCtx.DB, row.BotID),
		CreatedAt:  row.CreatedAt.Unix(),
	}
	if userRes, err := svcCtx.UserRpc.UserInfo(ctx, &user_rpc.UserInfoReq{UserID: row.BotID}); err == nil && userRes.UserInfo != nil {
		info.Name = userRes.UserInfo.NickName
		info.Avatar = userRes.UserInfo.Avatar
		info.Description = userRes.UserInfo.Abstract
	}
	return info
}
