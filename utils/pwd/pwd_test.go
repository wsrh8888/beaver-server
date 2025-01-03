package pwd

import (
	"fmt"
	"testing"
)

func TestHashPad(t *testing.T) {
	hash := HahPwd("123456")
	fmt.Println(hash)
}

func TestCheckPad(t *testing.T) {
	ok := CheckPad("$2a$10$qq.QNZaDpKsNk1qUAJqDLeFAOPTYIB8fFeNU8KojWlUQJsTRpDY7y", "123456")
	fmt.Println(ok)
}
