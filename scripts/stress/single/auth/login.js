/*
 * 单接口压测：auth_api 邮箱密码登录
 *
 * 接口：POST {AUTH_API_BASE}/api/auth/auth_public/v1/email_password_login
 * 服务：auth_api（默认端口 20100）
 *
 * 响应结构：{ code:0, msg:"", result:{ token, userId } }
 * 鉴权：请求头需要 User-Agent: BeaverDesktop/1.0 (Windows) 才能通过设备类型校验
 *
 * 运行（在 scripts/stress/single/auth 目录下）：
 *   $env:AUTH_API_BASE="http://你的IP:20100"
 *   k6 run --out json=report.json login.js
 *
 * 前置：先在 scripts/stress/seed 跑 users + tokens 生成 ../../../accounts.json
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

const accounts = new SharedArray('accounts', function () {
  return JSON.parse(open('../../accounts.json'));
});

const BASE_URL = __ENV.AUTH_API_BASE;
if (!BASE_URL) {
  throw new Error('请设置环境变量 AUTH_API_BASE，例如：$env:AUTH_API_BASE="http://IP:20100"');
}

// 阶梯加压：10 → 30 → 50 并发（温和测试，看正常表现）
export const options = {
  stages: [
    { duration: '30s', target: 10 },    // 升到 10 并发
    { duration: '1m', target: 10 },     // 保持 10 并发 1 分钟
    { duration: '30s', target: 30 },    // 升到 30
    { duration: '1m', target: 30 },     // 保持 30
    { duration: '30s', target: 50 },    // 升到 50
    { duration: '2m', target: 50 },     // 保持 50 并发 2 分钟
    { duration: '30s', target: 0 },      // 降到 0
  ],
  thresholds: {
    'http_req_duration': ['p(99)<2000'],
    'http_req_failed': ['rate<0.05'],
  },
};

export default function () {
  const acct = accounts[__VU % accounts.length];
  const res = http.post(
    `${BASE_URL}/api/auth/auth_public/v1/email_password_login`,
    JSON.stringify({ email: acct.email, password: acct.password }),
    {
      headers: {
        'Content-Type': 'application/json',
        'deviceId': `stress-device-${__VU}`,
        'User-Agent': 'BeaverDesktop/1.0 (Windows)',
      },
    }
  );
  check(res, {
    '状态200': (r) => r.status === 200,
    'code=0': (r) => {
      try { return JSON.parse(r.body).code === 0; } catch (e) { return false; }
    },
    '返回token': (r) => {
      try { return !!JSON.parse(r.body).result.token; } catch (e) { return false; }
    },
  });
  sleep(Math.random() + 1);
}
