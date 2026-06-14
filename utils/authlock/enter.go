package authlock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/go-redis/redis"
)

const (
	maxFailures  = 5
	lockDuration = 15 * time.Minute
	failWindow   = 15 * time.Minute
)

var (
	ErrTooManyAttempts = errors.New("尝试次数过多，请15分钟后再试")
	ErrInternal        = errors.New("服务内部异常")
	log                = logger.New("authlock")
)

func CheckLocked(ctx context.Context, redisClient *redis.Client, lockKey string) error {
	exists, err := redisClient.Exists(lockKey).Result()
	if err != nil {
		log.Error(model.LogMsg{
			Text: "检查锁定状态失败",
			Data: map[string]interface{}{"lockKey": lockKey, "err": err.Error()},
		})
		return ErrInternal
	}
	if exists > 0 {
		return ErrTooManyAttempts
	}
	return nil
}

func RecordFailure(ctx context.Context, redisClient *redis.Client, failKey, lockKey, event, scope string) error {
	count, err := redisClient.Incr(failKey).Result()
	if err != nil {
		log.Error(model.LogMsg{
			Text: "递增失败计数失败",
			Data: map[string]interface{}{"failKey": failKey, "err": err.Error()},
		})
		return nil
	}
	if count == 1 {
		_ = redisClient.Expire(failKey, failWindow).Err()
	}

	log.Warn(model.LogMsg{
		Text: failureText(event),
		Data: map[string]interface{}{
			"scope": scope,
			"count": count,
		},
	})

	if count >= maxFailures {
		_ = redisClient.Set(lockKey, "1", lockDuration).Err()
		_ = redisClient.Del(failKey).Err()
		log.Error(model.LogMsg{
			Text: lockedText(event),
			Data: map[string]interface{}{
				"scope": scope,
			},
		})
		return ErrTooManyAttempts
	}
	return nil
}

func failureText(event string) string {
	switch event {
	case "login":
		return "登录失败"
	case "verify_code":
		return "验证码校验失败"
	default:
		return "认证失败"
	}
}

func lockedText(event string) string {
	switch event {
	case "login":
		return "登录已锁定"
	case "verify_code":
		return "验证码校验已锁定"
	default:
		return "认证已锁定"
	}
}

func ClearFailures(redisClient *redis.Client, failKey, lockKey string) {
	_ = redisClient.Del(failKey, lockKey).Err()
}

func LoginFailKey(account string) string {
	return fmt.Sprintf("auth_login_fail:%s", account)
}

func LoginLockKey(account string) string {
	return fmt.Sprintf("auth_login_lock:%s", account)
}

func VerifyFailKey(scope, identifier string) string {
	return fmt.Sprintf("auth_verify_fail:%s:%s", scope, identifier)
}

func VerifyLockKey(scope, identifier string) string {
	return fmt.Sprintf("auth_verify_lock:%s:%s", scope, identifier)
}

func VerifyStoredCode(ctx context.Context, redisClient *redis.Client, codeKey, scope, subject, inputCode string) error {
	lockKey := VerifyLockKey(scope, subject)
	if err := CheckLocked(ctx, redisClient, lockKey); err != nil {
		return err
	}

	storedCode, err := redisClient.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}
	if storedCode != inputCode {
		failKey := VerifyFailKey(scope, subject)
		if lockErr := RecordFailure(ctx, redisClient, failKey, lockKey, "verify_code", scope); lockErr != nil {
			return lockErr
		}
		return fmt.Errorf("验证码错误")
	}

	ClearFailures(redisClient, VerifyFailKey(scope, subject), lockKey)
	_ = redisClient.Del(codeKey).Err()
	return nil
}
