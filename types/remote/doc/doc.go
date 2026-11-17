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

// beaver-server-online S2S 出入参类型
//
// 手写 JSON DTO：字段名与 json tag 与对端 proto 的 JSON 编码逐字段一致，
// 不依赖 protobuf 运行时（无 protoimpl / protoreflect / state）。
package doc

// GetDocContentApi 取文档内容：Req + Res（本接口无专属子类型）。
type (
	GetDocContentReq struct {
		DocId string `json:"doc_id,omitempty"`
	}

	GetDocContentRes struct {
		DocId       string `json:"doc_id,omitempty"`
		Title       string `json:"title,omitempty"`
		ContentText string `json:"content_text,omitempty"`
		Version     int32  `json:"version,omitempty"`
		Status      int32  `json:"status,omitempty"`
	}
)

// GetDocApi 取文档元数据：Req + Res + 本接口专属子类型。
type (
	GetDocReq struct {
		DocId  string `json:"doc_id,omitempty"`
		UserId string `json:"user_id,omitempty"`
	}

	GetDocRes struct {
		Doc *GetDocDocItem `json:"doc,omitempty"`
	}

	// GetDocDocItem 文档元数据项。
	GetDocDocItem struct {
		Id           uint64 `json:"id,omitempty"`
		DocId        string `json:"doc_id,omitempty"`
		NodeId       string `json:"node_id,omitempty"`
		SpaceId      string `json:"space_id,omitempty"`
		Title        string `json:"title,omitempty"`
		ContentJson  string `json:"content_json,omitempty"`
		Version      int32  `json:"version,omitempty"`
		OwnerUserId  string `json:"owner_user_id,omitempty"`
		LastEditedBy string `json:"last_edited_by,omitempty"`
		Status       int32  `json:"status,omitempty"`
		UpdatedAt    string `json:"updated_at,omitempty"`
		CreatedAt    string `json:"created_at,omitempty"`
		Epoch        int32  `json:"epoch,omitempty"`
		DocType      string `json:"doc_type,omitempty"`
		Icon         string `json:"icon,omitempty"`
		Cover        string `json:"cover,omitempty"`
		RagEnabled   bool   `json:"rag_enabled,omitempty"`
		ContentText  string `json:"content_text,omitempty"`
	}
)
