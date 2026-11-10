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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/platform/platform_rpc/platform"

	"gorm.io/gorm"
)

// loadOwnedKnowledgeSource 先校验数字伙伴归属，再按 (source_id, digital_id) 取知识源。
// 不单独按 source_id 查：知识源是智能体的私有子实体，跨智能体不能互相读写。
func loadOwnedKnowledgeSource(db *gorm.DB, ownerID, digitalID, sourceID string) (*open_models.OpenKnowledgeSource, error) {
	if sourceID == "" {
		return nil, errors.New("knowledgeSourceId 不能为空")
	}
	if _, err := loadOwnedDigital(db, ownerID, digitalID); err != nil {
		return nil, err
	}
	var row open_models.OpenKnowledgeSource
	if err := db.Where("source_id = ? AND digital_id = ?", sourceID, digitalID).First(&row).Error; err != nil {
		return nil, errors.New("知识源不存在或无权限")
	}
	return &row, nil
}

// loadKnowledgeSourceIDs 把该数字伙伴所有知识源引入的知识库合并去重。
//
// 这是运行时检索列表的唯一来源：agent 侧只认 builtinConfig.knowledgeIds，
// 所有知识源变更后都会调用它重算，同一个知识库被多个知识源引入也只检索一次。
func loadKnowledgeSourceIDs(db *gorm.DB, digitalID string) []string {
	var rows []open_models.OpenKnowledgeSource
	if err := db.Where("digital_id = ?", digitalID).Find(&rows).Error; err != nil {
		return []string{}
	}
	seen := make(map[string]bool)
	ids := make([]string, 0)
	for _, row := range rows {
		for _, id := range row.KnowledgeIDs {
			id = strings.TrimSpace(id)
			if id != "" && !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
	}
	return ids
}

// syncDigitalKnowledgeIDs 知识源增删改后回写运行时扁平列表。
func syncDigitalKnowledgeIDs(db *gorm.DB, digitalID string) error {
	return writeDigitalKnowledgeIDs(db, digitalID, loadKnowledgeSourceIDs(db, digitalID))
}

// migrateLegacyKnowledgeSources 兼容旧版把知识源分组塞在 builtinConfig JSON 里的数据。
// 只在「这个数字伙伴还没有任何知识源实体」时执行，迁移完立刻清掉旧分组，
// 避免重复迁移出双份；只有扁平 knowledgeIds 的更老数据补一个默认知识源。
func migrateLegacyKnowledgeSources(db *gorm.DB, digitalID, ownerID string) error {
	var count int64
	if err := db.Model(&open_models.OpenKnowledgeSource{}).
		Where("digital_id = ?", digitalID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var cfg open_models.OpenDigitalConfig
	if err := db.Where("digital_id = ?", digitalID).First(&cfg).Error; err != nil {
		return nil
	}
	if cfg.BuiltinConfig == nil {
		return nil
	}
	legacy := cfg.BuiltinConfig.KnowledgeSources
	flatIDs := cfg.BuiltinConfig.KnowledgeIDs
	if len(legacy) == 0 && len(flatIDs) == 0 {
		return nil
	}

	rows := make([]open_models.OpenKnowledgeSource, 0, len(legacy)+1)
	for _, item := range legacy {
		id := strings.TrimSpace(item.ID)
		if id == "" {
			id = newKnowledgeSourceID()
		}
		name := strings.TrimSpace(item.Name)
		if name == "" {
			name = "默认知识源"
		}
		rows = append(rows, open_models.OpenKnowledgeSource{
			SourceID:      id,
			DigitalID:     digitalID,
			OwnerID:       ownerID,
			Name:          name,
			RetrievalTopK: open_models.DefaultKnowledgeSourceTopK,
			KnowledgeIDs:  open_models.KnowledgeIDList(uniqTrimmed(item.KnowledgeIDs)),
		})
	}
	if len(legacy) == 0 && len(flatIDs) > 0 {
		rows = append(rows, open_models.OpenKnowledgeSource{
			SourceID:      newKnowledgeSourceID(),
			DigitalID:     digitalID,
			OwnerID:       ownerID,
			Name:          "默认知识源",
			RetrievalTopK: open_models.DefaultKnowledgeSourceTopK,
			KnowledgeIDs:  open_models.KnowledgeIDList(uniqTrimmed(flatIDs)),
		})
	}
	if len(rows) == 0 {
		return nil
	}
	if err := db.Create(&rows).Error; err != nil {
		return err
	}
	return withBuiltinConfig(db, digitalID, func(bc *open_models.BuiltinConfig) {
		bc.KnowledgeSources = nil
		bc.KnowledgeIDs = loadKnowledgeSourceIDs(db, digitalID)
	})
}

func uniqTrimmed(ids []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

func newKnowledgeSourceID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("ks_%d", time.Now().UnixNano())
	}
	return "ks_" + hex.EncodeToString(buf)
}

// collectionNameMap 取平台知识库的 ID -> 名称，用于文档/召回结果标注来源。
func collectionNameMap(ctx context.Context, svcCtx *svc.ServiceContext) map[string]string {
	names := make(map[string]string)
	res, err := svcCtx.PlatformRpc.ListCollections(ctx, &platform.ListCollectionsReq{Page: 1, PageSize: 200})
	if err != nil {
		return names
	}
	for _, item := range res.Collections {
		names[item.CollectionId] = item.Name
	}
	return names
}

// countCollectionDocuments 累加若干知识库的文档数；单库查询失败时按 0 计，不影响整页。
func countCollectionDocuments(ctx context.Context, svcCtx *svc.ServiceContext, ids []string) int64 {
	var total int64
	for _, id := range ids {
		res, err := svcCtx.PlatformRpc.ListDocuments(ctx, &platform.ListDocumentsReq{
			CollectionId: id,
			Page:         1,
			PageSize:     1,
		})
		if err == nil {
			total += res.Total
		}
	}
	return total
}
