package logic

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ---- ProseMirror JSON -> 导出格式 ----

type pmNode struct {
	Type    string         `json:"type"`
	Attrs   map[string]any `json:"attrs"`
	Content []pmNode       `json:"content"`
	Text    string         `json:"text"`
	Marks   []pmMark       `json:"marks"`
}

type pmMark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs"`
}

func parsePM(contentJSON string) pmNode {
	var root pmNode
	if trim(contentJSON) == "" {
		return root
	}
	if err := json.Unmarshal([]byte(contentJSON), &root); err != nil {
		return pmNode{}
	}
	return root
}

func pmAttrString(attrs map[string]any, key string) string {
	if attrs == nil {
		return ""
	}
	if v, ok := attrs[key].(string); ok {
		return v
	}
	return ""
}

func pmAttrInt(attrs map[string]any, key string, def int) int {
	if attrs == nil {
		return def
	}
	switch v := attrs[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	}
	return def
}

func hasMark(n pmNode, typ string) bool {
	for _, m := range n.Marks {
		if m.Type == typ {
			return true
		}
	}
	return false
}

func markAttrs(n pmNode, typ string) map[string]any {
	for _, m := range n.Marks {
		if m.Type == typ {
			return m.Attrs
		}
	}
	return nil
}

func pmInline(nodes []pmNode, asHTML bool) string {
	var sb strings.Builder
	for _, n := range nodes {
		switch n.Type {
		case "text":
			t := n.Text
			if asHTML {
				t = htmlEscape(t)
			}
			if hasMark(n, "code") {
				t = wrap(t, "`", "<code>", asHTML)
			}
			if hasMark(n, "bold") || hasMark(n, "strong") {
				t = wrap(t, "**", "<strong>", asHTML)
			}
			if hasMark(n, "italic") || hasMark(n, "em") {
				t = wrap(t, "*", "<em>", asHTML)
			}
			if hasMark(n, "strike") {
				t = wrap(t, "~~", "<s>", asHTML)
			}
			if a := markAttrs(n, "link"); a != nil {
				href := pmAttrString(a, "href")
				if href != "" {
					if asHTML {
						t = `<a href="` + htmlEscape(href) + `">` + t + `</a>`
					} else {
						t = "[" + t + "](" + href + ")"
					}
				}
			}
			sb.WriteString(t)
		case "hardBreak":
			if asHTML {
				sb.WriteString("<br/>")
			} else {
				sb.WriteString("\n")
			}
		case "image":
			src := pmAttrString(n.Attrs, "src")
			alt := pmAttrString(n.Attrs, "alt")
			if asHTML {
				sb.WriteString(`<img src="` + htmlEscape(src) + `" alt="` + htmlEscape(alt) + `"/>`)
			} else {
				sb.WriteString("![" + alt + "](" + src + ")")
			}
		default:
			sb.WriteString(pmInline(n.Content, asHTML))
		}
	}
	return sb.String()
}

func wrap(s, md, open string, asHTML bool) string {
	if asHTML {
		closeTag := "</" + strings.TrimPrefix(strings.TrimSuffix(open, ">"), "<") + ">"
		return open + s + closeTag
	}
	return md + s + md
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;")
	return r.Replace(s)
}

// pmToMarkdown 把 ProseMirror 文档树转成 Markdown。
func pmToMarkdown(contentJSON string) string {
	root := parsePM(contentJSON)
	var sb strings.Builder
	pmBlocks(root.Content, &sb, false, 0)
	return strings.TrimSpace(sb.String()) + "\n"
}

// pmToHTML 把 ProseMirror 文档树转成 HTML 片段。
func pmToHTML(contentJSON string) string {
	root := parsePM(contentJSON)
	var sb strings.Builder
	pmBlocks(root.Content, &sb, true, 0)
	return strings.TrimSpace(sb.String())
}

func pmBlocks(nodes []pmNode, sb *strings.Builder, asHTML bool, depth int) {
	for _, n := range nodes {
		inner := pmInline(n.Content, asHTML)
		switch n.Type {
		case "paragraph":
			if asHTML {
				sb.WriteString("<p>" + inner + "</p>\n")
			} else {
				sb.WriteString(inner + "\n\n")
			}
		case "heading":
			level := pmAttrInt(n.Attrs, "level", 1)
			if level < 1 {
				level = 1
			}
			if level > 6 {
				level = 6
			}
			if asHTML {
				fmt.Fprintf(sb, "<h%d>%s</h%d>\n", level, inner, level)
			} else {
				sb.WriteString(strings.Repeat("#", level) + " " + inner + "\n\n")
			}
		case "bulletList", "orderedList":
			for i, item := range n.Content {
				prefix := "- "
				if n.Type == "orderedList" {
					prefix = fmt.Sprintf("%d. ", i+1)
				}
				if asHTML {
					sb.WriteString("<li>" + pmInline(item.Content, true) + "</li>\n")
				} else {
					sb.WriteString(prefix + pmBlocksText(item.Content) + "\n")
				}
			}
			if !asHTML {
				sb.WriteString("\n")
			}
		case "blockquote":
			if asHTML {
				sb.WriteString("<blockquote>" + inner + "</blockquote>\n")
			} else {
				for _, line := range strings.Split(inner, "\n") {
					sb.WriteString("> " + line + "\n")
				}
				sb.WriteString("\n")
			}
		case "codeBlock":
			code := inner
			if asHTML {
				sb.WriteString("<pre><code>" + htmlEscape(code) + "</code></pre>\n")
			} else {
				lang := pmAttrString(n.Attrs, "language")
				sb.WriteString("```" + lang + "\n" + code + "\n```\n\n")
			}
		case "horizontalRule":
			if asHTML {
				sb.WriteString("<hr/>\n")
			} else {
				sb.WriteString("---\n\n")
			}
		default:
			if len(n.Content) > 0 {
				pmBlocks(n.Content, sb, asHTML, depth+1)
			} else if inner != "" {
				sb.WriteString(inner + "\n\n")
			}
		}
	}
}

func pmBlocksText(nodes []pmNode) string {
	var sb strings.Builder
	for _, n := range nodes {
		if n.Type == "paragraph" {
			sb.WriteString(pmInline(n.Content, false))
		} else {
			sb.WriteString(pmInline(n.Content, false))
		}
	}
	return sb.String()
}

// ---- 最小 OOXML (.docx) 生成：不引入第三方依赖，纯标准库 zip + XML ----

func buildDocx(title, text string) ([]byte, error) {
	var body strings.Builder
	body.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	body.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`)
	lines := append([]string{title}, strings.Split(text, "\n")...)
	for _, line := range lines {
		body.WriteString(`<w:p><w:r><w:t xml:space="preserve">` + htmlEscape(line) + `</w:t></w:r></w:p>`)
	}
	body.WriteString(`<w:sectPr/></w:body></w:document>`)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	parts := []struct{ name, data string }{
		{"[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`},
		{"_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`},
		{"word/document.xml", body.String()},
	}
	for _, p := range parts {
		w, err := zw.Create(p.name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(p.data)); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ---- 行级 diff：按行集合差给出新增/删除行，用于版本对比概览 ----

func lineDiff(base, target string) (added, removed []string) {
	baseSet := map[string]int{}
	for _, line := range strings.Split(base, "\n") {
		if l := strings.TrimSpace(line); l != "" {
			baseSet[l]++
		}
	}
	targetSet := map[string]int{}
	for _, line := range strings.Split(target, "\n") {
		if l := strings.TrimSpace(line); l != "" {
			targetSet[l]++
		}
	}
	for line, n := range targetSet {
		for i := n - baseSet[line]; i > 0; i-- {
			added = append(added, line)
		}
	}
	for line, n := range baseSet {
		for i := n - targetSet[line]; i > 0; i-- {
			removed = append(removed, line)
		}
	}
	return added, removed
}
