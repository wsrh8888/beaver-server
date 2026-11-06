package logic

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os/exec"
	"strings"
)

// textToContentJSON 把导入解析出的纯文本按行包装为 ProseMirror 多段落 doc JSON。
// 空文本回退骨架（单空段落），保证文档结构合法、编辑器可挂载。
func textToContentJSON(text string) string {
	text = normalizeExtractedText(text)
	if text == "" {
		return skeletonContentJSON
	}
	type pmText struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type pmBlock struct {
		Type    string   `json:"type"`
		Content []pmText `json:"content,omitempty"`
	}
	type pmDoc struct {
		Type    string    `json:"type"`
		Content []pmBlock `json:"content"`
	}
	doc := pmDoc{Type: "doc"}
	for _, ln := range strings.Split(text, "\n") {
		ln = strings.TrimSpace(ln)
		blk := pmBlock{Type: "paragraph"}
		if ln != "" {
			blk.Content = []pmText{{Type: "text", Text: ln}}
		}
		doc.Content = append(doc.Content, blk)
	}
	b, err := json.Marshal(doc)
	if err != nil {
		return skeletonContentJSON
	}
	return string(b)
}

// parseDocx 解析 .docx：archive/zip 读取 word/document.xml，抽取 <w:t> 文本，
// 段落 </w:p> 结束补换行，<w:tab> 记为制表符，<w:br/> 记为换行。纯标准库，无第三方依赖。
func parseDocx(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}
	var target io.ReadCloser
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			target = rc
			break
		}
	}
	if target == nil {
		return "", errors.New("docx 缺少 word/document.xml")
	}
	defer target.Close()

	dec := xml.NewDecoder(target)
	var sb strings.Builder
	var cur strings.Builder // 收集当前 <w:t> 文本
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
	return normalizeExtractedText(sb.String()), nil
}

// parsePdf 通过外部 pdftotext（poppler）子进程抽取正文：stdin 进、stdout 出（`pdftotext - -`）。
// 依赖缺失或执行失败时返回 parsed=false，调用方据此降级为"仅登记附件不解析正文"。
func parsePdf(pdfBin string, data []byte) (text string, parsed bool, msg string) {
	if strings.TrimSpace(pdfBin) == "" {
		pdfBin = "pdftotext"
	}
	if _, err := exec.LookPath(pdfBin); err != nil {
		return "", false, "pdftotext 不可用，已仅登记附件未解析正文"
	}
	cmd := exec.Command(pdfBin, "-", "-")
	cmd.Stdin = bytes.NewReader(data)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", false, "pdftotext 执行失败：" + firstLine(errBuf.String())
	}
	return normalizeExtractedText(out.String()), true, ""
}

// normalizeExtractedText 规整空白：统一换行、去除每行首尾空白、压缩连续空行为最多一个。
func normalizeExtractedText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	var out []string
	blank := 0
	for _, ln := range lines {
		ln = strings.TrimRight(ln, " \t")
		if ln == "" {
			blank++
			continue
		}
		if blank > 0 && len(out) > 0 {
			out = append(out, "")
		}
		blank = 0
		out = append(out, ln)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}
