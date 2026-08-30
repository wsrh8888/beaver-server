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

package core

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter IP限流器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter 创建一个新的限流器
func NewRateLimiter(r float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(r),
		burst:    burst,
	}

	// 启动清理协程
	go rl.cleanupLoop()

	return rl
}

// Allow 检查是否允许请求通过
func (rl *RateLimiter) Allow(key string) bool {
	limiter := rl.getLimiter(key)
	return limiter.Allow()
}

// getLimiter 获取指定key的限流器，如果不存在则创建
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 双重检查
	if limiter, exists = rl.limiters[key]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[key] = limiter
	return limiter
}

// cleanupLoop 定期清理未使用的限流器
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup 清理长时间未使用的限流器
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 如果限流器数量超过阈值，进行清理
	if len(rl.limiters) > 10000 {
		// 创建新的map
		newLimiters := make(map[string]*rate.Limiter)

		// 只保留最近使用的限流器
		for key, limiter := range rl.limiters {
			// 使用 Tokens() 方法来近似判断限流器的最近使用时间
			// 如果令牌数量小于burst，说明最近有被使用过
			if limiter.Tokens() < float64(rl.burst) {
				newLimiters[key] = limiter
			}
		}

		rl.limiters = newLimiters
	}
}
