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

package moderation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)


type HandleModerationCaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewHandleModerationCaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleModerationCaseLogic {
	return &HandleModerationCaseLogic{logger: logger.New("handle_moderation_case"), ctx: ctx, svcCtx: svcCtx}
}

func (l *HandleModerationCaseLogic) HandleModerationCase(req *types.HandleModerationCaseReq) (resp *types.HandleModerationCaseRes, err error) {
	if req.CaseID == 0 {
		return nil, errors.New("工单ID不能为空")
	}
	if req.Status < backend_models.CaseStatusPending || req.Status > backend_models.CaseStatusRejected {
		return nil, errors.New("无效的工单状态")
	}

	var c backend_models.AdminModerationCase
	if err = l.svcCtx.DB.Where("id = ?", req.CaseID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		logx.WithContext(l.ctx).Errorf("查询工单失败: %v", err)
		return nil, err
	}

	for _, act := range req.Actions {
		if actErr := executeControlAction(l.ctx, l.svcCtx, req.UserID, uint64(c.Id), act); actErr != nil {
			logx.WithContext(l.ctx).Errorf("执行管控动作失败 action=%s: %v", act.Action, actErr)
			return nil, actErr
		}
	}

	actionsJSON := ""
	if len(req.Actions) > 0 {
		if b, mErr := json.Marshal(req.Actions); mErr == nil {
			actionsJSON = string(b)
		}
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":        req.Status,
		"handler_id":    req.UserID,
		"handle_remark": req.HandleRemark,
		"handle_time":   &now,
		"actions_taken": actionsJSON,
	}
	if err = l.svcCtx.DB.Model(&c).Updates(updates).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("更新工单失败: %v", err)
		return nil, err
	}

	if req.Status == backend_models.CaseStatusResolved || req.Status == backend_models.CaseStatusRejected {
		reportRes, listErr := l.svcCtx.PlatformRpc.ListContentReports(l.ctx, &platform_rpc.ListContentReportsReq{
			Page:       1,
			PageSize:   100,
			TargetType: int32(c.TargetType),
			TargetId:   c.TargetID,
		})
		if listErr == nil && reportRes != nil {
			ids := make([]uint64, 0)
			for _, r := range reportRes.List {
				if r.CaseId == uint64(c.Id) || r.CaseId == 0 {
					ids = append(ids, r.Id)
				}
			}
			if len(ids) > 0 {
				action := int32(3)
				if req.Status == backend_models.CaseStatusRejected {
					action = 2
				}
				_, _ = l.svcCtx.PlatformRpc.UpdateContentReports(l.ctx, &platform_rpc.UpdateContentReportsReq{
					Ids:          ids,
					Action:       action,
					HandlerId:    req.UserID,
					HandleRemark: req.HandleRemark,
				})
			}
		}
	}

	l.svcCtx.RecordOperation(req.UserID, "handle_case", "case", fmt.Sprintf("%d", c.Id), uint64(c.Id),
		fmt.Sprintf("处置工单 status=%d remark=%s", req.Status, req.HandleRemark), "success", "")

	l.logger.Info(model.LogMsg{
		Text: "审核工单处置成功",
		Data: map[string]interface{}{
			"caseId":     req.CaseID,
			"operatorId": req.UserID,
			"status":     req.Status,
		},
	})

	return &types.HandleModerationCaseRes{}, nil
}
