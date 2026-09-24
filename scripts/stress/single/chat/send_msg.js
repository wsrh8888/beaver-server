/*
 * 单接口压测：chat_api 发送消息（TODO 占位）
 *
 * 接口：POST {CHAT_API_BASE}/api/chat/v1/sendMsg
 * 服务：chat_api（默认端口 20300）
 *
 * 前置：需要先造会话和好友关系（seed/friends.go + seed/conversations.go 完成后）
 *       账号需要有 token（seed/tokens.go）
 *
 * 运行（在 scripts/stress/single/chat 目录下）：
 *   $env:CHAT_API_BASE="http://你的IP:20300"
 *   k6 run --out json=report.json send_msg.js
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

const accounts = new SharedArray('accounts', function () {
  return JSON.parse(open('../../../accounts.json'));
});

const BASE_URL = __ENV.CHAT_API_BASE;
if (!BASE_URL) {
  throw new Error('请设置环境变量 CHAT_API_BASE，例如：$env:CHAT_API_BASE="http://IP:20300"');
}

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
  const acct = accounts[__VU % accounts.length];
  if (!acct.token) { sleep(1); return; }

  // TODO: conversationId 需要从 seed 造的会话数据里取
  const res = http.post(
    `${BASE_URL}/api/chat/v1/sendMsg`,
    JSON.stringify({
      conversationId: 'TODO_CONVERSATION_ID',
      messageId: `msg-${__VU}-${Date.now()}`,
      msg: { type: 1, textMsg: { content: 'stress test message' } },
    }),
    {
      headers: {
        'Content-Type': 'application/json',
        'Beaver-User-Id': acct.userId,
        'Authorization': acct.token,
      },
    }
  );
  check(res, { '状态200': (r) => r.status === 200 });
  sleep(Math.random() + 1);
}
