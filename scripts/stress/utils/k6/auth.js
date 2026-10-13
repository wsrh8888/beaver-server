/*
 * k6 真实业务 flow：auth_api 的注册与登录。
 *
 * 对应 auth_api 接口，给 single/scenarios 压测复用，同一条链路不在脚本里重写。
 * 走真实 HTTP 接口，不绕过校验。
 */

import http from 'k6/http';

const UA = 'BeaverDesktop/1.0 (Windows)';

// 邮箱密码登录，返回 { token, userId } 或 null。
//   - base：auth_api 地址
//   - email/password
//   - deviceId：设备标识
export function login(base, email, password, deviceId) {
  const res = http.post(
    `${base}/api/auth/auth_public/v1/email_password_login`,
    JSON.stringify({ email, password }),
    {
      headers: {
        'Content-Type': 'application/json',
        'deviceId': deviceId,
        'User-Agent': UA,
      },
    }
  );
  if (res.status !== 200) return null;
  try {
    const b = JSON.parse(res.body);
    if (b.code !== 0 || !b.result || !b.result.token) return null;
    return { token: b.result.token, userId: b.result.userId };
  } catch (e) {
    return null;
  }
}
