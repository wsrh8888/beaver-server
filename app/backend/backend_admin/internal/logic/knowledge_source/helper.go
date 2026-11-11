/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 *
 * 场景知识源的公共工具：归属校验、索引集合定位、运行时列表回写、在线文档地址解析。
 */

package knowledge_source

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_models"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

// sceneMaxLen 场景标识最大长度，与 agent_apps.scene 列宽保持一致。
const sceneMaxLen = 32

// requireScene 场景是知识源唯一的归属维度，空值一律拒绝。
func requireScene(scene string) (string, error) {
	s := strings.TrimSpace(scene)
	if s == "" {
		return "", errors.New("scene 不能为空")
	}
	if len(s) > sceneMaxLen {
		return "", errors.New("scene 不合法")
	}
	return s, nil
}

// loadSource 按 (source_id, scene) 取知识源。不单独按 source_id 查：
// 知识源是场景的私有子实体，跨场景不能互相读写。
func loadSource(db *gorm.DB, scene, sourceID string) (*backend_models.AgentKnowledgeSource, error) {
	if sourceID = strings.TrimSpace(sourceID); sourceID == "" {
		return nil, errors.New("knowledgeSourceId 不能为空")
	}
	var row backend_models.AgentKnowledgeSource
	if err := db.Where("source_id = ? AND scene = ?", sourceID, scene).First(&row).Error; err != nil {
		return nil, errors.New("知识源不存在")
	}
	return &row, nil
}

// sourceCollectionID 取知识源持有的 platform 索引集合 ID。
// 每行只放 1 个集合 ID（新建知识源时由本服务创建），取首个非空值即可定位集合。
func sourceCollectionID(ids backend_models.KnowledgeIDList) string {
	for _, id := range ids {
		if s := strings.TrimSpace(id); s != "" {
			return s
		}
	}
	return ""
}

// loadSceneKnowledgeIDs 把该场景所有知识源引入的索引集合合并去重。
//
// 这是运行时检索列表的唯一来源：agent 侧只认 agent_apps.knowledge_ids，
// 知识源任何增删改后都调用它重算，同一个集合被多个知识源引入也只检索一次。
func loadSceneKnowledgeIDs(db *gorm.DB, scene string) []string {
	var rows []backend_models.AgentKnowledgeSource
	if err := db.Where("scene = ?", scene).Find(&rows).Error; err != nil {
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

// syncSceneKnowledgeIDs 知识源增删改后回写运行时扁平列表。
// agent_apps 里没有该场景记录时补建一条（只带 scene 与 knowledge_ids，其余字段留空）。
func syncSceneKnowledgeIDs(ctx context.Context, svcCtx *svc.ServiceContext, logger *beaverlog.Logger, scene string) error {
	payload, marshalErr := json.Marshal(loadSceneKnowledgeIDs(svcCtx.DB, scene))
	if marshalErr != nil {
		return marshalErr
	}
	listRes, listErr := svcCtx.PlatformRpc.ListAgentApps(ctx, &platform.ListAgentAppsReq{
		Scene:    scene,
		Page:     1,
		PageSize: 20,
	})
	if listErr != nil {
		logger.Error(model.LogMsg{Text: "查询场景应用失败", Data: map[string]interface{}{"scene": scene, "err": listErr.Error()}})
		return errors.New("同步检索配置失败")
	}
	if len(listRes.Apps) == 0 {
		if _, createErr := svcCtx.PlatformRpc.CreateAgentApp(ctx, &platform.CreateAgentAppReq{
			Name:         scene,
			Scene:        scene,
			KnowledgeIds: string(payload),
			Status:       1,
		}); createErr != nil {
			logger.Error(model.LogMsg{Text: "补建场景应用失败", Data: map[string]interface{}{"scene": scene, "err": createErr.Error()}})
			return errors.New("同步检索配置失败")
		}
		return nil
	}
	if _, updateErr := svcCtx.PlatformRpc.UpdateAgentApp(ctx, &platform.UpdateAgentAppReq{
		AppId:        listRes.Apps[0].AppId,
		KnowledgeIds: string(payload),
	}); updateErr != nil {
		logger.Error(model.LogMsg{Text: "回写场景知识列表失败", Data: map[string]interface{}{"scene": scene, "err": updateErr.Error()}})
		return errors.New("同步检索配置失败")
	}
	return nil
}

// collectionDocCount 统计知识源索引集合下的文档数；集合为空时返回 0。
func collectionDocCount(ctx context.Context, svcCtx *svc.ServiceContext, logger *beaverlog.Logger, collectionID string) int64 {
	if collectionID == "" {
		return 0
	}
	res, rpcErr := svcCtx.PlatformRpc.ListDocuments(ctx, &platform.ListDocumentsReq{
		CollectionId: collectionID,
		Page:         1,
		PageSize:     1,
	})
	if rpcErr != nil {
		logger.Error(model.LogMsg{Text: "统计知识源文档数失败", Data: map[string]interface{}{"err": rpcErr.Error(), "collectionId": collectionID}})
		return 0
	}
	return res.Total
}

// parseBeaverDocID 从粘贴的海狸内部在线文档地址中宽容解析出 doc_id。
//
// 兼容形态：裸 doc_id、完整 URL（/d/xxx、/doc/xxx、#/doc/xxx）、带 query 的地址。
// 解析失败或为空返回 ""，由调用方计入失败列表。
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

// parseRFC3339Unix platform 侧时间以 RFC3339 字符串返回，前端统一用秒级时间戳。
func parseRFC3339Unix(raw string) int64 {
	if strings.TrimSpace(raw) == "" {
		return 0
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Unix()
	}
	return 0
}
