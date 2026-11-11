/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func zipOf(t *testing.T, files map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestParseKnowledgeDocx(t *testing.T) {
	doc := `<?xml version="1.0"?><w:document xmlns:w="x"><w:body>` +
		`<w:p><w:r><w:t>你好</w:t></w:r><w:r><w:t>世界</w:t></w:r></w:p>` +
		`<w:p><w:r><w:t>第二行</w:t></w:r></w:p></w:body></w:document>`
	got, err := parseKnowledgeDocx(zipOf(t, map[string]string{"word/document.xml": doc}))
	if err != nil {
		t.Fatal(err)
	}
	if got != "你好世界\n第二行" {
		t.Fatalf("docx 解析结果不符: %q", got)
	}
}

func TestParseKnowledgeXlsx(t *testing.T) {
	shared := `<?xml version="1.0"?><sst><si><t>姓名</t></si><si><t>年龄</t></si><si><t>张三</t></si></sst>`
	sheet := `<?xml version="1.0"?><worksheet><sheetData>` +
		`<row><c t="s"><v>0</v></c><c t="s"><v>1</v></c></row>` +
		`<row><c t="s"><v>2</v></c><c><v>30</v></c></row>` +
		`</sheetData></worksheet>`
	got, err := parseKnowledgeXlsx(zipOf(t, map[string]string{
		"xl/sharedStrings.xml":     shared,
		"xl/worksheets/sheet1.xml": sheet,
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "姓名") || !strings.Contains(got, "张三") || !strings.Contains(got, "30") {
		t.Fatalf("xlsx 解析结果不符: %q", got)
	}
	if !strings.Contains(got, "\t") || !strings.Contains(got, "\n") {
		t.Fatalf("xlsx 未保留行列分隔: %q", got)
	}
}
