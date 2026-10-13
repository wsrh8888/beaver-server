/*
 * k6 真实业务 flow：friend_api 的好友操作。
 * 对应 friend_api 接口，给 single/scenarios 压测复用。
 * 鉴权只需 Beaver-User-Id header（见 friend_api/friend.go，无 auth 中间件）。
 */

import http from 'k6/http';

const BASE_PATH = '/api/friend/v1';

// 发送好友验证请求。source: email/qrcode/userId
export function addFriend(base, userID, friendID, source = 'userId') {
  return http.post(
    `${base}${BASE_PATH}/add_friend`,
    JSON.stringify({ friendId: friendID, source }),
    { headers: userHeaders(userID) }
  );
}

// 拉好友列表（GET）
export function friendList(base, userID, page = 1, limit = 20) {
  return http.get(
    `${base}${BASE_PATH}/friend_list?page=${page}&limit=${limit}`,
    { headers: userHeaders(userID) }
  );
}

// 拉好友详情（GET，拿 conversationId）
export function friendInfo(base, userID, friendID) {
  return http.get(
    `${base}${BASE_PATH}/friend_info?friendId=${friendID}`,
    { headers: userHeaders(userID) }
  );
}

// 业务接口请求头：只需 Beaver-User-Id（业务 API 不校验 token）
function userHeaders(userID) {
  return {
    'Content-Type': 'application/json',
    'Beaver-User-Id': userID,
  };
}
