/*
 * 造数据：给 accounts.json 里的账号批量登录拿 token
 *
 * 作用：很多接口（发消息、好友列表等）需要鉴权 token，
 *       这个脚本批量登录已有账号，把 token 回填到 accounts.json。
 *
 * 运行（在 scripts/stress/seed 目录下）：
 *   go run tokens.go
 *
 * 前置：先跑 users.go 生成 accounts.json
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"beaver/scripts/stress/lib"
)

func main() {
	cfg, err := lib.Load()
	if err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}

	data, err := os.ReadFile("../accounts.json")
	if err != nil {
		fmt.Println("❌ 读 accounts.json 失败（请先跑 users.go）:", err)
		os.Exit(1)
	}
	var accounts []lib.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		fmt.Println("❌ 解析 accounts.json 失败:", err)
		os.Exit(1)
	}

	success := 0
	for i := range accounts {
		body, _, err := lib.PostJSON(
			cfg.Auth.APIBase+"/api/auth/auth_public/v1/email_password_login",
			map[string]string{"email": accounts[i].Email, "password": accounts[i].Password},
			map[string]string{
				"deviceId":   fmt.Sprintf("stress-device-%03d", i),
				"User-Agent": "BeaverDesktop/1.0 (Windows)",
			},
		)
		if err != nil {
			fmt.Printf("  [%d] 登录失败 %s: %v\n", i, accounts[i].Email, err)
			continue
		}
		var res struct {
			Code   int    `json:"code"`
			Result struct {
				Token  string `json:"token"`
				UserID string `json:"userId"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(body), &res); err != nil || res.Result.Token == "" {
			fmt.Printf("  [%d] 解析登录响应失败\n", i)
			continue
		}
		accounts[i].Token = res.Result.Token
		accounts[i].UserID = res.Result.UserID
		success++
		if (i+1)%50 == 0 {
			fmt.Printf("  已处理 %d/%d\n", i+1, len(accounts))
		}
	}

	out, _ := json.MarshalIndent(accounts, "", "  ")
	if err := os.WriteFile("../accounts.json", out, 0644); err != nil {
		fmt.Println("❌ 写回 accounts.json 失败:", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 完成: 成功登录 %d/%d，token 已回填到 accounts.json\n", success, len(accounts))
}
