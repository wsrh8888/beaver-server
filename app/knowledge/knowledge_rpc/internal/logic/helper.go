/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package logic

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	knowledge_models "beaver/app/knowledge/knowledge_models"
	"beaver/app/knowledge/knowledge_rpc/internal/svc"
	"beaver/app/knowledge/knowledge_rpc/types/knowledge_rpc"
	platformcli "beaver/app/platform/platform_rpc/platform"

	"github.com/zeromicro/go-zero/core/logx"
)

// zeroDeletedAt 软删哨兵：deleted_at 等于该时刻表示未删，与 knowledge 模型的默认值对齐。
const zeroDeletedAt = "1970-01-01 00:00:00"

// randHex 业务 ID 随机段，与 platform 域同法（crypto/rand）。
func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// jsonStrList 把 AgentApp.KnowledgeIDs / index_config 里的 JSON 数组解析成 []string。
// 空串或坏 JSON 返回 nil：绑定视为「未配」，与运行时降级语义一致，不报错。
func jsonStrList(raw string) []string {
	raw = trim(raw)
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	filtered := out[:0]
	for _, s := range out {
		if trim(s) != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// trim 去首尾空白，业务上判断「是否配置了」用它。
func trim(s string) string {
	return strings.TrimSpace(s)
}

func toCollectionItem(row *knowledge_models.KnowledgeCollection) *knowledge_rpc.KnowledgeCollectionItem {
	if row == nil {
		return nil
	}
	return &knowledge_rpc.KnowledgeCollectionItem{
		Id:           uint64(row.Id),
		CollectionId: row.CollectionID,
		ParentId:     row.ParentID,
		Name:         row.Name,
		Kind:         row.Kind,
		StorageRef:   row.StorageRef,
		AclMode:      row.AclMode,
		IndexConfig:  row.IndexConfig,
		Description:  row.Description,
		OwnerUserId:  row.OwnerUserID,
		Status:       int32(row.Status),
		Sort:         int32(row.Sort),
	}
}

func toDocumentItem(row *knowledge_models.KnowledgeDocument) *knowledge_rpc.KnowledgeDocumentItem {
	if row == nil {
		return nil
	}
	item := &knowledge_rpc.KnowledgeDocumentItem{
		Id:              uint64(row.Id),
		DocumentId:      row.DocumentID,
		CollectionId:    row.CollectionID,
		Title:           row.Title,
		MediaType:       row.MediaType,
		DocType:         row.DocType,
		ByteSize:        row.ByteSize,
		SourceRef:       row.SourceRef,
		MetadataText:    row.MetadataText,
		RagEnabled:      row.RagEnabled,
		IndexStatus:     row.IndexStatus,
		IndexRetryCount: int32(row.IndexRetryCount),
		IndexError:      row.IndexError,
		Version:         int32(row.Version),
		Status:          int32(row.Status),
		OwnerUserId:     row.OwnerUserID,
		UpdatedBy:       row.UpdatedBy,
	}
	if row.IndexedAt != nil {
		item.IndexedAt = row.IndexedAt.Format(time.RFC3339)
	}
	return item
}

// resolveBoundCollectionIDs 解析 app 的知识库绑定：只认 AgentApp.KnowledgeIDs，
// 不接受调用方指定分区（防越权：绕过配置读未绑库）。经 PlatformRpc 回查，跨域不直连 platform 库。
// 无 app_id/scene、查不到、或绑定为空一律返回 nil（调用方据此返回空命中）。
func resolveBoundCollectionIDs(ctx context.Context, svcCtx *svc.ServiceContext, appID, scene string) []string {
	if trim(appID) == "" && trim(scene) == "" {
		return nil
	}
	req := &platformcli.GetAgentAppReq{AppId: appID, Scene: scene}
	res, err := svcCtx.PlatformRpc.GetAgentApp(ctx, req)
	if err != nil {
		logx.Errorf("解析 AgentApp 知识绑定失败 app=%s scene=%s: %v", appID, scene, err)
		return nil
	}
	app := res.GetApp()
	if app == nil || app.GetStatus() != 1 {
		return nil
	}
	return jsonStrList(app.GetKnowledgeIds())
}

// clearDocIndex 同步清掉某文档在 OpenSearch 的分段与 DB chunk 行。
// 停用、删除、关闭检索等「让检索立即不再命中」的操作走这里，不等 worker——
// 检索侧已按 rag_enabled/status 过滤，但立即清理能避免旧分段在边界条件下被读到。
func clearDocIndex(ctx context.Context, svcCtx *svc.ServiceContext, documentIDs []string) {
	if len(documentIDs) == 0 {
		return
	}
	if svcCtx.OpenSearch != nil {
		if err := svcCtx.OpenSearch.DeleteByDocumentIDs(ctx, svcCtx.ChunkIndex, documentIDs); err != nil {
			logx.Errorf("清理文档分段索引失败 %v: %v", documentIDs, err)
		}
	}
	if err := svcCtx.DB.Where("document_id IN ?", documentIDs).Delete(&knowledge_models.KnowledgeChunk{}).Error; err != nil {
		logx.Errorf("清理文档分段行失败 %v: %v", documentIDs, err)
	}
}

// nowPtr 返回当前时间指针，写入可空 *time.Time 列。
func nowPtr() *time.Time {
	n := time.Now()
	return &n
}
