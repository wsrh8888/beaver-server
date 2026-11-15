| 服务            | API 端口 | RPC 端口 | 所属仓库 | 备注 |
|-----------------|---------:|---------:|----------|------|
| gateway_api     | 20800    | -        | beaver-server | 客户端 API 网关 |
| account         | 20100    | 30100    | beaver-server | 用户 + 认证（已合并） |
| relation        | 20200    | 30200    | beaver-server | 好友 + 群组（已合并） |
| message         | 20300    | 30300    | beaver-server | 聊天（chat 更名） |
| ws              | 20400    | -        | beaver-server | WebSocket |
| asset           | 20600    | 30600    | beaver-server | 文件 + 表情（已合并） |
| notification    | 21000    | 31000    | beaver-server | 通知 |
| platform        | 21600    | 31100    | beaver-server | 平台（字典/反馈/埋点/更新已并入） |
| datasync        | 21700    | -        | beaver-server | 数据同步 |
| moment          | 20900    | 30800    | beaver-social | 圈子动态 |
| call            | 21800    | 31800    | beaver-call | 通话 |
| open_api        | 21900    | 30900    | beaver-open | 开放平台 API |
| circle          | 22100    | 32100    | beaver-social | 圈子 |
| ai              | 22200    | 32200    | beaver-ai | AI 服务 |
| agent           | 22300    | -        | beaver-ai | 智能体（Python，静态 X-Internal-Token） |
| doc             | 22400    | 31300    | beaver-online | 在线文档 |

<!-- 其他 -->
| backend_admin   | 40000    | -        | beaver-backend | 后台管理（自带 JWT 鉴权 / CORS） |
