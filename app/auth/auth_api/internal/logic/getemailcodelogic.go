package logic

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/email"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmailCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailCodeLogic {
	return &GetEmailCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEmailCodeLogic) GetEmailCode(req *types.GetEmailCodeReq) (resp *types.GetEmailCodeRes, err error) {
	// 生成6位数字验证码
	code := email.GenerateCode()

	// 发送验证码邮件
	err = l.sendEmail(req.Email, code, req.Type)
	if err != nil {
		logx.Errorf("发送邮件失败: %v", err)
		return nil, errors.New("发送验证码失败，请稍后重试")
	}

	// 存储验证码到Redis（5分钟有效期）
	codeKey := fmt.Sprintf("email_code_%s_%s", req.Email, req.Type)
	err = l.svcCtx.Redis.Set(codeKey, code, 5*time.Minute).Err()
	if err != nil {
		logx.Errorf("存储验证码失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 设置发送频率限制（60秒）
	rateLimitKey := fmt.Sprintf("email_rate_limit_%s", req.Email)
	err = l.svcCtx.Redis.Set(rateLimitKey, "1", 60*time.Second).Err()
	if err != nil {
		logx.Errorf("设置发送频率限制失败: %v", err)
	}

	return &types.GetEmailCodeRes{
		Message: "验证码发送成功",
	}, nil
}

// 发送邮件
func (l *GetEmailCodeLogic) sendEmail(to, code, codeType string) error {
	var host, username, password string
	var port int

	// 直接使用系统配置的邮箱服务商
	host = l.svcCtx.Config.Email.QQ.Host
	port = l.svcCtx.Config.Email.QQ.Port
	username = l.svcCtx.Config.Email.QQ.Username
	password = l.svcCtx.Config.Email.QQ.Password
	logx.Infof("QQ邮箱配置: Host=%s, Port=%d, Username=%s", host, port, username)

	// 构建邮件内容
	subject := email.GetEmailSubject(codeType)
	body := email.GetEmailBody(code, codeType)

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = username // 发送方邮箱
	headers["To"] = to         // 接收方邮箱（用户传入的）
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// 构建邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	logx.Infof("准备发送邮件: From=%s, To=%s, Subject=%s", username, to, subject)

	// 根据端口选择发送方式
	var err error
	if port == 465 {
		// SSL连接
		err = l.sendEmailSSL(host, port, username, password, to, message)
	} else {
		// 普通SMTP连接
		auth := smtp.PlainAuth("", username, password, host)
		addr := fmt.Sprintf("%s:%d", host, port)
		err = smtp.SendMail(addr, auth, username, []string{to}, []byte(message))
	}

	if err != nil {
		logx.Errorf("SMTP发送失败: %v, host=%s, port=%d, username=%s, to=%s", err, host, port, username, to)
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	logx.Infof("邮件发送成功: To=%s", to)
	return nil
}

// SSL方式发送邮件
func (l *GetEmailCodeLogic) sendEmailSSL(host string, port int, username, password, to, message string) error {
	// 连接到SMTP服务器
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: host,
	})
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 认证
	auth := smtp.PlainAuth("", username, password, host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	// 设置发件人
	if err := client.Mail(username); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取邮件写入器失败: %v", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	return nil
}
