/*
 * k6 压测公共库：所有 k6 脚本（single/、scenarios/）共用
 *
 * 用法（在任一 k6 脚本顶部，按脚本所在目录调整相对层级）：
 *   import { loadAccounts, think, userHeaders, checkOk } from '../../utils/k6/shared.js';
 *   import { load as loadConfig } from '../../utils/k6/config.js';
 *   const { authBase, friendBase, chatBase } = loadConfig('../../.env');
 *
 * 服务地址统一从 .env 读（utils/k6/config.js），不在脚本里读 __ENV，也不手动设环境变量。
 *
 * 鉴权说明：friend_api/chat_api 等业务接口不校验 token，只需 Beaver-User-Id header
 * （见各服务 main.go 的中间件注册：仅 UA 解析，无 auth 中间件）。
 * 所以 userHeaders 只带 Beaver-User-Id + Content-Type，不带 Authorization。
 */

import { sleep } from 'k6';
import { SharedArray } from 'k6/data';

// 账号池：SharedArray 在所有 VU 间共享一份内存，数据文件只读一次。
// accounts.json 相对 shared.js（utils/k6/）固定为 ../../accounts.json，
// 调用方无需传参，避免 k6 open() 路径解析随调用脚本目录变化的问题。
export function loadAccounts() {
  return new SharedArray('accounts', function () {
    return JSON.parse(open('../../accounts.json'));
  });
}

// 好友关系池：同理，friends.json 相对 shared.js 固定路径。
export function loadFriendships() {
  return new SharedArray('friendships', function () {
    return JSON.parse(open('../../friends.json'));
  });
}

// think time：模拟真实用户两次操作间的停顿，min/max 单位秒
export function think(min = 1, max = 2) {
  sleep(Math.random() * (max - min) + min);
}

// 业务接口请求头：Beaver-User-Id + Content-Type。
// 业务 API 不校验 token（见各服务 main.go），无需 Authorization。
export function userHeaders(userID) {
  return {
    'Content-Type': 'application/json',
    'Beaver-User-Id': userID,
  };
}

// 通用结果校验：HTTP 200 + 业务 code==0
export function checkOk(r) {
  if (r.status !== 200) return false;
  try {
    return JSON.parse(r.body).code === 0;
  } catch (e) {
    return false;
  }
}

