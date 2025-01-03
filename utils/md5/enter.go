package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
