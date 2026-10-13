/*
 * 场景压测：模拟真实用户日常聊天
 *
 * 链路：拉好友列表 → 发消息 → 拉最近会话列表
 * 跨服务：friend_api + chat_api（登录已在 setup 拿过 token，链路里不重复登录）
 *
 * 与单接口压测的区别：
 *   - 单接口：看单个接口性能上限，数据隔离
 *   - 场景：模拟真实用户行为，多接口按比例混合，看整体链路瓶颈
 *
 * flow 复用：每个接口调 utils/k6 里的 flow（friendList / sendMsg / recentList），
 * 与单接口压测共用同一份 flow，不在脚本里重写 HTTP 细节。
 * conversationId 从 setup/friends 产出的 friends.json 取。
 *
 * 运行（在 scripts/stress/scenarios/daily_chat 目录下）：
 *   k6 run --out influxdb=http://localhost:8086/k6 daily_chat.js
 *   服务地址从 .env 读，无需手动设环境变量
 *
 * 前置：setup 全部跑完（users + tokens + friends）
 */

import { check } from 'k6';
import { loadAccounts, loadFriendships, think, checkOk } from '../../utils/k6/shared.js';
import { load as loadConfig } from '../../utils/k6/config.js';
import { friendList } from '../../utils/k6/friend.js';
import { sendMsg, recentList } from '../../utils/k6/chat.js';

const { chatBase, friendBase } = loadConfig();
const accounts = loadAccounts();
const friendships = loadFriendships();

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
  const idx = __VU % friendships.length;
  const fs = friendships[idx];
  const acct = accounts[idx];
  if (!acct.userId) { think(1, 2); return; }

  // 1. 拉好友列表（friend_api）
  const friendRes = friendList(friendBase, acct.userId, 1, 20);
  check(friendRes, { '好友列表200': (r) => r.status === 200, '好友列表code=0': checkOk });
  think(1, 2);

  // 2. 发消息（chat_api）—— 用好友关系自动产生的会话
  const sendRes = sendMsg(chatBase, acct.userId, fs.conversationId, `msg-${__VU}-${Date.now()}-${__ITER}`, 'stress test message');
  check(sendRes, { '发消息200': (r) => r.status === 200, '发消息code=0': checkOk });
  think(1, 2);

  // 3. 拉最近会话列表（chat_api，读接口）
  const recentRes = recentList(chatBase, acct.userId, 1, 20);
  check(recentRes, { '会话列表200': (r) => r.status === 200, '会话列表code=0': checkOk });
  think(1, 3);
}
