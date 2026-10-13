/*
 * 建状态：批量注册测试用户
 *
 * 编排：worker pool + ticker 限速，每个 worker 调本文件的 registerOne flow 走真实注册链路。
 * flow 就在本文件里——Go flow 唯一消费者是 setup，不独立成层。
 *
 * 注册链路：直连 auth Redis 写验证码（绕过邮件发送，但注册接口校验链路照走）
 *           → 调 /email_register（服务端从 Redis 读验证码校验 → 查重 → 建用户 → 写凭证）
 *
 * 运行（在 scripts/stress/setup 目录下）：
 *   go run ./users -count 500 -c 10 -rps 5
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"beaver/scripts/stress/utils/config"
	"beaver/scripts/stress/utils/httputil"
	"beaver/scripts/stress/utils/types"

	"github.com/go-redis/redis"
)

func main() {
	count := flag.Int("count", 500, "要注册的测试账号数量")
	start := flag.Int("start", 0, "邮箱序号起始值（避免重复注册冲突）")
	concurrency := flag.Int("c", 10, "并发 worker 数")
	rps := flag.Int("rps", 5, "每秒提交注册数上限（保护服务端不被打爆）")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fatal("%v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Auth.RedisAddr,
		Password: cfg.Auth.RedisPassword,
		DB:       cfg.Auth.RedisDB,
	})
	if err := rdb.Ping().Err(); err != nil {
		fatal("连接 auth Redis 失败: %v", err)
	}
	fmt.Printf("✅ auth Redis 连接成功，目标 auth_api: %s（并发=%d 限速=%drps）\n", cfg.Auth.APIBase, *concurrency, *rps)

	type job struct {
		idx   int
		email string
	}
	jobs := make(chan job, *concurrency)
	type result struct {
		acct types.Account
		ok   bool
	}
	results := make(chan result, *concurrency)

	var wg sync.WaitGroup
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				acct, err := registerOne(rdb, cfg.Auth.APIBase, j.email, fmt.Sprintf("stress-device-%03d", j.idx))
				if err != nil {
					fmt.Printf("  [%d] 注册失败 %s: %v\n", j.idx, j.email, err)
					results <- result{ok: false}
					continue
				}
				results <- result{acct: acct, ok: true}
			}
		}()
	}

	go func() {
		tick := time.NewTicker(time.Second / time.Duration(*rps))
		defer tick.Stop()
		for i := *start; i < *start+*count; i++ {
			<-tick.C
			jobs <- job{idx: i, email: fmt.Sprintf("stress_test_%d@beaver.test", i)}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	accounts := make([]types.Account, 0, *count)
	var success, fail int32
	for r := range results {
		if r.ok {
			atomic.AddInt32(&success, 1)
			accounts = append(accounts, r.acct)
		} else {
			atomic.AddInt32(&fail, 1)
		}
		if n := atomic.LoadInt32(&success) + atomic.LoadInt32(&fail); n%50 == 0 && n > 0 {
			fmt.Printf("  已处理 %d/%d\n", n, int32(*count))
		}
	}

	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].Email < accounts[j].Email
	})

	out, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		fatal("序列化 accounts.json 失败: %v", err)
	}
	if err := os.WriteFile(config.AccountsPath(), out, 0644); err != nil {
		fatal("写 accounts.json 失败: %v", err)
	}
	fmt.Printf("\n✅ 完成: 成功 %d, 失败 %d\n", len(accounts), int(fail))
	fmt.Printf("✅ 账号已写入 %s\n", config.AccountsPath())
}

// registerOne 注册一个账号（真实业务 flow）。
// 走完整注册链路：写验证码 → 调 /email_register（服务端从 Redis 读验证码校验 → 查重 → 建用户 → 写凭证）。
// 绕过邮件发送但保留校验链路，既不依赖 SMTP 又覆盖查重/建用户/写凭证。
func registerOne(rdb *redis.Client, authAPIBase, email, deviceID string) (types.Account, error) {
	code := "123456"
	codeKey := fmt.Sprintf("email_code_%s_%s", email, "register")
	if err := rdb.Set(codeKey, code, 5*time.Minute).Err(); err != nil {
		return types.Account{}, fmt.Errorf("写验证码: %w", err)
	}

	body, status, err := httputil.PostJSON(
		authAPIBase+"/api/auth/auth_public/v1/email_register",
		map[string]string{"email": email, "password": types.TestPassword, "code": code},
		map[string]string{"deviceId": deviceID, "User-Agent": "BeaverDesktop/1.0 (Windows)"},
	)
	if err != nil {
		return types.Account{}, fmt.Errorf("注册: %w", err)
	}
	if status != 200 || strings.Contains(body, "\"code\":1") || strings.Contains(body, "失败") || strings.Contains(body, "已注册") {
		return types.Account{}, fmt.Errorf("注册失败: status=%d body=%s", status, body)
	}
	return types.Account{Email: email, Password: types.TestPassword}, nil
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", a...)
	os.Exit(1)
}
