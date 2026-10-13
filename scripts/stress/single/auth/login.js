/*
 * 单接口压测：auth_api 邮箱密码登录
 *
 * 接口：POST {AUTH_API_BASE}/api/auth/auth_public/v1/email_password_login
 * 服务：auth_api（默认端口 20100）
 *
 * flow 复用：压测动作直接调 utils/k6/auth.js 的 login()，与 setup/tokens 走同一条真实链路，
 * 不在脚本里重写 HTTP 细节。
 *
 * 运行（在 scripts/stress/single/auth 目录下）：
 *   k6 run --out influxdb=http://localhost:8086/k6 login.js
 *   服务地址从 .env 读，无需手动设环境变量
 *
 * 前置：先在 scripts/stress/setup 跑 users + tokens 生成 accounts.json
 */

import { check } from 'k6';
import { loadAccounts, think } from '../../utils/k6/shared.js';
import { load as loadConfig } from '../../utils/k6/config.js';
import { login } from '../../utils/k6/auth.js';

const { authBase } = loadConfig();
const accounts = loadAccounts();

// 阶梯加压：10 → 30 → 50 并发（温和测试，看正常表现）
export const options = {
  stages: [
    { duration: '30s', target: 10 },    // 升到 10 并发
    { duration: '1m', target: 10 },      // 保持 10 并发 1 分钟
    { duration: '30s', target: 30 },     // 升到 30
    { duration: '1m', target: 30 },      // 保持 30
    { duration: '30s', target: 50 },     // 升到 50
    { duration: '2m', target: 50 },      // 保持 50 并发 2 分钟
    { duration: '30s', target: 0 },      // 降到 0
  ],
  thresholds: {
    'http_req_duration': ['p(99)<2000'],
    'http_req_failed': ['rate<0.05'],
  },
};

export default function () {
  const acct = accounts[__VU % accounts.length];
  const r = login(authBase, acct.email, acct.password, `stress-device-${__VU}`);
  check(r, {
    '登录成功': (x) => x !== null && !!x.token,
  });
  think(1, 2);
}
