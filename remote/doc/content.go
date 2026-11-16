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

package doc

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	doctype "beaver/types/remote/doc"
)

// GetDoc 读取在线文档元信息（标题 / 类型 / 同步状态等）。
func GetDoc(ctx context.Context, header http.Header, body *doctype.GetDocReq) (*doctype.GetDocRes, error) {
	return Post[doctype.GetDocRes](ctx, commonconfig.DocUrl+"/beaver-doc/api/doc/machine/v1/get_doc", header, body)
}

// GetDocContent 读取在线文档正文（知识导入 / 文档源校验）。
func GetDocContent(ctx context.Context, header http.Header, body *doctype.GetDocContentReq) (*doctype.GetDocContentRes, error) {
	return Post[doctype.GetDocContentRes](ctx, commonconfig.DocUrl+"/beaver-doc/api/doc/machine/v1/get_doc_content", header, body)
}
