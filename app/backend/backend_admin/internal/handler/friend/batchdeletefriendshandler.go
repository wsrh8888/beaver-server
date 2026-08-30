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

package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/friend"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BatchDeleteFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BatchDeleteFriendsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if len(req.Ids) == 0 {
			response.Response(r, w, nil, errors.New("删除的好友关系ID列表不能为空"))
			return
		}

		// 验证ID格式
		for _, idStr := range req.Ids {
			_, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				response.Response(r, w, nil, errors.New("无效的好友关系ID: "+idStr))
				return
			}
		}

		l := logic.NewBatchDeleteFriendsLogic(r.Context(), svcCtx)
		resp, err := l.BatchDeleteFriends(&req)
		response.Response(r, w, resp, err)
	}
}
