package sse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Emitter 向客户端推一帧 SSE 事件。
type Emitter func(event string, payload any) error

// WriteHeaders 设置 SSE 响应头（须在首次写 body 前调用）。
func WriteHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

// Upgrade 校验 Flusher 并写好 SSE 头，返回可复用的 Emitter。
func Upgrade(w http.ResponseWriter) (Emitter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming unsupported")
	}
	WriteHeaders(w)
	return func(event string, payload any) error {
		return Write(w, flusher, event, payload)
	}, nil
}

// Write 写一帧：event + data(JSON) + flush。
func Write(w http.ResponseWriter, flusher http.Flusher, event string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}
