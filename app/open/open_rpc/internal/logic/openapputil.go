/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package logic

import (
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"
)

func toOpenOAuthItem(app open_models.OpenOAuth) *open_rpc.OpenOAuthItem {
	return &open_rpc.OpenOAuthItem{
		OauthId:     app.OAuthID,
		Name:        app.Name,
		Description: app.Description,
		Icon:        app.Icon,
		OwnerUserId: app.OwnerID,
		Status:      int32(app.Status),
		CreatedAt:   time.Time(app.CreatedAt).UnixMilli(),
	}
}
