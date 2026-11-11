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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"beaver/app/open/open_models"

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

// sourceCollectionID 取知识源持有的 platform 索引集合 ID。
// 每行只放 1 个集合 ID（新建知识源时由本服务创建），取首个非空值即可定位集合。
// 调用方需自行处理「集合为空」的提示，不放这里以免掩盖业务含义。
func sourceCollectionID(ids open_models.KnowledgeIDList) string {
	for _, id := range ids {
		if s := strings.TrimSpace(id); s != "" {
			return s
		}
	}
	return ""
}

// parseBeaverDocID 从用户粘贴的海狸内部在线文档地址中宽容解析出 doc_id。
//
// 兼容形态：
//   - 裸 doc_id（如 doc_xxx）
//   - 完整 URL（https://doc.beaver.im/d/doc_xxx、/doc/、/#/doc/ 等）
//   - 带 query / hash 的地址
//
// 解析失败或为空返回 ""，由调用方计入失败列表。这里只负责「认出 doc_id」，
// 真正的鉴权与正文读取交给 doc_rpc / platform_rpc。
func parseBeaverDocID(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	// 先剥离 fragment（#/doc/xxx 这种 SPA 路由在 hash 里）
	if i := strings.Index(s, "#"); i >= 0 {
		s = s[i+1:]
	}
	// 有 scheme 的当成 URL，只保留 path + query
	if u, err := url.Parse(s); err == nil && u.Scheme != "" {
		s = u.Path
		if u.RawQuery != "" {
			s += "?" + u.RawQuery
		}
	}
	// /d/ 或 /doc/ 之后那一段就是 doc_id
	for _, prefix := range []string{"/d/", "/doc/"} {
		if idx := strings.Index(s, prefix); idx >= 0 {
			rest := s[idx+len(prefix):]
			if q := strings.IndexAny(rest, "/?"); q >= 0 {
				rest = rest[:q]
			}
			return strings.TrimSpace(rest)
		}
	}
	// 兜底：取最后一段 path，再去掉 query
	s = strings.TrimRight(s, "/")
	if i := strings.LastIndex(s, "/"); i >= 0 {
		s = s[i+1:]
	}
	if q := strings.Index(s, "?"); q >= 0 {
		s = s[:q]
	}
	return strings.TrimSpace(s)
}

func newKnowledgeSourceID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("ks_%d", time.Now().UnixNano())
	}
	return "ks_" + hex.EncodeToString(buf)
}
