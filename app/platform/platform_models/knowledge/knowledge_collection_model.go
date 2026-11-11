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

package knowledge_models

import "beaver/common/models"

// KnowledgeCollection 知识库集合（目录树节点）。表 knowledge_collections。
// 单租户：无 org_id / product_scope；owner_user_id 可空（系统内置库为空）。
// parent_id 支撑两层目录树（根为空串）；index_config 承载分段/清洗/检索与 embedding 模型配置。
type KnowledgeCollection struct {
	models.Model
	CollectionID  string `gorm:"column:collection_id;type:varchar(64);uniqueIndex;not null;comment:集合业务ID" json:"collectionId"`
	ParentID      string `gorm:"column:parent_id;type:varchar(64);not null;default:'';index;comment:父集合ID，空为根；仅支持两层" json:"parentId"`
	Name          string `gorm:"type:varchar(128);not null;comment:集合名" json:"name"`
	Kind          string `gorm:"type:varchar(32);not null;default:'docs';comment:类型 docs文档库 faq问答对 ..." json:"kind"`
	StorageRef    string `gorm:"column:storage_ref;type:varchar(255);not null;default:'';comment:对象存储前缀/桶引用" json:"storageRef"`
	AclMode       string `gorm:"column:acl_mode;type:varchar(32);not null;default:'internal';comment:user_scoped按用户隔离 internal内部共享" json:"aclMode"`
	IndexConfig   string `gorm:"column:index_config;type:text;not null;comment:分段/清洗/检索配置 JSON {embedding_model_id,chunk_size,overlap,top_k,...}；整块读整块写，不做 JSON 路径查询，故用 text" json:"indexConfig"`
	Description   string `gorm:"type:varchar(512);not null;default:'';comment:描述" json:"description"`
	OwnerUserID   string `gorm:"column:owner_user_id;type:varchar(64);not null;default:'';index;comment:归属用户，系统内置为空" json:"ownerUserId"`
	Status        int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
	Sort          int    `gorm:"type:int;not null;default:0;comment:排序" json:"sort"`
}

func (KnowledgeCollection) TableName() string { return "knowledge_collections" }
