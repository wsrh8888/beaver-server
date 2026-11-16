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

package open

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	opentype "beaver/types/remote/open"
)

// ListDigitals 列数字伙伴（数字人）。
func ListDigitals(ctx context.Context, header http.Header, body *opentype.ListDigitalsReq) (*opentype.ListDigitalsRes, error) {
	return Post[opentype.ListDigitalsRes](ctx, commonconfig.OpenUrl+"/beaver-open/api/open/machine/v1/list_digitals", header, body)
}

// CreateDigital 创建数字伙伴。
func CreateDigital(ctx context.Context, header http.Header, body *opentype.CreateDigitalReq) (*opentype.CreateDigitalRes, error) {
	return Post[opentype.CreateDigitalRes](ctx, commonconfig.OpenUrl+"/beaver-open/api/open/machine/v1/create_digital", header, body)
}

// UpdateDigital 更新数字伙伴。
func UpdateDigital(ctx context.Context, header http.Header, body *opentype.UpdateDigitalReq) (*opentype.UpdateDigitalRes, error) {
	return Post[opentype.UpdateDigitalRes](ctx, commonconfig.OpenUrl+"/beaver-open/api/open/machine/v1/update_digital", header, body)
}
