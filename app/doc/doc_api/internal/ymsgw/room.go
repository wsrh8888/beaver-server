package ymsgw

import (
	"strconv"
	"strings"
)

// roomSep 是协同房间名中 docId 与 epoch 的分隔符。
// 不能用 "#"：y-websocket 客户端会把房间名拼进 URL path，"#" 会被浏览器当作
// fragment，导致服务端收不到 epoch 部分；"@" 在 path 中是安全的可打印字符。
// 房间名格式：docId@epoch，例如 doc_9f3a...@2。
const roomSep = "@"

// BuildRoom 构造带纪元的协同房间名。epoch<=0 时归一为 1。
func BuildRoom(docID string, epoch int32) string {
	if epoch <= 0 {
		epoch = 1
	}
	return docID + roomSep + EpochString(epoch)
}

// EpochString 把纪元格式化为房间名/query 中的十进制字符串。
func EpochString(epoch int32) string {
	if epoch <= 0 {
		epoch = 1
	}
	return strconv.Itoa(int(epoch))
}

// ParseRoom 从房间名拆出 docId 与 epoch。
// 兼容旧房间名（纯 docId，无分隔符）：返回 (docId, 0)，由 doc_rpc 按当前纪元处理。
func ParseRoom(room string) (string, int32) {
	if i := strings.LastIndex(room, roomSep); i > 0 {
		if n, err := strconv.Atoi(room[i+1:]); err == nil && n > 0 {
			return room[:i], int32(n)
		}
	}
	return room, 0
}
