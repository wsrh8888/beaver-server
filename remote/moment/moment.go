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

package moment

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	momenttype "beaver/types/remote/moment"
)

// ListMoments 列圈子动态（管理后台 / 运营概览）。
func ListMomentsApi(ctx context.Context, header http.Header, body *momenttype.ListMomentsReq) (*momenttype.ListMomentsRes, error) {
	return Post[momenttype.ListMomentsRes](ctx, commonconfig.SocialUrl+"/beaver-social/api/moment/machine/v1/list_moments", header, body)
}

// UpdateMoment 更新圈子动态（状态 / 内容变更）。
func UpdateMomentApi(ctx context.Context, header http.Header, body *momenttype.UpdateMomentReq) (*momenttype.UpdateMomentRes, error) {
	return Post[momenttype.UpdateMomentRes](ctx, commonconfig.SocialUrl+"/beaver-social/api/moment/machine/v1/update_moment", header, body)
}

// ListMomentComments 列动态下的评论。
func ListMomentCommentsApi(ctx context.Context, header http.Header, body *momenttype.ListMomentCommentsReq) (*momenttype.ListMomentCommentsRes, error) {
	return Post[momenttype.ListMomentCommentsRes](ctx, commonconfig.SocialUrl+"/beaver-social/api/moment/machine/v1/list_comments", header, body)
}

// DeleteMomentComments 批量删除动态下的评论。
func DeleteMomentCommentsApi(ctx context.Context, header http.Header, body *momenttype.DeleteMomentCommentsReq) (*momenttype.DeleteMomentCommentsRes, error) {
	return Post[momenttype.DeleteMomentCommentsRes](ctx, commonconfig.SocialUrl+"/beaver-social/api/moment/machine/v1/delete_comments", header, body)
}
