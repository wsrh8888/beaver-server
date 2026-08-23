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
	"time"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/types/file_rpc"
)

func toFileItem(f file_models.FileModel) *file_rpc.FileItem {
	return &file_rpc.FileItem{
		Id:           uint64(f.Id),
		FileKey:      f.FileKey,
		OriginalName: f.OriginalName,
		Size:         f.Size,
		Path:         f.Path,
		Md5:          f.Md5,
		Type:         f.Type,
		Source:       string(f.Source),
		CreatedAt:    time.Time(f.CreatedAt).Format(time.RFC3339),
		UpdatedAt:    time.Time(f.UpdatedAt).Format(time.RFC3339),
	}
}
