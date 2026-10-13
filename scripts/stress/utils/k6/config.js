/*
 * k6 配置：从 scripts/stress/.env 读取所有服务地址，与 Go 端共用一份配置。
 *
 * k6 v2.3.0 没有 --env-file，用 init context 的 open() 自己读 .env 解析。
 * k6 v2 的 open() 相对当前模块文件（config.js）所在目录解析路径，
 * 所以 .env 路径相对 config.js（utils/k6/）固定为 ../../.env，与调用脚本所在目录无关。
 *
 * 用法（任一 k6 脚本顶部，无需关心相对路径）：
 *   import { load } from '<rel>/utils/k6/config.js';
 *   const { authBase, friendBase, chatBase } = load();
 */

// load 读 .env 并解析，返回各服务地址。
// .env 路径相对 config.js 固定，调用方无需传参。
export function load() {
  const env = parseEnv(open('../../.env'));
  return {
    authBase: mustGet(env, 'AUTH_API_BASE'),
    friendBase: mustGet(env, 'FRIEND_API_BASE'),
    chatBase: mustGet(env, 'CHAT_API_BASE'),
    userBase: env['USER_API_BASE'] || '',
    groupBase: env['GROUP_API_BASE'] || '',
    momentBase: env['MOMENT_API_BASE'] || '',
  };
}

// parseEnv 解析 .env 文本为键值 map，跳过注释和空行，trim 引号
function parseEnv(text) {
  const env = {};
  for (const line of text.split('\n')) {
    const l = line.trim();
    if (l === '' || l.startsWith('#')) continue;
    const eq = l.indexOf('=');
    if (eq < 0) continue;
    const key = l.slice(0, eq).trim();
    let val = l.slice(eq + 1).trim();
    val = val.replace(/^["']|["']$/g, ''); // 去首尾引号
    env[key] = val;
  }
  return env;
}

// mustGet 必填项，缺了直接抛错终止 k6
function mustGet(env, key) {
  const v = env[key];
  if (!v) {
    throw new Error(`.env 中 ${key} 未配置或为空`);
  }
  return v;
}
