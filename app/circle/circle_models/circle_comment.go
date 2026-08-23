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

package circle_models

import "beaver/common/models"

// CircleCommentModel 圈子帖子评论表
type CircleCommentModel struct {
	models.Model
	CommentID        string `gorm:"column:comment_id;size:64;uniqueIndex;not null" json:"commentId"` // 评论唯一ID
	PostID           string `gorm:"size:64;not null;index" json:"postId"`                            // 所属帖子ID
	CircleID         string `gorm:"size:64;not null;index" json:"circleId"`                          // 所属圈子ID（冗余便于查询）
	UserID           string `gorm:"size:64;not null;index" json:"userId"`                            // 评论用户ID
	Content          string `gorm:"type:text;not null" json:"content"`                               // 评论内容
	ParentID         string `gorm:"size:64;index" json:"parentId"`                                   // 父评论ID，空表示一级评论
	ReplyToCommentID string `gorm:"size:64" json:"replyToCommentId"`                                 // 被回复的评论ID
	ReplyToUserID    string `gorm:"size:64" json:"replyToUserId"`                                    // 被回复用户ID
	IsDeleted        bool   `gorm:"not null;default:false;index" json:"isDeleted"`                   // 软删除标记
}
