package main

import (
	"beaver/app/chat/chat_models"
	"beaver/app/file/file_models"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_models"
	"beaver/core"
	"flag"
	"fmt"
)

type Options struct {
	DB bool
}

// go run main.go  -db
func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()

	if opt.DB {
		db := core.InitGorm("root:123456@tcp(127.0.0.1:1800)/beaver?charset=utf8mb4&parseTime=True&loc=Local")
		err := db.AutoMigrate(
			&user_models.UserModel{},
			&friend_models.FriendModel{},
			&friend_models.FriendVerifyModel{},
			&chat_models.ChatModel{},
			&chat_models.ChatUserConversationModel{},
			&file_models.FileModel{},
		)
		if err != nil {
			fmt.Println("表结构生成失败")
		}
		fmt.Println("表结构生成成功")
	}
}
