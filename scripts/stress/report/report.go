/*
 * 把 k6 输出的 report.json 转成一张 HTML 报告表
 *
 * 运行（在 scripts/stress/report 目录下）：
 *   go run report.go <input.json> <output.html>
 *
 * 例：
 *   go run report.go ../single/auth/report.json ../single/auth/report.html
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type metricPoint struct {
	Metric string  `json:"metric"`
	Type   string  `json:"type"`
	Data   struct {
		Time   string  `json:"time"`
		Value  float64 `json:"value"`
		Tags   map[string]string `json:"tags"`
	} `json:"data"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("用法: go run report.go <input.json> <output.html>")
		fmt.Println("例:   go run report.go ../single/auth/report.json ../single/auth/report.html")
		os.Exit(1)
	}
	in, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("读 JSON 失败:", err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(string(in)), "\n")

	// 按指标收集所有数据点
	httpReqs := 0.0
	durations := []float64{}
	failedCount := 0.0
	totalCount := 0.0
	statusDist := map[string]int{}

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		var p metricPoint
		if err := json.Unmarshal([]byte(l), &p); err != nil {
			continue
		}
		if p.Type != "Point" {
			continue
		}
		switch p.Metric {
		case "http_reqs":
			httpReqs += p.Data.Value
			totalCount += 1
			if s, ok := p.Data.Tags["status"]; ok {
				statusDist[s]++
			}
		case "http_req_duration":
			durations = append(durations, p.Data.Value)
		case "http_req_failed":
			failedCount += p.Data.Value
		}
	}

	// 计算延迟分位数
	avgDur := 0.0
	p95 := 0.0
	p99 := 0.0
	maxDur := 0.0
	if len(durations) > 0 {
		sum := 0.0
		for _, d := range durations {
			sum += d
			if d > maxDur {
				maxDur = d
			}
		}
		avgDur = sum / float64(len(durations))

		// 排序算分位数
		sorted := make([]float64, len(durations))
		copy(sorted, durations)
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i] > sorted[j] {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		p95 = sorted[int(float64(len(sorted))*0.95)]
		p99 = sorted[int(float64(len(sorted))*0.99)]
	}

	failRate := 0.0
	if totalCount > 0 {
		failRate = (failedCount / totalCount) * 100
	}

	// 生成 HTML
	var sb strings.Builder
	sb.WriteString(`<!DOCTYPE html><html><head><meta charset="utf-8"><title>k6 压测报告</title>
<style>
body{font-family:sans-serif;margin:40px;background:#f5f5f5}
h1{color:#333}
.summary{display:flex;flex-wrap:wrap;gap:16px;margin:20px 0}
.card{background:#fff;padding:20px;border-radius:8px;box-shadow:0 2px 4px rgba(0,0,0,0.1);min-width:180px}
.card .label{color:#888;font-size:14px;margin-bottom:8px}
.card .value{font-size:28px;font-weight:bold;color:#333}
.card.warn .value{color:#e74c3c}
.card.ok .value{color:#27ae60}
table{background:#fff;border-collapse:collapse;width:100%;max-width:600px;margin-top:20px;box-shadow:0 2px 4px rgba(0,0,0,0.1)}
td,th{border:1px solid #eee;padding:10px 16px;text-align:left}
th{background:#fafafa}
</style></head><body>`)

	sb.WriteString(`<h1>k6 压测报告</h1>`)
	sb.WriteString(`<p style="color:#888">接口：auth_api 邮箱密码登录</p>`)

	// 汇总卡片
	sb.WriteString(`<div class="summary">`)
	card(&sb, "总请求数", fmt.Sprintf("%.0f", httpReqs), false)
	card(&sb, "平均延迟", fmt.Sprintf("%.0f ms", avgDur), avgDur > 1000)
	card(&sb, "P95 延迟", fmt.Sprintf("%.0f ms", p95), p95 > 1000)
	card(&sb, "P99 延迟", fmt.Sprintf("%.0f ms", p99), p99 > 1000)
	card(&sb, "最大延迟", fmt.Sprintf("%.0f ms", maxDur), false)
	card(&sb, "错误率", fmt.Sprintf("%.1f%%", failRate), failRate > 1)
	sb.WriteString(`</div>`)

	// 状态码分布
	sb.WriteString(`<h2>HTTP 状态码分布</h2><table><tr><th>状态码</th><th>次数</th><th>占比</th></tr>`)
	for status, cnt := range statusDist {
		pct := float64(cnt) / totalCount * 100
		cls := ""
		if status != "200" {
			cls = ` style="color:#e74c3c"`
		}
		sb.WriteString(fmt.Sprintf("<tr><td%s>%s</td><td>%d</td><td>%.1f%%</td></tr>", cls, status, cnt, pct))
	}
	sb.WriteString(`</table>`)

	sb.WriteString(`<p style="margin-top:20px;color:#888;font-size:13px">数据来源：`)
	sb.WriteString(os.Args[1])
	sb.WriteString(`</p>`)

	sb.WriteString(`</body></html>`)

	if err := os.WriteFile(os.Args[2], []byte(sb.String()), 0644); err != nil {
		fmt.Println("写 HTML 失败:", err)
		os.Exit(1)
	}
	fmt.Println("✅ 报告已生成:", os.Args[2])
}

func card(sb *strings.Builder, label, value string, warn bool) {
	cls := "card"
	if warn {
		cls += " warn"
	}
	sb.WriteString(fmt.Sprintf(`<div class="%s"><div class="label">%s</div><div class="value">%s</div></div>`, cls, label, value))
}
