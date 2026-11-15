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

// Package convert 在「同形不同型」的两套结构之间做一次 JSON 往返转换。
//
// 用在哪：业务 API 的 HTTP 出入参（.api 生成的 types.XxxReq/XxxRes）与同域 RPC 的
// proto 出入参，字段名与 json tag 是刻意对齐的（同一份业务语义，两套 Go 类型）。
// 逐字段手写赋值既啰嗦又容易漏，这里统一走 JSON：序列化 → 反序列化。
//
// 不用在：形状不同的两套结构（那样转换必然丢字段，应该在设计上对齐，而不是在这里打补丁）。
//
// 代价是一次内存内的 marshal/unmarshal，相对一次跨进程 RPC 可忽略。
package convert

import "encoding/json"

// To 把 in 按 JSON 语义转换成 *T。
//
// 转换失败（如字段类型不兼容）返回错误，调用方应把它当业务错误抛给上层，
// 不要静默返回零值——那会让「契约漂移」变成线上难查的空响应。
func To[T any](in any) (*T, error) {
	if in == nil {
		return new(T), nil
	}
	raw, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	out := new(T)
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, err
	}
	return out, nil
}
