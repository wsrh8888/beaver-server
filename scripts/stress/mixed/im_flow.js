/*
 * 混搭压测：模拟真实用户 IM 操作链路（TODO 占位）
 *
 * 链路：登录 → 拉好友列表 → 发消息 → 拉消息列表
 * 跨服务：auth_api + friend_api + chat_api
 *
 * 与单接口压测的区别：
 *   - 单接口：看单个接口性能上限，数据隔离
 *   - 混搭：模拟真实用户行为，多个接口按比例混合，看整体链路瓶颈
 *
 * 运行（在 scripts/stress/mixed 目录下）：
 *   $env:AUTH_API_BASE="http://IP:20100"
 *   $env:CHAT_API_BASE="http://IP:20300"
 *   $env:FRIEND_API_BASE="http://IP:..."
 *   k6 run --out json=report.json im_flow.js
 *
 * 前置：seed 全部跑完（users + tokens + friends + conversations）
 */

import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

const accounts = new SharedArray('accounts', function () {
  return JSON.parse(open('../accounts.json'));
});

const AUTH = __ENV.AUTH_API_BASE;
const CHAT = __ENV.CHAT_API_BASE;
const FRIEND = __ENV.FRIEND_API_BASE;

export const options = {
  scenarios: {
    im_flow: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 50 },
        { duration: '3m', target: 50 },
        { duration: '2m', target: 100 },
        { duration: '3m', target: 100 },
        { duration: '1m', target: 0 },
      ],
    },
  },
  thresholds: {
    'http_req_duration': ['p(99)<1500'],
    'http_req_failed': ['rate<0.02'],
  },
};

export default function () {
  const acct = accounts[__VU % accounts.length];
  if (!acct.token) { sleep(1); return; }

  const headers = {
    'Content-Type': 'application/json',
    'Beaver-User-Id': acct.userId,
    'Authorization': acct.token,
  };

  // 1. 拉好友列表（friend_api）—— TODO: 确认接口路径
  // const friendRes = http.get(`${FRIEND}/api/friend/v1/list`, { headers });

  // 2. 发消息（chat_api）—— TODO: conversationId 需从 seed 取
  // const sendRes = http.post(`${CHAT}/api/chat/v1/sendMsg`, ...);

  // 3. 拉消息列表（chat_api）—— TODO

  // 占位：链路待 seed 完成后补全
  check(true, { '占位': () => true });
  sleep(Math.random() * 2 + 1);
}
