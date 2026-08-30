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

package conversation

// CityData 城市数据
type CityData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetDefaultCities 获取默认城市列表
func GetDefaultCities() []CityData {
	return []CityData{
		{Code: "ALL", Name: "全国"},
		{Code: "010", Name: "北京"},
		{Code: "021", Name: "上海"},
		{Code: "020", Name: "广州"},
		{Code: "0755", Name: "深圳"},
		{Code: "0571", Name: "杭州"},
		{Code: "028", Name: "成都"},
		{Code: "027", Name: "武汉"},
		{Code: "029", Name: "西安"},
		{Code: "025", Name: "南京"},
		{Code: "023", Name: "重庆"},
		{Code: "022", Name: "天津"},
		{Code: "0512", Name: "苏州"},
		{Code: "0731", Name: "长沙"},
		{Code: "0532", Name: "青岛"},
		{Code: "0510", Name: "无锡"},
		{Code: "0574", Name: "宁波"},
		{Code: "0371", Name: "郑州"},
		{Code: "0757", Name: "佛山"},
		{Code: "0769", Name: "东莞"},
	}
}
