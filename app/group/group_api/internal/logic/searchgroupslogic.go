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
	"strings"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type SearchGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 搜索群组
func NewSearchGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGroupsLogic {
	return &SearchGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("search_groups", ctx),
	}
}

func (l *SearchGroupsLogic) SearchGroups(req *types.GroupSearchReq) (resp *types.GroupSearchRes, err error) {
	resp = &types.GroupSearchRes{
		List: []types.GroupSearchItem{},
	}

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100 // 限制最大每页数量
	}
	offset := (page - 1) * limit

	// 构建查询条件
	query := l.svcCtx.DB.Model(&group_models.GroupModel{}).Where("status = 1") // 只搜索正常状态的群组

	// 如果有搜索关键词，按群组名称模糊搜索（大小写不敏感）
	if req.Keyword != "" {
		// 使用ILIKE进行大小写不敏感的搜索，如果数据库不支持，则回退到LIKE
		keyword := strings.TrimSpace(req.Keyword)
		if keyword != "" {
			// 尝试使用ILIKE（PostgreSQL风格的大小写不敏感搜索）
			query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+keyword+"%")
		}
	}

	// 获取总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询群组总数失败", Data: map[string]any{"keyword": req.Keyword, "err": err.Error()}})
		return nil, err
	}

	// 分页查询群组信息
	var groups []group_models.GroupModel
	err = query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&groups).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询群组列表失败", Data: map[string]any{"keyword": req.Keyword, "err": err.Error()}})
		return nil, err
	}

	// 为每个群组获取成员数量
	for _, group := range groups {
		var memberCount int64
		err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
			Where("group_id = ? AND status = 1", group.GroupID).
			Count(&memberCount).Error
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "查询群组成员数量失败", Data: map[string]any{"groupId": group.GroupID, "err": err.Error()}})
			memberCount = 0
		}

		resp.List = append(resp.List, types.GroupSearchItem{
			GroupID:     group.GroupID,
			Title:       group.Title,
			Avatar:      group.Avatar,
			MemberCount: int(memberCount),
			JoinType:    group.JoinType,
			CreatorID:   group.CreatorID,
		})
	}

	resp.Count = total

	l.logger.Info(model.LogMsg{Text: "搜索群组完成", Data: map[string]interface{}{"keyword": req.Keyword, "count": len(resp.List), "total": total}})
	return resp, nil
}
