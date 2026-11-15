# 踩坑清单

路由段数规则（5 段为准、消灭 6 段、qrcode/org 现状偏差）见 `api-conventions.md`，此处不重复。

## CRLF / 批量编辑

- 仓库文件多为 **CRLF**：python/sed 的 `\n` 匹配会静默失败。用 `splitlines(keepends=True)` 行过滤、正则配 `\r?\n`、跳过块终止条件用 `line.rstrip() == "}"`。
- bash 传 `python -c` 时 `\n` 被损坏成字面 `/n`：换行用 `chr(10)` 或写文件。sed 替换侧 `&` 需转义 `\&`。
- MSYS 路径 `/f/...` Windows Python 打不开，用 `F:/...`；`python /dev/stdin` 不可用。
- 收费仓在 `F:/code/mine/github/beaver-server-package/` 下；本仓 import 是短名 `beaver/...`，grep 别用全 module 名（假阴性）。

## goctl

- api 生成：`goctl api go -api xxx.api -dir . --home ../../../template`（本仓自带模板，response 信封由模板+中间件统一包裹，**不手改 handler**）。
- `.api` 不接受 C 风格块注释 `/** */`（用 `//`）；同文件 service 名唯一；新增 service 名可能误建额外 main 文件（删本次误生成的，不删已有）。
- `types.go` 撞名 `duplicate type expression` 的根因是接口间共享类型——按规范每个 `@handler` 独立 `type()` 块；临时规避才改名。
- 改行为请改 `template/` 后重新生成，不要手改生成文件。
- rpc 生成后检查 pb 落盘目录（由 `go_package` 决定）。

## pb 类型映射

`.api` types 与 proto types 是两套 Go 类型，用 `common/convert.To[T]`（JSON 往返）转换，字段名/json tag 需逐字一致；转换失败当业务错误抛出，不要静默返回零值。常见名不一致：`.api` `UserLookup/UserList/UserUpdate/UserDelete/UserUpdateStatus/UserSearch/ResetUserPassword` → proto `UserInfo/ListUsers/UpdateUsers/DeleteUsers/UpdateUsersStatus/SearchUser/ResetPassword`。新增 DTO 先对照 `beaver-backend/types/remote/server/*.go` 的 json tag。

## 其他

- 鉴权失败用真实 401，别用 200+code=1（否则 S2S 客户端无法清票重换）。
- agent_api 用静态 `X-Internal-Token`（非 OAuth），与调用方 yaml `InternalToken` 同值——特例。
- 所有 Go 文件保留完整 `beaver-server-header` 版权头，生成的新文件也要补。
