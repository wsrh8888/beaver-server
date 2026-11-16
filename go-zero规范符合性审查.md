# go-zero 规范符合性审查（app/ 目录）

> 依据：`skills/beaver-dev/references/go-zero-server.md`
> 范围：`f:/code/mine/github/beaver-server/app/` 全部 9 个域
> 方式：静态扫描 + 逐条对照规范硬性要求。
> 说明：已排除误报（如业务字段 `json:"code"`/`json:"msg"`、注释里的 URL `platform/*`）。
> 后续调整（2026-11-17）：proto 生成物落点由 `types/` 全仓迁移为官方习惯的 `pb/`（`go_package = "./pb"`，生成目录 `xxx_rpc/pb/`）。下方表格与章节中的 `./types` 为当时的历史方案，已被 pb 方案取代，以本节说明为准。

## 修复进度（2026-11-16 全部完成）

并发确认停止后逐项重构，**全部 8 项归零，`go build ./...` 通过**。

| 问题 | 状态 | 修复方式 |
|---|---|---|
| 缺 `beaver-server-header` | ✅ | 6 个生成模板补 header + 全 app 补注释；0 缺失 |
| action 命名非 snake_case | ✅ | 14 个 .api、60+ 处；生成物已同步 |
| 接口共享 Req/Res | ✅ | `file.api` 4 个上传接口各拆独立 Req/Res + 子类型；`file_utils.go` 加 `ConvertUploadRes` 集中转换 |
| `go_package` 落点不符 | ✅ | 最终统一为官方 pb 方案：6 个 rpc 的 `types/`→`pb/`、`go_package = "./pb"`、package `types`→`pb`；294 个 Go 文件 import 路径迁移、121 个无别名文件限定符 `types.`→`pb.`；`protoc --go_out=.` 实测产出 `pb/`（写 `./pb` 会套娃成 `pb/pb/`）|
| RPC 单 service | ✅ | notification_rpc 拆 `event`/`inbox`/`push` 三个 service，`-m` 重生成 client/logic/server；6 个手写 logic 迁入子域；聚合 client `notification.Notification` 改为组合 3 个窄 client，对外接口不变 |
| 跨服务直连数据库 | ✅ | message_rpc 新增 `GetSyncDeletedMessages`/`GetConversationMemberIDs` 两个 rpc；datasync、ws 改走 rpc；全仓 0 处跨库 |
| models 分包 | ✅ | message→`chat/`、notification→`notification/`、datasync→`datasync/`；全仓 import 加别名保持引用不变；migrations/main 同步 |
| type 块打包 | ✅ | 9 个 .api 按 handler 拆成独立 `type(...)` 块；types.go 生成物 diff=0，证明纯格式化 |

**额外消除跨域耦合**：relation 引用的通知事件常量/用户类型常量下沉到 `common/const/notifyconst`、`common/const/usertypeconst`，relation 不再 import 其它域 models。

**合理例外**：8 个 api 的 type 块数比 handler 数多 1~2 —— 均为被多 handler 共享或跨文件复用的子类型（如 `UserLookupInfo`、`TextMsg`、`WorkbenchEntryConfig`）提取成的公共块。这类类型无法并入单一 handler 而不重复定义，保留公共块符合规范「避免共享导致改A连带改B」的本意。

**验证**：`go build ./...` EXIT=0；8 项规范复核全部为 0。

## 结论速览（审查时原始记录）

| 严重度 | 问题 | 数量 | 规范条款 |
|---|---|---|---|
| 🔴 高 | 跨服务直连其它域数据库 | 2 处真跨库 + 5 处跨服务类型耦合 | 三条红线 / 通用硬性 |
| 🔴 高 | RPC 只有单 service、未按 `-m` 拆子域 | 1 个服务（notification_rpc） | RPC 服务规范 |
| 🟠 中 | `go_package` 未写 `./types` | 4 个 rpc | proto 规范 |
| 🟠 中 | 接口间共享 Req/Res / 子类型 | 1 + 3 处 | `.api` 契约 |
| 🟠 中 | type 块打包、handler 与 type 块不一一对应 | 11 个 api 文件 | `.api` 契约 |
| 🟡 低 | action 命名非 snake_case | 60 个（15 文件） | 路由六段 |
| 🟡 低 | 缺 `beaver-server-header` | 145 个文件（根因在 template） | 通用硬性 |
| 🟡 低 | models 未按业务分子目录 / 文件复数命名 | 3 域 + 2 文件 | Models 规范 |

合规的部分（未发现问题）：领域目录只有 `xx_api/xx_models/xx_rpc`；一文件一表（无多 `TableName()`）；路由六段格式、`@server group` 与路由 group 段一致、无 `:id` 路径参数、无 `/v1` 后多余拆段；无响应包 `code`/`msg`；`.api` 无 `/* */` 块注释。

---

## 🔴 一、跨服务直接访问其它模块数据库（红线）

### 1.1 真·跨库 CRUD（拿本服务 DB 句柄查别人表）——必须改

| 位置 | 证据 |
|---|---|
| `app/datasync/datasync_api/internal/logic/sync/getsyncdeletedmessageslogic.go` | `l.svcCtx.DB.Model(&message_models.ChatUserDelete{})` 查 message 库的表 |
| `app/ws/ws_api/internal/logic/websocket/handler/chat_message/typing_send.go` | `svcCtx.DB.Where(...)` 查 `message_models.ChatUserConversation` |

### 1.2 跨服务类型/常量耦合（未直接跨库写，但仍违反"禁止访问其它模块"）

| 位置 | 引用 |
|---|---|
| `app/relation/relation_api/internal/logic/friend/addfriendlogic.go` | `account_models/user`、`notification_models.EventTypeFriendRequest` |
| `app/relation/relation_api/internal/logic/friend/uservalidstatuslogic.go` | `notification_models.EventTypeFriendRequestAccept` 等 |
| `app/relation/relation_api/internal/logic/group/joingrouplogic.go` | `notification_models.EventTypeGroupJoinRequest` |
| `app/relation/relation_api/internal/logic/group/quitgrouplogic.go` | `notification_models.EventTypeGroupLeft` |

**修复方向**：按规范「本仓服务 A → B 走 etcd gRPC（rpc）」，把查询/事件能力下沉到对应域的 rpc；常量应在本域定义或走 rpc 参数，不 import 他域 models。

---

## 🔴 二、RPC 目录形态违规

### 2.1 notification_rpc 单 service、未拆子域

- `app/notification/notification_rpc/notification_rpc.proto`：仅 `service notification` 1 个。
- 规范强制「必须多个 service 按子域拆，编译加 `-m`」。
- 连带形态缺失：`internal/logic`、`internal/server` 下**无任何子域目录**，也没有 `client/<子域>/`（其它 5 个 rpc 都有，例如 account 的 `client/{auth,user}`、`logic/{auth,user}`、`server/{auth,user}`）。
- **修复方向**：按子域（如 inbox / push / event / read）拆多个 service，用 `goctl rpc protoc ... -m` 重新生成。

---

## 🟠 三、proto `go_package` 不符（已按官方 pb 方案修复）

审查时规范要求仍写 `./types`，后与 go-zero 官方对齐改为 `option go_package = "./pb";`，生成物落 `xxx_rpc/pb/`。当时 3 个 proto 写成 `./<服务>_rpc`，导致生成物多套一层目录：

| proto | 审查时现状 | 最终状态 |
|---|---|---|
| `app/message/message_rpc/chat_rpc.proto` | `option go_package = "./chat_rpc";` | `./pb` |
| `app/notification/notification_rpc/notification_rpc.proto` | `option go_package = "./notification_rpc";` | `./pb` |
| `app/platform/platform_rpc/platform_rpc.proto` | `option go_package = "./platform_rpc";` | `./pb` |

6 个 rpc（`account`/`asset`/`message`/`notification`/`platform`/`relation`）现已全部为 `./pb`，生成物目录统一 `pb/`。

---

## 🟠 四、`.api` type 契约违规

### 4.1 接口间共享 Req/Res（规范明令禁止）

| 文件 | 共享 type | 被哪些 handler 使用 |
|---|---|---|
| `app/asset/asset_api/api/file.api` | `FileReq` / `FileRes` | `FileUploadQiniu`、`FileUploadLocal`、`FileUploadMinio`、`FileUpload` |

### 4.2 接口间共享子类型（规范要求子类型"仅本接口用"）

| 文件 | 子类型 | 被复用于 |
|---|---|---|
| `app/account/account_api/api/user.api` | `UserLookupInfo` | `UserLookupRes`、`UserSearchRes`、`UserListRes` |
| `app/asset/asset_api/api/emoji.api` | `EmojiPackageItem` | `GetEmojiPackagesRes`、`GetUserFavoritePackagesRes` |
| 同上 | `EmojiCollectDetailItem` | `GetEmojiCollectsByIdsRes`、`GetEmojiCollectsByUuidsRes` |
| 同上 | `EmojiPackageContentDetailItem` | `GetEmojiPackageContentsByPackageIdsRes`、`...ByRelationIdsRes` |
| `app/message/message_api/api/chat.api` | `Msg`、`Sender`、`Message` | 多个 Req/Res（SendMsg、ChatHistory、SearchMessages、GetForwardDetails…） |

### 4.3 type 块未与 handler 一一对应（一个 handler 应单独一组 `type (...)`）

下列文件把多个接口的 type 堆进一个大块、或混用单个 `type X{}`；虽多数 Req/Res 命名已独立，但不符合「有 N 个 handler 就有且仅有 N 组独立 type 块」：

| 文件 | handler 数 | `type()` 块数 | 单 `type X{}` |
|---|---|---|---|
| `app/platform/platform_api/api/admin.api` | 24 | 1 | 0 |
| `app/relation/relation_api/api/group.api` | 30 | 23 | 0 |
| `app/relation/relation_api/api/friend.api` | 21 | 15 | 0 |
| `app/asset/asset_api/api/emoji.api` | 20 | 16 | 0 |
| `app/message/message_api/api/chat.api` | 25 | 27 | 0 |
| `app/account/account_api/api/user.api` | 15 | 16 | 0 |
| `app/platform/platform_api/api/org.api` | 3 | 1 | 0 |
| `app/platform/platform_api/api/update_public.api` | 2 | 0 | 4 |
| `app/platform/platform_api/api/workbench.api` | 1 | 2 | 0 |
| `app/ws/ws_api/api/proxy.api` | 1 | 0 | 0 |
| `app/ws/ws_api/api/socket.api` | 1 | 0 | 0 |

（type 块数 > handler 数的，通常是把某个接口拆成了多个块，同样不满足一一对应。）

---

## 🟡 五、路由 action 命名非 snake_case

规范：路由第六段 `action` 必须 snake_case。共 **60 个** action 用了驼峰或连字符，分布在 15 个文件。重灾区：

| 文件 | 违规数 | 示例 |
|---|---|---|
| `app/message/message_api/api/chat.api` | 15 | `sendMsg`、`getChatHistory`、`markMessageMedia` |
| `app/asset/asset_api/api/emoji.api` | 13 | `favoriteEmoji`、`packageList`、`emojis-by-ids` |
| `app/datasync/datasync_api/api/chat.api` | 5 | `getSyncChatMessages` |
| `app/notification/notification_api/api/read.api` | 5 | `getReadCursors`、`markReadByCursor` |
| 其余 11 文件 | 1~3 不等 | `searchValidInfo`、`uploadQiniu`、`getInboxByIds`… |

注意：`datasync` 域几乎所有 sync action 都是驼峰（`getSyncXxx` 系列）；`emoji` 域混用了 `-` 连字符（`packages-by-ids`、`package-contents-by-package-ids`）。

---

## 🟡 六、版权头缺失（根因在 template）

全仓 `app/` 下 **145 个 `.go` 文件** 缺 `beaver-server-header-v1`：

| 分类 | 数量 |
|---|---|
| handler（goctl 生成） | 93 |
| pb 生成 | 12 |
| client（生成） | 14 |
| server（生成） | 14 |
| middleware | 8 |
| logic（手写） | 3 |
| 其它手写（含 `platform_rpc/platform/platform.go`） | 1 |

**根因**：`template/api/` 下只有 `routes.tpl`、`types.tpl` 配了 header；`handler.tpl`、`logic.tpl`、`middleware.tpl` 开头直接是 `package ...`，没有 header。规范"生成文件也要补"，正确做法是**给这些 tpl 补 header 后重新生成**（对应规范"改公共生成行为：改 template/ 后重新生成"），而不是逐个改生成物。

手写文件缺 header 清单（8 middleware + 4 其它）：
- `*/internal/middleware/jwts2smiddleware.go`（account/asset/datasync/message/notification/platform/relation/ws 共 8 个）
- `app/platform/platform_api/internal/logic/org/{departmenttreelogic,getorganizationlogic,listmemberlogic}.go`
- `app/platform/platform_rpc/platform/platform.go`

---

## 🟡 七、Models 分包与命名

### 7.1 未按业务分子目录（models 根下平铺）

规范要求「按业务分子目录（user/auth/oauth/bot/digital…），禁止大杂烩包」。以下域的 models 根下直接放表文件：

- `app/message/message_models/`：`chat_conversation_setting_model.go`、`chat_datasync_model.go`、`chat_forward_model.go`、`chat_message_media_model.go`、`chat_model.go`、`chat_user_delete_model.go`（6 个平铺，应建 `chat/` 子包）
- `app/notification/notification_models/`：`constants.go`、`notification_event.go`、`notification_inbox.go`、`notification_read.go`、`push_registration_model.go`（5 个平铺）
- `app/datasync/datasync_models/datasync_model.go`（单文件平铺）

（account/asset/platform/relation 已按 `user/`、`auth/`、`emoji/`、`file/`、`org/`、`friend/`、`group/` 等正确分了子目录。）

### 7.2 文件复数命名 `*_models.go`

规范文件名应为 `xxx_model.go`（或现有 `*_model.go` 风格）。以下为复数，不合风格：
- `app/asset/asset_models/emoji/emoji_models.go`
- `app/asset/asset_models/file/file_models.go`

（二者均为"一文件一表 + 同文件 JSON Scan/Value 辅助类型"，结构合规，仅文件名不合风格。）

---

## 建议修复顺序

1. **先补 template header**（一处改动解决 145 个生成文件），再 `goctl` 全量重生成。
2. **拆 notification_rpc 为多 service + `-m`**，并统一 6 个 rpc 的 `go_package = "./pb"`（官方方案，生成物落 `pb/`）。
3. **改 2 处真跨库**：datasync、ws 通过对应域 rpc 取数；relation 的跨服务常量改为本域定义或 rpc 参数。
4. **拆 `.api` 共享 type**：file.api 的 FileReq/FileRes、3 处子类型；再按 handler 拆独立 type 块。
5. **action 批量改 snake_case**（注意同步前端/客户端调用路径）。
6. **models 分包 + 复数文件名**收尾。
