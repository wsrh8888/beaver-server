/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type DeleteDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDigitalLogic {
	return &DeleteDigitalLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_digital", ctx)}
}

// DeleteDigital 删除数字伙伴：主表、编排配置、成员、令牌一并软删除，IM 用户同步软删除。
// 会话与历史消息保留（人已删，历史仍可读），重建同名数字伙伴会拿到新的 DigitalID。
func (l *DeleteDigitalLogic) DeleteDigital(req *types.DeleteDigitalReq) (resp *types.DeleteDigitalRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("digital_id = ?", partner.DigitalID).Delete(&open_models.OpenDigitalConfig{}).Error; err != nil {
			return err
		}
		if err := tx.Where("digital_id = ?", partner.DigitalID).Delete(&open_models.OpenDigitalMember{}).Error; err != nil {
			return err
		}
		if err := tx.Where("digital_id = ?", partner.DigitalID).Delete(&open_models.OpenDigitalToken{}).Error; err != nil {
			return err
		}
		return tx.Delete(partner).Error
	})
	if err != nil {
		return nil, errors.New("删除数字伙伴失败")
	}

	// IM 用户是另一个域的实体，失败不回滚：主表已删，鉴权层查不到就是不可用
	if _, err := l.svcCtx.UserRpc.DeleteUsers(l.ctx, &user_rpc.DeleteUsersReq{
		UserIds: []string{partner.DigitalID},
	}); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "删除数字伙伴的 IM 用户失败",
			Data: map[string]interface{}{"digitalId": partner.DigitalID, "err": err.Error()},
		})
	}
	return &types.DeleteDigitalRes{}, nil
}
