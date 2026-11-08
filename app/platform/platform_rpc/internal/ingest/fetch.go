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

// Package ingest 解析知识文档来源：doc 走 doc_rpc 读取本平台在线文档快照，file 由上传正文提供。
package ingest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	doccli "beaver/app/doc/doc_rpc/doc"
)

// Source 知识文档来源。Type 只有两种：doc（本平台在线文档）、file（上传文件）。
type Source struct {
	Type  string
	DocID string
}

// Result 来源正文快照。Revision 为 doc 的 version，用于自动同步增量判断。
type Result struct {
	Content  string
	Revision string
	Title    string
}

// Fetch 按来源类型读取正文。doc 通过 doc_rpc 按文档 ID 读取；file 需调用方直接提供 content_text。
func Fetch(ctx context.Context, docRpc doccli.Doc, src Source) (*Result, error) {
	switch strings.ToLower(strings.TrimSpace(src.Type)) {
	case "doc":
		docID := strings.TrimSpace(src.DocID)
		if docID == "" {
			return nil, errors.New("source_doc_id 不能为空")
		}
		res, err := docRpc.GetDocContent(ctx, &doccli.GetDocContentReq{DocId: docID})
		if err != nil {
			return nil, fmt.Errorf("读取在线文档失败: %w", err)
		}
		if res.Status != 1 {
			return nil, fmt.Errorf("在线文档已停用: %s", docID)
		}
		content := strings.TrimSpace(res.ContentText)
		if content == "" {
			return nil, fmt.Errorf("在线文档正文为空: %s", docID)
		}
		return &Result{
			Content:  content,
			Revision: strconv.Itoa(int(res.Version)),
			Title:    strings.TrimSpace(res.Title),
		}, nil
	case "file":
		return nil, errors.New("file 来源需由调用方直接传 content_text")
	default:
		return nil, fmt.Errorf("不支持的 source_type: %s", src.Type)
	}
}

// Digest 正文指纹，用于导入判重与自动同步增量判断。
func Digest(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
