/*
 * 单接口压测：chat_api 发送消息
 *
 * 接口：POST {CHAT_API_BASE}/api/chat/v1/sendMsg
 * 服务：chat_api（默认端口 20300）
 * 鉴权：只需 Beaver-User-Id header（见 chat_api/chat.go，无 auth 中间件）
 *
 * flow 复用：压测动作调 utils/k6/chat.js 的 sendMsg()，与 setup 走同一条真实链路。
 * conversationId 从 setup/friends 产出的 friends.json 取（会话随好友关系自动创建）。
 *
 * 运行（在 scripts/stress/single/chat 目录下）：
 *   k6 run --out influxdb=http://localhost:8086/k6 send_msg.js
 *   服务地址从 .env 读，无需手动设环境变量
 *
 * 前置：先跑 setup/users + setup/tokens + setup/friends 生成 friends.json
 */

import { check } from 'k6';
import { loadAccounts, loadFriendships, think, checkOk } from '../../utils/k6/shared.js';
import { load as loadConfig } from '../../utils/k6/config.js';
import { sendMsg } from '../../utils/k6/chat.js';

const { chatBase } = loadConfig();
const accounts = loadAccounts();
const friendships = loadFriendships();

export const options = {
  stages: [
    { duration: '3m', target: 50 },
    { duration: '3m', target: 50 },
    { duration: '1m', target: 100 },
    { duration: '3m', target: 100 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    'http_req_duration': ['p(99)<1000'],
    'http_req_failed': ['rate<0.01'],
  },
};

export default function () {
  const idx = __VU % friendships.length;
  const fs = friendships[idx];
  const acct = accounts[idx];
  if (!acct.userId) { think(1, 2); return; }

  const messageId = `msg-${__VU}-${__ITER}-${Math.random().toString(36).slice(2, 10)}`;
  const res = sendMsg(chatBase, acct.userId, fs.conversationId, messageId, 'stress test message');
  const ok = check(res, { '状态200': (r) => r.status === 200, 'code=0': checkOk });
  if (!ok) {
    console.error(`sendMsg 失败: status=${res.status} body=${res.body}`);
  }
  think(1, 2);
}
