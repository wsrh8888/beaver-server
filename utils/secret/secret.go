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

// Package secret 提供可逆的对称加密，用于「接入通道」API Key 这类需要在运行时还原明文的机密。
//
// 与 utils/pwd（bcrypt 单向哈希，只用于口令校验、不可还原）相反：模型运行必须拿真实 Key 去调厂商，
// 所以这里用 AES-256-GCM 做可加可解密文。密文带 `v1:` 版本前缀，将来换算法/轮换密钥可凭前缀平滑迁移。
// 主密钥由部署配置注入（32 字节 base64），不进代码库、不进日志。
package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const (
	// encPrefix 密文版本前缀：解密时校验，未来 v2 可据前缀选择不同算法。
	encPrefix = "v1:"
	// keyBytes AES-256 需要 32 字节主密钥。
	keyBytes = 32
)

var (
	ErrKeyLength   = errors.New("secret: 主密钥必须为 32 字节（AES-256）")
	ErrCipherEmpty = errors.New("secret: 密文为空")
	ErrCipherForm  = errors.New("secret: 密文格式非法（缺少 v1: 前缀或 base64 解码失败）")
	ErrKeyMissing  = errors.New("secret: 未配置主密钥（SecretKey）")
)

// Cipher 绑定一把主密钥的 AES-256-GCM 加解密器。构造后并发安全（cipher.AEAD 自身安全）。
type Cipher struct {
	aead cipher.AEAD
}

// NewCipher 用 base64（标准编码）主密钥构造 Cipher。解码后必须恰好 32 字节，否则报 ErrKeyLength。
// 配置项 SecretKey 为空时返回 ErrKeyMissing，让服务启动即失败，而不是运行到调用点才崩。
func NewCipher(base64Key string) (*Cipher, error) {
	trimmed := strings.TrimSpace(base64Key)
	if trimmed == "" {
		return nil, ErrKeyMissing
	}
	raw, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, ErrKeyLength
	}
	if len(raw) != keyBytes {
		return nil, ErrKeyLength
	}
	block, err := aes.NewCipher(raw)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{aead: aead}, nil
}

// Encrypt 把明文加密成 `v1:base64(nonce||ciphertext||tag)`。空明文按空串返回（调用点自行判空）。
func (c *Cipher) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := c.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return encPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt 还原 Encrypt 产出的密文。非 v1 前缀 / 长度不足 / tag 校验失败都返回错误，绝不返回半成品。
func (c *Cipher) Decrypt(token string) (string, error) {
	if token == "" {
		return "", ErrCipherEmpty
	}
	if !strings.HasPrefix(token, encPrefix) {
		return "", ErrCipherForm
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(token, encPrefix))
	if err != nil {
		return "", ErrCipherForm
	}
	ns := c.aead.NonceSize()
	if len(raw) < ns {
		return "", ErrCipherForm
	}
	nonce, ciphertext := raw[:ns], raw[ns:]
	plain, err := c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// KeyHint 生成脱敏回显串：保留末 4 位、其余打码，供后台/客户端展示「这把 Key 长什么样」而不泄明文。
// 短于 8 位则整体打码，避免把短 Key 的绝大部分暴露出去。
func KeyHint(plaintext string) string {
	if plaintext == "" {
		return ""
	}
	r := []rune(plaintext)
	if len(r) < 8 {
		return strings.Repeat("*", len(r))
	}
	tail := string(r[len(r)-4:])
	return strings.Repeat("*", 4) + tail
}
