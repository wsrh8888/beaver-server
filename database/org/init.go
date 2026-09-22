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

package org

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	org_models "beaver/app/platform/platform_models/org"
	"beaver/app/user/user_models"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

const (
	defaultOrgName = "海狸"
	ownerUserID    = "100000" // 与 database/user/seed.go 的 ownerUserID 相同
)

// Org 当前安装的组织。
type Org struct {
	OrgID   string
	OrgName string
}

// ReadOrgName 读取平台配置里的组织名。文件不存在或未配置时用缺省名。
func ReadOrgName(path string) (string, error) {
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			log.Printf("未找到部署配置 %s，组织名使用缺省「%s」", path, defaultOrgName)
			return defaultOrgName, nil
		}
		return "", fmt.Errorf("读取部署配置失败: %w", readErr)
	}

	var file struct {
		Deploy struct {
			OrgName string `yaml:"OrgName"`
		} `yaml:"Deploy"`
	}
	if err := yaml.Unmarshal(data, &file); err != nil {
		return "", fmt.Errorf("解析部署配置失败: %w", err)
	}
	name := strings.TrimSpace(file.Deploy.OrgName)
	if name == "" {
		name = defaultOrgName
	}
	return name, nil
}

// Init 在用户种子之后执行。还没有组织时写入一个组织、一个根部门和拥有者。已有组织则跳过。
func Init(orgDB, userDB *gorm.DB, orgName string) error {
	if orgName == "" {
		orgName = defaultOrgName
	}

	var existing org_models.Organization
	err := orgDB.Take(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询组织失败: %w", err)
	}

	var owner user_models.UserModel
	if err = userDB.Where("user_id = ?", ownerUserID).Take(&owner).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("默认用户 %s 不存在，请先初始化用户", ownerUserID)
		}
		return fmt.Errorf("查询默认用户失败: %w", err)
	}

	return orgDB.Transaction(func(tx *gorm.DB) error {
		var n int64
		if err = tx.Model(&org_models.Organization{}).Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			return nil
		}

		orgID := strings.ReplaceAll(uuid.New().String(), "-", "")
		deptID := strings.ReplaceAll(uuid.New().String(), "-", "")
		if err = tx.Create(&org_models.Organization{
			OrgID: orgID,
			Name:  orgName,
		}).Error; err != nil {
			return fmt.Errorf("写入组织失败: %w", err)
		}
		if err = tx.Create(&org_models.Department{
			DeptID:   deptID,
			Name:     orgName,
			ParentID: "",
			Sort:     0,
			Path:     "/" + deptID,
		}).Error; err != nil {
			return fmt.Errorf("写入根部门失败: %w", err)
		}
		if err = tx.Create(&org_models.Staff{
			UserID: ownerUserID,
			Status: org_models.StaffActive,
		}).Error; err != nil {
			return fmt.Errorf("写入员工档案失败: %w", err)
		}
		if err = tx.Create(&org_models.DeptUser{
			UserID:   ownerUserID,
			DeptID:   deptID,
			IsMain:   true,
			IsLeader: true,
		}).Error; err != nil {
			return fmt.Errorf("写入部门关系失败: %w", err)
		}
		if err = tx.Create(&org_models.Admin{
			UserID: ownerUserID,
			Role:   org_models.AdminOwner,
		}).Error; err != nil {
			return fmt.Errorf("写入组织创建者失败: %w", err)
		}
		log.Printf("组织初始化完成: org=%s", orgName)
		return nil
	})
}

// Load 读取当前组织。还没有组织时返回 gorm.ErrRecordNotFound。
func Load(db *gorm.DB) (*Org, error) {
	var n int64
	if err := db.Model(&org_models.Organization{}).Count(&n).Error; err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if n != 1 {
		return nil, fmt.Errorf("组织数据异常")
	}

	var row org_models.Organization
	if err := db.Take(&row).Error; err != nil {
		return nil, err
	}
	return &Org{OrgID: row.OrgID, OrgName: row.Name}, nil
}
