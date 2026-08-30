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

package list_query

import (
	"beaver/common/models"
	"fmt"

	"gorm.io/gorm"
)

type Option struct {
	PageInfo models.PageInfo
	Where    *gorm.DB //高级查询
	Likes    []string //模糊查询
	Joins    string
	Debug    bool
	Preload  []string             //预加载
	Table    func() (string, any) //子查询
	Groups   []string             //分组查询
}

func ListQuery[T any](db *gorm.DB, model T, option Option) (list []T, count int64, err error) {
	if option.Debug {
		db = db.Debug()
	}
	query := db.Where(model) //把结构体自己的查询条件查了

	// 模糊查询
	if option.PageInfo.Key != "" && len(option.Likes) > 0 {
		likeQuery := db.Where("")
		for index, column := range option.Likes {
			if index == 0 {
				likeQuery = likeQuery.Where(fmt.Sprintf("%s like '%%?%%'", column), option.PageInfo.Key)
			} else {
				likeQuery = likeQuery.Or(fmt.Sprintf("%s like '%%?%%'", column), option.PageInfo.Key)
			}
		}
		query.Where(likeQuery)
	}

	if option.Table != nil {
		table, data := option.Table()
		query = query.Table(table, data)
	}

	if len(option.Groups) > 0 {
		for _, group := range option.Groups {
			query = query.Group(group)
		}
	}

	if option.Joins != "" {
		query = query.Joins(option.Joins)
	}

	if option.Where != nil {
		query = query.Where(option.Where)
	}

	// 求总数
	query.Model(model).Count(&count)

	// 预加载
	for _, s := range option.Preload {
		query = query.Preload(s)
	}

	//分页查询
	if option.PageInfo.Page <= 0 {
		option.PageInfo.Page = 1
	}
	if option.PageInfo.Limit <= 0 && option.PageInfo.Limit != -1 {
		option.PageInfo.Limit = 10
	}
	fmt.Println(option.PageInfo, option.PageInfo.Limit)
	offset := (option.PageInfo.Page - 1) * option.PageInfo.Limit

	if option.PageInfo.Sort != "" {
		query.Order(option.PageInfo.Sort)
	}
	err = query.Limit(option.PageInfo.Limit).Offset(offset).Find(&list).Error

	return
}
