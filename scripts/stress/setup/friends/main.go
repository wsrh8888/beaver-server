/*
 * 建状态：配对账号互加好友，产出 friends.json（含会话ID）
 *
 * 编排：worker pool 并发处理账号对，每对调本文件的 makeFriends flow 走真实加好友链路，
 * 成后调 friendInfo 拿 conversationId。flow 就在本文件里——Go flow 唯一消费者是 setup，不独立成层。
 *
 * 加好友链路（双向异步）：
 *   A add_friend(B) → B valid_list 拿 verifyID → B valid(verifyID, 同意) → 自动建会话
 * 会话随好友关系自动创建（friend_info 返回 conversationId），无独立建会话接口。
 *
 * 配对策略：环形配对，account[i] 加 account[(i+1)%N] 为好友，
 * 每个账号恰有 1 个好友和 1 个会话，简单够压测用。
 *
 * 运行（在 scripts/stress/setup 目录下）：
 *   go run ./friends -c 5
 *
 * 前置：先跑 users + tokens（账号要有 userId）
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	"beaver/scripts/stress/utils/config"
	"beaver/scripts/stress/utils/httputil"
	"beaver/scripts/stress/utils/types"
)

func main() {
	concurrency := flag.Int("c", 5, "并发 worker 数（加好友涉及双向请求，不宜过高）")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fatal("%v", err)
	}

	data, err := os.ReadFile(config.AccountsPath())
	if err != nil {
		fatal("读 accounts.json 失败（请先跑 users + tokens）: %v", err)
	}
	var accounts []types.Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		fatal("解析 accounts.json: %v", err)
	}
	if len(accounts) < 2 {
		fatal("账号数 < 2，无法配对，请先多造些账号")
	}

	type pair struct {
		a, b   types.Account
		aIndex int
	}
	// 环形配对：i 加 (i+1)%N
	pairs := make([]pair, 0, len(accounts))
	for i := 0; i < len(accounts); i++ {
		j := (i + 1) % len(accounts)
		pairs = append(pairs, pair{a: accounts[i], b: accounts[j], aIndex: i})
	}

	jobs := make(chan pair, *concurrency)
	type result struct {
		friendship types.Friendship
		ok         bool
		aIndex     int
	}
	results := make(chan result, *concurrency)

	var wg sync.WaitGroup
	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if err := makeFriends(cfg.Friend.APIBase, p.a.UserID, p.b.UserID); err != nil {
					fmt.Printf("  [%d] 加好友失败 %s↔%s: %v\n", p.aIndex, p.a.UserID, p.b.UserID, err)
					results <- result{ok: false, aIndex: p.aIndex}
					continue
				}
				// 拿会话ID（会话随好友关系自动创建）
				info, err := friendInfo(cfg.Friend.APIBase, p.a.UserID, p.b.UserID)
				if err != nil || info.ConversationID == "" {
					fmt.Printf("  [%d] 拿会话ID失败 %s↔%s: %v\n", p.aIndex, p.a.UserID, p.b.UserID, err)
					results <- result{ok: false, aIndex: p.aIndex}
					continue
				}
				results <- result{
					friendship: types.Friendship{
						UserID:         p.a.UserID,
						FriendID:       p.b.UserID,
						ConversationID: info.ConversationID,
					},
					ok:     true,
					aIndex: p.aIndex,
				}
			}
		}()
	}

	go func() {
		for _, p := range pairs {
			jobs <- p
		}
		close(jobs)
	}()
	go func() {
		wg.Wait()
		close(results)
	}()

	friendships := make([]types.Friendship, 0, len(pairs))
	var success, fail int32
	for r := range results {
		if r.ok {
			atomic.AddInt32(&success, 1)
			friendships = append(friendships, r.friendship)
		} else {
			atomic.AddInt32(&fail, 1)
		}
		if n := atomic.LoadInt32(&success) + atomic.LoadInt32(&fail); n%20 == 0 && n > 0 {
			fmt.Printf("  已处理 %d/%d\n", n, int32(len(pairs)))
		}
	}

	out, err := json.MarshalIndent(friendships, "", "  ")
	if err != nil {
		fatal("序列化 friends.json 失败: %v", err)
	}
	friendsPath := config.FriendsPath()
	if err := os.WriteFile(friendsPath, out, 0644); err != nil {
		fatal("写 friends.json 失败: %v", err)
	}
	fmt.Printf("\n✅ 完成: 成功 %d, 失败 %d\n", len(friendships), int(fail))
	fmt.Printf("✅ 好友关系已写入 %s\n", friendsPath)
}

// ===== 加好友真实业务 flow（对应 app/friend/friend_api/friend_api.api）=====
//
// 鉴权：只需 Beaver-User-Id header，不校验 token（见 friend_api/friend.go 中间件注册）。

// friendValidInfo 好友验证请求记录（valid_list 返回项）
type friendValidInfo struct {
	ID     string `json:"id"`     // 验证记录ID
	UserID string `json:"userId"` // 发送方用户ID
	Status int8   `json:"status"` // 0未处理 1已通过 2已拒绝
}

// friendInfoRes 好友详情（含会话ID）
type friendInfoRes struct {
	ConversationID string `json:"conversationId"`
	IsFriend       bool   `json:"isFriend"`
}

// addFriend A 向 B 发送好友验证请求。source: email / qrcode / userId
func addFriend(friendAPIBase, userID, friendID, source string) error {
	_, status, err := httputil.PostJSON(
		friendAPIBase+"/api/friend/v1/add_friend",
		map[string]string{"friendId": friendID, "source": source},
		map[string]string{"Beaver-User-Id": userID},
	)
	if err != nil {
		return fmt.Errorf("add_friend: %w", err)
	}
	if status != 200 {
		return fmt.Errorf("add_friend 失败: status=%d", status)
	}
	return nil
}

// validList 拉取 userID 的好友验证请求列表，用于找到发送方的 verifyID 以便同意。
func validList(friendAPIBase, userID string) ([]friendValidInfo, error) {
	body, status, err := httputil.PostJSON(
		friendAPIBase+"/api/friend/v1/valid_list",
		map[string]any{"page": 1, "limit": 50},
		map[string]string{"Beaver-User-Id": userID},
	)
	if err != nil {
		return nil, fmt.Errorf("valid_list: %w", err)
	}
	if status != 200 {
		return nil, fmt.Errorf("valid_list 失败: status=%d body=%s", status, body)
	}
	var res struct {
		Code   int               `json:"code"`
		Result struct {
			List []friendValidInfo `json:"list"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, fmt.Errorf("解析 valid_list: %w", err)
	}
	return res.Result.List, nil
}

// valid 处理好友验证请求。status: 1=同意 2=拒绝
func valid(friendAPIBase, userID, verifyID string, status int8) error {
	_, st, err := httputil.PostJSON(
		friendAPIBase+"/api/friend/v1/valid",
		map[string]any{"verifyId": verifyID, "status": status},
		map[string]string{"Beaver-User-Id": userID},
	)
	if err != nil {
		return fmt.Errorf("valid: %w", err)
	}
	if st != 200 {
		return fmt.Errorf("valid 失败: status=%d", st)
	}
	return nil
}

// friendInfo 拉好友详情，主要拿 conversationId（会话随好友关系自动创建）。
func friendInfo(friendAPIBase, userID, friendID string) (friendInfoRes, error) {
	body, status, err := httputil.GetJSON(
		friendAPIBase+"/api/friend/v1/friend_info?friendId="+friendID,
		map[string]string{"Beaver-User-Id": userID},
	)
	if err != nil {
		return friendInfoRes{}, fmt.Errorf("friend_info: %w", err)
	}
	if status != 200 {
		return friendInfoRes{}, fmt.Errorf("friend_info 失败: status=%d body=%s", status, body)
	}
	var res struct {
		Code   int            `json:"code"`
		Result friendInfoRes `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return friendInfoRes{}, fmt.Errorf("解析 friend_info: %w", err)
	}
	return res.Result, nil
}

// makeFriends 完整加好友 flow：A add_friend(B) → B valid_list 找到 A 的请求 → B valid(同意)。
func makeFriends(friendAPIBase, aUserID, bUserID string) error {
	if err := addFriend(friendAPIBase, aUserID, bUserID, "userId"); err != nil {
		return fmt.Errorf("A 发起加好友: %w", err)
	}
	list, err := validList(friendAPIBase, bUserID)
	if err != nil {
		return fmt.Errorf("B 拉验证列表: %w", err)
	}
	var verifyID string
	for _, v := range list {
		if v.UserID == aUserID && v.Status == 0 {
			verifyID = v.ID
			break
		}
	}
	if verifyID == "" {
		return fmt.Errorf("B 未找到来自 %s 的未处理验证请求", aUserID)
	}
	if err := valid(friendAPIBase, bUserID, verifyID, 1); err != nil {
		return fmt.Errorf("B 同意加好友: %w", err)
	}
	return nil
}

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", a...)
	os.Exit(1)
}
