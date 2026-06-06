package main

import (
	"beaver/database"
	"flag"
	"fmt"
)

type Options struct {
	DB bool
}

// go run main.go -db
func main() {
	var opt Options
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()

	if !opt.DB {
		return
	}

	dbMap, err := database.RunMigrations()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := database.InitDefaultUser(
		dbMap["beaver_user"],
		dbMap["beaver_auth"],
		dbMap["beaver_open"],
	); err != nil {
		fmt.Printf("默认用户初始化失败: %v\n", err)
		return
	}

	fmt.Println("所有库表结构生成完成")
}
