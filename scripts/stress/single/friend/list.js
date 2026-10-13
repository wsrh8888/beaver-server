/*
 * 单接口压测：friend_api 好友列表
 *
 * 接口：GET {FRIEND_API_BASE}/api/friend/v1/friend_list
 * 服务：friend_api（默认端口 20200）
 * 鉴权：只需 Beaver-User-Id header（见 friend_api/friend.go，无 auth 中间件）
 *
 * flow 复用：压测动作调 utils/k6/friend.js 的 friendList()，与 setup 走同一条真实链路。
 * 读接口，重点看缓存命中和 DB 查询性能。
 *
 * 运行（在 scripts/stress/single/friend 目录下）：
 *   k6 run --out influxdb=http://localhost:8086/k6 list.js
 *   服务地址从 .env 读，无需手动设环境变量
 *
 * 前置：先跑 setup/users + setup/tokens + setup/friends 造好友关系
 */

import { check } from 'k6';
import { loadAccounts, think, checkOk } from '../../utils/k6/shared.js';
import { load as loadConfig } from '../../utils/k6/config.js';
import { friendList } from '../../utils/k6/friend.js';

const { friendBase } = loadConfig();
const accounts = loadAccounts();

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
  if (!acct.userId) { think(1, 2); return; }

  const res = friendList(friendBase, acct.userId, 1, 20);
  check(res, { '状态200': (r) => r.status === 200, 'code=0': checkOk });
  think(1, 2);
}
