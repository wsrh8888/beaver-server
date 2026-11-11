/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package agent

import (
	"encoding/json"
	"errors"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// skillDoc 一份 SKILL.md 解析后的结果。
type skillDoc struct {
	Name        string
	Description string
	Version     string
	Body        string
	ConfigJSON  string
}

// parseSkillFile 解析 SKILL.md：frontmatter 取 name/description/version，其余键收进配置 JSON，
// 正文原样保留。没有 frontmatter 时用文件名当名称、正文首行当描述。
func parseSkillFile(fileName, content string) (*skillDoc, error) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.TrimPrefix(content, "\ufeff")
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, errors.New("文件是空的")
	}

	meta := map[string]any{}
	body := trimmed
	if strings.HasPrefix(trimmed, "---\n") || trimmed == "---" {
		rest := strings.TrimPrefix(trimmed, "---\n")
		if idx := strings.Index(rest, "\n---"); idx >= 0 {
			front := rest[:idx]
			body = strings.TrimSpace(strings.TrimPrefix(rest[idx+4:], "\n"))
			if err := yaml.Unmarshal([]byte(front), &meta); err != nil {
				return nil, errors.New("frontmatter 不是合法的 YAML：" + err.Error())
			}
		}
	}

	doc := &skillDoc{Body: body}
	doc.Name = metaString(meta, "name")
	if doc.Name == "" {
		doc.Name = strings.TrimSuffix(path.Base(fileName), path.Ext(fileName))
	}
	doc.Name = strings.TrimSpace(doc.Name)
	if doc.Name == "" {
		return nil, errors.New("未能确定技能名，请在 frontmatter 里写上 name")
	}

	doc.Description = metaString(meta, "description")
	doc.Version = metaString(meta, "version")
	if doc.Description == "" {
		doc.Description = firstLine(body)
	}
	doc.Description = truncateRunes(doc.Description, 200)
	doc.Version = truncateRunes(doc.Version, 32)

	// name/description/version 已单独成列，其余 frontmatter 键（allowed-tools 等）留在配置里
	delete(meta, "name")
	delete(meta, "description")
	delete(meta, "version")
	if len(meta) > 0 {
		if raw, err := json.Marshal(meta); err == nil {
			doc.ConfigJSON = string(raw)
		}
	}

	if strings.TrimSpace(doc.Body) == "" {
		return nil, errors.New("SKILL.md 没有正文内容")
	}
	return doc, nil
}

func metaString(meta map[string]any, key string) string {
	raw, ok := meta[key]
	if !ok || raw == nil {
		return ""
	}
	switch v := raw.(type) {
	case string:
		return strings.TrimSpace(v)
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			if s := strings.TrimSpace(strings.TrimSpace(toString(item))); s != "" {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, ", ")
	default:
		return strings.TrimSpace(toString(raw))
	}
}

func toString(v any) string {
	if raw, err := json.Marshal(v); err == nil {
		s := string(raw)
		if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		return s
	}
	return ""
}

// firstLine 取正文里第一段有意义的文字，用来给没有 description 的技能补一句说明。
func firstLine(body string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(line), "#>-*"))
		if line != "" {
			return line
		}
	}
	return ""
}

func truncateRunes(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit])
}
