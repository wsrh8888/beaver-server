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
	"encoding/xml"
	"errors"
	"io"
	"strconv"
	"strings"
)

// parseKnowledgeDocx 从 .docx（OOXML zip）抽取纯文本。只依赖标准库：
// 读 word/document.xml 里的 <w:t>，遇到段落 </w:p> 补换行，制表/软换行分别还原。
func parseKnowledgeDocx(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	var target io.ReadCloser
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			target, err = f.Open()
			if err != nil {
				return "", err
			}
			break
		}
	}
	if target == nil {
		return "", errors.New("docx 缺少 word/document.xml")
	}
	defer target.Close()

	dec := xml.NewDecoder(target)
	var sb, cur strings.Builder
	inT := false
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "t":
				inT = true
				cur.Reset()
			case "tab":
				sb.WriteString("\t")
			case "br", "cr":
				sb.WriteString("\n")
			}
		case xml.CharData:
			if inT {
				cur.Write(se)
			}
		case xml.EndElement:
			switch se.Name.Local {
			case "t":
				if inT {
					sb.WriteString(cur.String())
					inT = false
				}
			case "p":
				sb.WriteString("\n")
			}
		}
	}
	return normalizeKnowledgeText(sb.String()), nil
}

// parseKnowledgeXlsx 从 .xlsx（OOXML zip）抽取纯文本。只依赖标准库：
// 先读 xl/sharedStrings.xml 得到共享字符串表，再顺序遍历 xl/worksheets/sheet*.xml，
// 每个 <row> 内单元格用制表符分隔、行之间换行；t="s" 的 <v> 是共享串下标。
func parseKnowledgeXlsx(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	shared, err := readXlsxSharedStrings(zr)
	if err != nil {
		return "", err
	}

	var sheetNames []string
	for _, f := range zr.File {
		if strings.HasPrefix(f.Name, "xl/worksheets/") &&
			strings.HasSuffix(f.Name, ".xml") &&
			!strings.Contains(f.Name, "/_rels/") {
			sheetNames = append(sheetNames, f.Name)
		}
	}
	if len(sheetNames) == 0 {
		return "", errors.New("xlsx 缺少工作表")
	}

	var out strings.Builder
	for _, name := range sheetNames {
		var text string
		text, err = readXlsxSheet(zr, name, shared)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(text) != "" {
			if out.Len() > 0 {
				out.WriteString("\n")
			}
			out.WriteString(text)
		}
	}
	return normalizeKnowledgeText(out.String()), nil
}

func readXlsxSharedStrings(zr *zip.Reader) ([]string, error) {
	for _, f := range zr.File {
		if f.Name != "xl/sharedStrings.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		dec := xml.NewDecoder(rc)
		var stringsOut []string
		var siDepth int
		var cur strings.Builder
		for {
			tok, err := dec.Token()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return nil, err
			}
			switch se := tok.(type) {
			case xml.StartElement:
				if se.Name.Local == "si" {
					siDepth++
					cur.Reset()
				}
			case xml.CharData:
				if siDepth > 0 {
					cur.Write(se)
				}
			case xml.EndElement:
				if se.Name.Local == "si" && siDepth > 0 {
					stringsOut = append(stringsOut, cur.String())
					siDepth--
				}
			}
		}
		return stringsOut, nil
	}
	return []string{}, nil
}

func readXlsxSheet(zr *zip.Reader, name string, shared []string) (string, error) {
	var f *zip.File
	for _, zf := range zr.File {
		if zf.Name == name {
			f = zf
			break
		}
	}
	if f == nil {
		return "", nil
	}
	rc, err := f.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	dec := xml.NewDecoder(rc)
	var rows []string
	var rowCells []string
	var inRow, inCell bool
	var cellIsShared bool
	var cellValue strings.Builder
	var inValue bool

	flushRow := func() {
		if len(rowCells) > 0 {
			rows = append(rows, strings.Join(rowCells, "\t"))
			rowCells = nil
		}
	}
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		switch se := tok.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "row":
				inRow = true
				rowCells = nil
			case "c":
				inCell = true
				cellIsShared = false
				cellValue.Reset()
				for _, attr := range se.Attr {
					if attr.Name.Local == "t" && attr.Value == "s" {
						cellIsShared = true
					}
				}
			case "v", "t":
				if inCell {
					inValue = true
					cellValue.Reset()
				}
			}
		case xml.CharData:
			if inValue {
				cellValue.Write(se)
			}
		case xml.EndElement:
			switch se.Name.Local {
			case "v", "t":
				if inValue {
					inValue = false
				}
			case "c":
				if inCell {
					val := strings.TrimSpace(cellValue.String())
					if cellIsShared {
						if idx, convErr := strconv.Atoi(val); convErr == nil && idx >= 0 && idx < len(shared) {
							val = shared[idx]
						}
					}
					rowCells = append(rowCells, val)
					inCell = false
				}
			case "row":
				if inRow {
					flushRow()
					inRow = false
				}
			}
		}
	}
	flushRow()
	return strings.Join(rows, "\n"), nil
}

func normalizeKnowledgeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	clean := make([]string, 0, len(lines))
	for _, ln := range lines {
		if trimmed := strings.TrimRight(ln, " \t"); trimmed != "" {
			clean = append(clean, trimmed)
		} else if len(clean) > 0 && clean[len(clean)-1] != "" {
			clean = append(clean, "")
		}
	}
	for len(clean) > 0 && clean[len(clean)-1] == "" {
		clean = clean[:len(clean)-1]
	}
	return strings.Join(clean, "\n")
}
