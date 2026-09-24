/*
 * 造数据：批量注册测试用户
 *
 * 作用：调 auth_api 的真实注册接口批量创建测试账号，
 *       走完整注册链路（验证码校验→查重→建用户→写凭证）。
 *       注册完导出 accounts.json 给所有压测脚本共用。
 *
 * 流程：对每个账号
 *   1) 调 /emailcode 接口让服务端往 Redis 写入验证码
 *   2) 直连 Redis 读出验证码
 *   3) 调 /email_register 接口完成注册
 *
 * 运行（在 scripts/stress/seed 目录下）：
 *   go run users.go -count 500
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"beaver/scripts/stress/lib"

	"github.com/go-redis/redis"
)

func main() {
	count := flag.Int("count", 500, "要注册的测试账号数量")
	start := flag.Int("start", 0, "邮箱序号起始值（避免重复注册冲突）")
	flag.Parse()

	cfg, err := lib.Load()
	if err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}

	// 用 auth 服务端自己的 Redis（验证码存在 auth 的 Redis 里）
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Auth.RedisAddr,
		Password: cfg.Auth.RedisPassword,
		DB:       cfg.Auth.RedisDB,
	})
	if err := rdb.Ping().Err(); err != nil {
		fmt.Printf("❌ 连接 auth Redis 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ auth Redis 连接成功，目标 auth_api:", cfg.Auth.APIBase)

	accounts := make([]lib.Account, 0, *count)
	fail := 0

	for i := *start; i < *start+*count; i++ {
		email := fmt.Sprintf("stress_test_%d@beaver.test", i)
		acct, err := registerOne(rdb, cfg.Auth.APIBase, email)
		if err != nil {
			fmt.Printf("  [%d] 注册失败 %s: %v\n", i, email, err)
			fail++
			continue
		}
		accounts = append(accounts, acct)
		fmt.Printf("  [%d] 注册成功 %s\n", i, email)
		time.Sleep(2 * time.Second) // 限速，避免触发服务端频率限制
	}

	out, _ := json.MarshalIndent(accounts, "", "  ")
	if err := os.WriteFile("../accounts.json", out, 0644); err != nil {
		fmt.Printf("❌ 写 accounts.json 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()
	fmt.Printf("✅ 完成: 成功 %d, 失败 %d\n", len(accounts), fail)
	fmt.Printf("✅ 账号已写入 scripts/stress/accounts.json\n")
}

func registerOne(rdb *redis.Client, baseURL, email string) (lib.Account, error) {
	// 1) 直接往 Redis 写验证码（绕过邮件发送，但注册接口的校验链路照走）
	//    这样既不依赖 SMTP 配置，又完整走了 注册接口的验证码校验→查重→建用户→写凭证
	code := "123456"
	codeKey := fmt.Sprintf("email_code_%s_%s", email, "register")
	if err := rdb.Set(codeKey, code, 5*time.Minute).Err(); err != nil {
		return lib.Account{}, fmt.Errorf("写入验证码: %w", err)
	}

	// 2) 调注册接口（服务端会从 Redis 读验证码校验）
	body, status, err := lib.PostJSON(
		baseURL+"/api/auth/auth_public/v1/email_register",
		map[string]string{"email": email, "password": lib.TestPassword, "code": code},
		nil,
	)
	if err != nil {
		return lib.Account{}, fmt.Errorf("注册: %w", err)
	}
	if status != 200 || strings.Contains(body, "\"code\":1") || strings.Contains(body, "失败") || strings.Contains(body, "已注册") {
		return lib.Account{}, fmt.Errorf("注册失败: status=%d body=%s", status, body)
	}

	return lib.Account{Email: email, Password: lib.TestPassword}, nil
}
