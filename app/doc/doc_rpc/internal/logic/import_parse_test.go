package logic

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// buildDocx 构造一个最小 .docx（zip 内含 word/document.xml），段落由 paras 提供。
func buildDocx(t *testing.T, paras []string) []byte {
	t.Helper()
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	sb.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`)
	for _, p := range paras {
		sb.WriteString(`<w:p><w:r><w:t xml:space="preserve">` + p + `</w:t></w:r></w:p>`)
	}
	sb.WriteString(`</w:body></w:document>`)

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("zip create: %v", err)
	}
	if _, err := w.Write([]byte(sb.String())); err != nil {
		t.Fatalf("zip write: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

func TestParseDocx_ExtractsText(t *testing.T) {
	data := buildDocx(t, []string{"第一行 协同", "Hello — 海狸 🦫", ""})
	got, err := parseDocx(data)
	if err != nil {
		t.Fatalf("parseDocx err: %v", err)
	}
	if !strings.Contains(got, "第一行 协同") {
		t.Errorf("缺少第一段文本，got=%q", got)
	}
	if !strings.Contains(got, "Hello — 海狸 🦫") {
		t.Errorf("缺少第二段文本，got=%q", got)
	}
	// 段落间应有换行分隔。
	if lines := strings.Split(got, "\n"); len(lines) < 2 {
		t.Errorf("段落未以换行分隔，lines=%d got=%q", len(lines), got)
	}
}

func TestParseDocx_MissingDocumentXML(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	_ = zw.Close()
	if _, err := parseDocx(buf.Bytes()); err == nil {
		t.Fatal("缺少 word/document.xml 应返回错误")
	}
}

func TestTextToContentJSON_ValidProseMirror(t *testing.T) {
	got := textToContentJSON("标题行\n正文一段\n\n另一段")
	var doc struct {
		Type    string `json:"type"`
		Content []struct {
			Type    string `json:"type"`
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"content"`
	}
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("content json 非法: %v (%s)", err, got)
	}
	if doc.Type != "doc" {
		t.Errorf("root type=%q want doc", doc.Type)
	}
	// 空行被 normalizeExtractedText 压缩，应至少 3 个段落且首段文本正确。
	if len(doc.Content) < 3 {
		t.Fatalf("段落数不足: %d json=%s", len(doc.Content), got)
	}
	if len(doc.Content[0].Content) == 0 || doc.Content[0].Content[0].Text != "标题行" {
		t.Errorf("首段文本不符: %+v", doc.Content[0])
	}
}

func TestTextToContentJSON_EmptyIsSkeleton(t *testing.T) {
	if got := textToContentJSON("   \n  "); got != skeletonContentJSON {
		t.Errorf("空文本应回退骨架，got=%q", got)
	}
}

func TestExtractTextFromJSON_RoundTrip(t *testing.T) {
	text := "甲\n乙丙"
	js := textToContentJSON(text)
	back := extractTextFromJSON(js)
	if !strings.Contains(back, "甲") || !strings.Contains(back, "乙丙") {
		t.Errorf("往返抽取文本丢失，back=%q", back)
	}
}
