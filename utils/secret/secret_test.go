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

package secret

import (
	"encoding/base64"
	"strings"
	"testing"
)

// testKey 一把合法的 32 字节 base64 主密钥，仅测试用。
func testKey() string {
	return base64.StdEncoding.EncodeToString(make([]byte, 32))
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	c, err := NewCipher(testKey())
	if err != nil {
		t.Fatalf("NewCipher: %v", err)
	}
	plain := "sk-abcdef1234567890"
	token, err := c.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !strings.HasPrefix(token, "v1:") {
		t.Fatalf("密文应带 v1: 前缀, got %q", token)
	}
	if token == plain {
		t.Fatal("密文不应等于明文")
	}
	got, err := c.Decrypt(token)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plain {
		t.Fatalf("回环失败: want %q got %q", plain, got)
	}
}

func TestEncryptEmptyAndRandomNonce(t *testing.T) {
	c, _ := NewCipher(testKey())
	if s, err := c.Encrypt(""); err != nil || s != "" {
		t.Fatalf("空明文应返回空串, got %q err %v", s, err)
	}
	// 同一明文两次加密应因随机 nonce 而不同，但都能解回。
	plain := "sk-same-input"
	a, _ := c.Encrypt(plain)
	b, _ := c.Encrypt(plain)
	if a == b {
		t.Fatal("两次密文不应相同（nonce 应随机）")
	}
	da, _ := c.Decrypt(a)
	db, _ := c.Decrypt(b)
	if da != plain || db != plain {
		t.Fatal("两次密文都应能还原")
	}
}

func TestDecryptRejectsTamper(t *testing.T) {
	c, _ := NewCipher(testKey())
	token, _ := c.Encrypt("sk-secret-value")
	// 篡改末位字符（破坏 GCM tag）应报错。
	bad := token[:len(token)-2] + "AA"
	if _, err := c.Decrypt(bad); err == nil {
		t.Fatal("篡改密文应解密失败")
	}
	// 缺前缀应报 ErrCipherForm。
	if _, err := c.Decrypt("v2:whatever"); err != ErrCipherForm {
		t.Fatalf("非 v1 前缀应报 ErrCipherForm, got %v", err)
	}
	if _, err := c.Decrypt(""); err != ErrCipherEmpty {
		t.Fatalf("空密文应报 ErrCipherEmpty, got %v", err)
	}
}

func TestNewCipherKeyValidation(t *testing.T) {
	if _, err := NewCipher(""); err != ErrKeyMissing {
		t.Fatalf("空密钥应报 ErrKeyMissing, got %v", err)
	}
	shortKey := base64.StdEncoding.EncodeToString(make([]byte, 16))
	if _, err := NewCipher(shortKey); err != ErrKeyLength {
		t.Fatalf("16 字节应报 ErrKeyLength, got %v", err)
	}
	if _, err := NewCipher("not-valid-base64!!!"); err != ErrKeyLength {
		t.Fatalf("非法 base64 应报 ErrKeyLength, got %v", err)
	}
}

func TestKeyHint(t *testing.T) {
	if KeyHint("") != "" {
		t.Fatal("空串 hint 应为空")
	}
	if h := KeyHint("sk-proj-abcdefgh1234"); !strings.HasSuffix(h, "1234") || !strings.HasPrefix(h, "****") {
		t.Fatalf("长 Key hint 应 ****+末4位, got %q", h)
	}
	if h := KeyHint("abc"); strings.ContainsAny(h, "abc") {
		t.Fatalf("短 Key hint 应全打码, got %q", h)
	}
}
