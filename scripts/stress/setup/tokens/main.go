/*
 * 建状态：给 accounts.json 里的账号批量登录拿 token
 *
 * 编排：worker pool，每个 worker 调本文件的 loginOne flow 走真实登录链路。
 * flow 就在本文件里——Go flow 唯一消费者是 setup，不独立成层。
 *
 * 登录链路：调 /email_password_login，返回 token + userId。
 *
 * 运行（在 scripts/stress/setup 目录下）：
 *   go run ./tokens -c 10
 *
 * 前置：先跑 users 生成 accounts.json
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"

	"beaver/scripts/stress/utils/config"
	"beaver/scripts/stress/utils/httputil"
	"beaver/scripts/stress/utils/types"
)

func main() {
	concurrency := flag.Int("c", 10, "并发登录 worker 数")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fatal("%v", err)
	}

	data, err := os.ReadFile(config.AccountsPath())
	if err != nil {
		fatal("读 accounts.json 失败（请先跑 users）: %v", err)
	}
	var accounts []types.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		fatal("解析 accounts.json: %v", err)
	}

	type job struct{ i int }
	jobs := make(chan job, *concurrency)

	var wg sync.WaitGroup
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				i := j.i
				token, userID, err := loginOne(cfg.Auth.APIBase, accounts[i].Email, accounts[i].Password, fmt.Sprintf("stress-device-%03d", i))
				if err != nil {
					fmt.Printf("  [%d] 登录失败 %s: %v\n", i, accounts[i].Email, err)
					continue
				}
				accounts[i].Token = token
				accounts[i].UserID = userID
			}
		}()
	}

	go func() {
		for i := range accounts {
			jobs <- job{i}
		}
		close(jobs)
	}()
	wg.Wait()

	out, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		fatal("序列化 accounts.json 失败: %v", err)
	}
	if err := os.WriteFile(config.AccountsPath(), out, 0644); err != nil {
		fatal("写回 accounts.json 失败: %v", err)
	}
	fmt.Printf("✅ 完成: 共处理 %d 个账号，token 已回填到 %s\n", len(accounts), config.AccountsPath())
}

// loginOne 邮箱密码登录（真实业务 flow），返回 token + userId。
// 登录接口需带 BeaverDesktop UA 做设备类型校验。
func loginOne(authAPIBase, email, password, deviceID string) (token, userID string, err error) {
	body, status, err := httputil.PostJSON(
		authAPIBase+"/api/auth/auth_public/v1/email_password_login",
		map[string]string{"email": email, "password": password},
		map[string]string{
			"deviceId":   deviceID,
			"User-Agent": "BeaverDesktop/1.0 (Windows)",
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("登录: %w", err)
	}
	if status != 200 {
		return "", "", fmt.Errorf("登录失败: status=%d body=%s", status, body)
	}
	var res struct {
		Code   int    `json:"code"`
		Result struct {
			Token  string `json:"token"`
			UserID string `json:"userId"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return "", "", fmt.Errorf("解析登录响应: %w", err)
	}
	if res.Result.Token == "" {
		return "", "", fmt.Errorf("登录响应无 token: %s", body)
	}
	return res.Result.Token, res.Result.UserID, nil
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", a...)
	os.Exit(1)
}
