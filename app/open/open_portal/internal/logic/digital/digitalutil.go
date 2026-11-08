/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func maskSecret(secret string) string {
	if secret == "" {
		return ""
	}
	if len(secret) <= 16 {
		return "****"
	}
	return secret[:8] + "****" + secret[len(secret)-8:]
}

func loadOwnedDigital(db *gorm.DB, ownerID, digitalID string) (*open_models.OpenDigital, error) {
	if digitalID == "" {
		return nil, errors.New("digitalId 不能为空")
	}
	var partner open_models.OpenDigital
	if err := db.Where("digital_id = ? AND owner_id = ?", digitalID, ownerID).First(&partner).Error; err != nil {
		return nil, errors.New("数字伙伴不存在或无权限")
	}
	return &partner, nil
}

func loadDigitalAdminIDs(db *gorm.DB, digitalID string) []string {
	return loadDigitalMemberIDs(db, digitalID, open_models.DigitalMemberRoleAdmin)
}

func loadDigitalViewerIDs(db *gorm.DB, digitalID string) []string {
	return loadDigitalMemberIDs(db, digitalID, open_models.DigitalMemberRoleViewer)
}

func loadDigitalMemberIDs(db *gorm.DB, digitalID, role string) []string {
	var rows []open_models.OpenDigitalMember
	_ = db.Where("digital_id = ? AND role = ?", digitalID, role).Find(&rows).Error
	ids := make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.UserID)
	}
	return ids
}

// syncDigitalMembers diff 增删指定 role 的成员：只删移除的、只加新增的，未变更的不动。
// ownerID 不入成员表（Owner 天然管理权）；空 ID 跳过。
func syncDigitalMembers(tx *gorm.DB, digitalID, role string, wantIDs []string, ownerID string) error {
	want := make(map[string]struct{}, len(wantIDs))
	for _, uid := range wantIDs {
		if uid == "" || uid == ownerID {
			continue
		}
		want[uid] = struct{}{}
	}

	var existing []open_models.OpenDigitalMember
	if err := tx.Where("digital_id = ? AND role = ?", digitalID, role).Find(&existing).Error; err != nil {
		return err
	}
	have := make(map[string]struct{}, len(existing))
	for _, r := range existing {
		have[r.UserID] = struct{}{}
	}

	// 删除：有现无
	for uid := range have {
		if _, ok := want[uid]; !ok {
			if err := tx.Where("digital_id = ? AND user_id = ? AND role = ?", digitalID, uid, role).
				Delete(&open_models.OpenDigitalMember{}).Error; err != nil {
				return err
			}
		}
	}
	// 新增：现有无
	for uid := range want {
		if _, ok := have[uid]; !ok {
			if err := tx.Create(&open_models.OpenDigitalMember{DigitalID: digitalID, UserID: uid, Role: role}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func fillDigitalInfo(ctx context.Context, svcCtx *svc.ServiceContext, partner *open_models.OpenDigital) types.DigitalInfo {
	tools, knowledge, mcp := loadDigitalConfigCounts(svcCtx.DB, partner.DigitalID)
	info := types.DigitalInfo{
		DigitalID:      partner.DigitalID,
		WorkID:         partner.WorkID,
		Runtime:        partner.Runtime,
		Visibility:     partner.Visibility,
		Status:         partner.Status,
		AdminIDs:       loadDigitalAdminIDs(svcCtx.DB, partner.DigitalID),
		ViewerIDs:      loadDigitalViewerIDs(svcCtx.DB, partner.DigitalID),
		ToolCount:      tools,
		KnowledgeCount: knowledge,
		McpCount:       mcp,
		CreatedAt:      partner.CreatedAt.Unix(),
		UpdatedAt:      partner.UpdatedAt.Unix(),
	}
	if partner.Runtime == open_models.DigitalRuntimeExternal {
		info.AppKey = partner.AppKey
		info.AppSecret = maskSecret(partner.AppSecret)
	}
	if userRes, err := svcCtx.UserRpc.UserInfo(ctx, &user_rpc.UserInfoReq{UserID: partner.DigitalID}); err == nil && userRes.UserInfo != nil {
		info.Name = userRes.UserInfo.NickName
		info.Avatar = userRes.UserInfo.Avatar
		info.Description = userRes.UserInfo.Abstract
	}
	return info
}

// loadDigitalConfigCounts 取编排配置里引用的工具 / 知识库 / MCP 数量。
// 只有 builtin（提示词编排）会在平台侧挂这些能力，其余 runtime 不统计，返回 0。
func loadDigitalConfigCounts(db *gorm.DB, digitalID string) (tools, knowledge, mcp int) {
	var cfg open_models.OpenDigitalConfig
	// 用 Find 而不是 First：还没配过的数字伙伴没有配置行，First 会刷一条 not found 日志
	if err := db.Where("digital_id = ?", digitalID).Limit(1).Find(&cfg).Error; err != nil {
		return 0, 0, 0
	}
	if cfg.BuiltinConfig == nil {
		return 0, 0, 0
	}
	return len(cfg.BuiltinConfig.Tools), len(cfg.BuiltinConfig.KnowledgeIDs), len(cfg.BuiltinConfig.McpServerIDs)
}

func generateExternalCredentials() (appKey, appSecret string) {
	appKey = fmt.Sprintf("dig_%s", uuid.New().String()[:8])
	appSecret = uuid.New().String() + uuid.New().String()
	return appKey, appSecret
}

// applyWorkID 分配一个未被占用的工号：DG + 6 位数字，冲突重试。
func applyWorkID(db *gorm.DB) (string, error) {
	for i := 0; i < 30; i++ {
		candidate := fmt.Sprintf("DG%06d", rand.Intn(1000000))
		var count int64
		if err := db.Model(&open_models.OpenDigital{}).Where("work_id = ?", candidate).Count(&count).Error; err != nil {
			return "", errors.New("工号生成失败，请重试")
		}
		if count == 0 {
			return candidate, nil
		}
	}
	return "", errors.New("工号生成失败，请重试")
}

// checkWorkIDAvailable 自定义工号时校验唯一：排除自己，被别人占用则拒绝。
func checkWorkIDAvailable(db *gorm.DB, workID, selfDigitalID string) error {
	var count int64
	if err := db.Model(&open_models.OpenDigital{}).
		Where("work_id = ? AND digital_id <> ?", workID, selfDigitalID).
		Count(&count).Error; err != nil {
		return errors.New("工号校验失败，请重试")
	}
	if count > 0 {
		return errors.New("该工号已被占用")
	}
	return nil
}
