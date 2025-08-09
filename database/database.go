package database

import (
	"fmt"

	"gorm.io/gorm"
)

// 主初始化函数 - 在main.go中调用
func InitAllData(db *gorm.DB) error {
	fmt.Println("=== 开始初始化所有表的默认数据 ===")

	// 按顺序初始化各表数据
	initializers := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"文件表", InitFileData},
	}

	// 逐个执行初始化
	for _, init := range initializers {
		fmt.Printf("正在初始化%s...\n", init.name)
		if err := init.fn(db); err != nil {
			fmt.Printf("%s初始化失败: %v\n", init.name, err)
			// 继续执行其他表初始化，不中断流程
		} else {
			fmt.Printf("%s初始化成功\n", init.name)
		}
	}

	fmt.Println("=== 所有表初始化完成 ===")
	return nil
}
