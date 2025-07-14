package pwd

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HahPwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func CheckPad(hashPad string, pwd string) bool {
	byteHash := []byte(hashPad)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(pwd))
	if err != nil {
		fmt.Println("11111111111111111111111111111111", hashPad, pwd)
		log.Println(err)
		return false
	}
	return true
}
