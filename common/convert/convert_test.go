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

package convert

import (
	"testing"

	acctpb "beaver/app/account/account_rpc/pb"
)

// 模拟 go-zero 从 .api 生成的 HTTP DTO：lowerCamel tag + 缩写全大写（UserID）。
type httpVerifyPasswordReq struct {
	UserID   string `json:"userId,optional"`
	Password string `json:"password,optional"`
}

type httpUserLookupInfo struct {
	UserID   string `json:"userId"`
	NickName string `json:"nickName"`
	Version  int64  `json:"version"`
}

type httpUserSearchRes struct {
	UserInfo httpUserLookupInfo `json:"userInfo"`
}

// TestTo_ApiToPb_IgnoresWireTagStyle 验证：HTTP 层 lowerCamel + UserID 与
// pb 层 snake_case + UserId 之间，按 Go 字段名正确映射（不再依赖 json tag 一致）。
func TestTo_ApiToPb_IgnoresWireTagStyle(t *testing.T) {
	in := &httpVerifyPasswordReq{UserID: "100000", Password: "e10adc3949ba59abbe56e057f20f883e"}
	out, err := To[acctpb.VerifyPasswordReq](in)
	if err != nil {
		t.Fatalf("To failed: %v", err)
	}
	if out.UserId != "100000" {
		t.Fatalf("UserId = %q, want 100000", out.UserId)
	}
	if out.Password != in.Password {
		t.Fatalf("Password = %q, want %q", out.Password, in.Password)
	}
}

// TestTo_PbToApi_SnakeToCamel 验证 pb→HTTP 方向：pb 的 snake_case tag / UserId
// 能映射回 HTTP 的 lowerCamel / UserID，且 int 宽度转换（int32→int64）正常。
func TestTo_PbToApi_SnakeToCamel(t *testing.T) {
	in := &acctpb.SearchUserRes{UserInfo: &acctpb.UserInfo{
		UserId:   "100000",
		NickName: "测试用户385",
		Version:  7,
	}}
	out, err := To[httpUserSearchRes](in)
	if err != nil {
		t.Fatalf("To failed: %v", err)
	}
	if out.UserInfo.UserID != "100000" {
		t.Fatalf("UserID = %q, want 100000", out.UserInfo.UserID)
	}
	if out.UserInfo.NickName != "测试用户385" {
		t.Fatalf("NickName = %q", out.UserInfo.NickName)
	}
	if out.UserInfo.Version != 7 {
		t.Fatalf("Version = %d, want 7", out.UserInfo.Version)
	}
}
