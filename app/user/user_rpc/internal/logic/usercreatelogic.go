package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/utils/pwd"
	utils "beaver/utils/rand"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	userIDKey = "user_id_counter"
	minUserID = 100000
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// generateUserID 生成递增用户ID
func (l *UserCreateLogic) generateUserID() (string, error) {
	// 使用Redis INCR原子性递增
	result := l.svcCtx.Redis.Incr(userIDKey)
	if result.Err() != nil {
		return "", fmt.Errorf("生成用户ID失败: %v", result.Err())
	}

	id := result.Val()

	// 如果小于最小值，设置为最小值
	if id < minUserID {
		l.svcCtx.Redis.Set(userIDKey, minUserID, 0)
		id = minUserID
	}

	return strconv.FormatInt(id, 10), nil
}

func (l *UserCreateLogic) UserCreate(in *user_rpc.UserCreateReq) (*user_rpc.UserCreateRes, error) {
	// 验证必填字段
	if in.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	if in.Phone == "" && in.Email == "" {
		return nil, errors.New("手机号或邮箱至少需要提供一个")
	}

	// 检查用户是否已存在
	var user user_models.UserModel
	var err error

	if in.Phone != "" && in.Email != "" {
		// 手机号和邮箱都提供，检查是否已存在
		err = l.svcCtx.DB.Take(&user, "phone = ? OR email = ?", in.Phone, in.Email).Error
	} else if in.Phone != "" {
		// 只提供手机号
		err = l.svcCtx.DB.Take(&user, "phone = ?", in.Phone).Error
	} else {
		// 只提供邮箱
		err = l.svcCtx.DB.Take(&user, "email = ?", in.Email).Error
	}

	if err == nil {
		return nil, errors.New("用户已存在")
	}

	// 加密密码
	hashedPassword := pwd.HahPwd(in.Password)

	// 生成随机昵称
	nickname := utils.GenerateRandomString(8)
	if in.NickName != "" {
		nickname = in.NickName
	}

	// 生成6位数字递增用户ID
	userID, err := l.generateUserID()
	if err != nil {
		logx.Errorf("生成用户ID失败: %v", err)
		return nil, errors.New("生成用户ID失败")
	}

	// 获取新版本号（用户独立递增，从1开始）
	version := l.svcCtx.VersionGen.GetNextVersion("users", "uuid", userID)
	if version == -1 {
		logx.Errorf("获取版本号失败")
		return nil, errors.New("获取版本号失败")
	}
	logx.Infof("获取用户版本号: userID=%s, version=%d", userID, version)

	user = user_models.UserModel{
		UUID:     userID,
		Password: hashedPassword,
		Email:    in.Email,
		Phone:    in.Phone,
		Source:   in.Source,
		NickName: nickname,
		Abstract: "",
		Version:  version,
	}

	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	// 记录用户创建的变更日志
	l.recordUserCreateLog(user.UUID, version)

	logx.Infof("用户创建成功: %s, version: %d", user.UUID, version)

	return &user_rpc.UserCreateRes{
		UserID: user.UUID,
	}, nil
}

// recordUserCreateLog 记录用户创建日志
func (l *UserCreateLogic) recordUserCreateLog(userID string, version int64) {
	// 创建用户创建的变更日志
	changeLog := user_models.UserChangeLogModel{
		UserID:     userID,
		ChangeType: "create",
		NewValue:   "",
		ChangeTime: time.Now().Unix(),
		Version:    version,
	}

	if err := l.svcCtx.DB.Create(&changeLog).Error; err != nil {
		logx.Errorf("记录用户创建日志失败: userID=%s, error=%v", userID, err)
	} else {
		logx.Infof("用户创建日志记录成功: userID=%s, version=%d", userID, version)
	}
}
