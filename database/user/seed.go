package userseed

import (
	"beaver/app/auth/auth_models"
	"beaver/app/open/open_models"
	"beaver/app/user/user_models"
	"beaver/utils/pwd"
	"fmt"
	"log"

	"gorm.io/gorm"
)

const (
	defaultUserEmail    = "751135385@qq.com"
	defaultUserPassword = "e10adc3949ba59abbe56e057f20f883e" // MD5("123456")
	defaultUserID       = "100000"
	defaultNickName     = "默认开发者"
)

// InitDefaultUser 初始化默认测试账号（用户 + 密码 + 已审核开发者）
func InitDefaultUser(userDB, authDB, openDB *gorm.DB) error {
	var user user_models.UserModel
	err := userDB.Where("email = ?", defaultUserEmail).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = user_models.UserModel{
			UserID:   defaultUserID,
			UserType: user_models.UserTypeNormal,
			NickName: defaultNickName,
			Email:    defaultUserEmail,
			Source:   user_models.SourceEmail,
			Status:   1,
			Version:  1,
		}
		if err := userDB.Create(&user).Error; err != nil {
			return fmt.Errorf("创建默认用户失败: %w", err)
		}
		log.Printf("创建默认用户成功: userId=%s, email=%s", user.UserID, user.Email)
	} else if err != nil {
		return fmt.Errorf("查询默认用户失败: %w", err)
	} else {
		log.Printf("默认用户已存在: userId=%s, email=%s", user.UserID, user.Email)
	}

	var credential auth_models.AuthCredentialModel
	err = authDB.Where("user_id = ?", user.UserID).First(&credential).Error
	if err == gorm.ErrRecordNotFound {
		credential = auth_models.AuthCredentialModel{
			UserID:   user.UserID,
			Password: pwd.HahPwd(defaultUserPassword),
		}
		if err := authDB.Create(&credential).Error; err != nil {
			return fmt.Errorf("创建默认用户凭证失败: %w", err)
		}
		log.Printf("创建默认用户凭证成功: userId=%s", user.UserID)
	} else if err != nil {
		return fmt.Errorf("查询默认用户凭证失败: %w", err)
	} else {
		log.Printf("默认用户凭证已存在: userId=%s", user.UserID)
	}

	var developer open_models.OpenDeveloper
	err = openDB.Where("user_id = ?", user.UserID).First(&developer).Error
	if err == gorm.ErrRecordNotFound {
		developer = open_models.OpenDeveloper{
			UserID:      user.UserID,
			RealName:    defaultNickName,
			CompanyName: "Beaver",
			Email:       defaultUserEmail,
			Description: "系统初始化默认开发者",
			Status:      1,
		}
		if err := openDB.Create(&developer).Error; err != nil {
			return fmt.Errorf("创建默认开发者失败: %w", err)
		}
		log.Printf("创建默认开发者成功: userId=%s", user.UserID)
	} else if err != nil {
		return fmt.Errorf("查询默认开发者失败: %w", err)
	} else if developer.Status != 1 {
		if err := openDB.Model(&developer).Update("status", 1).Error; err != nil {
			return fmt.Errorf("更新默认开发者状态失败: %w", err)
		}
		log.Printf("默认开发者已审核通过: userId=%s", user.UserID)
	} else {
		log.Printf("默认开发者已存在: userId=%s", user.UserID)
	}

	return nil
}
