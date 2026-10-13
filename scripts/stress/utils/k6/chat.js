/*
 * k6 真实业务 flow：chat_api 的发消息与会话列表。
 * 对应 chat_api 接口，给 single/scenarios 压测复用。
 * 鉴权只需 Beaver-User-Id header（见 chat_api/chat.go，无中间件）。
 */

import http from 'k6/http';

const BASE_PATH = '/api/chat/v1';

// 发送文本消息。messageID 客户端生成（幂等）。
export function sendMsg(base, userID, conversationID, messageID, text) {
  return http.post(
    `${base}${BASE_PATH}/sendMsg`,
    JSON.stringify({
      conversationId: conversationID,
      messageId: messageID,
      msg: { type: 1, textMsg: { content: text } },
    }),
    { headers: userHeaders(userID) }
  );
}

// 最近会话列表（GET，读接口，压缓存/DB 查询）
export function recentList(base, userID, page = 1, limit = 20) {
  return http.get(
    `${base}${BASE_PATH}/getRecentChatList?page=${page}&limit=${limit}`,
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
