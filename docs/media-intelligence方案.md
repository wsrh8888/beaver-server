# Media Intelligence（媒体智能）整体方案

> 目标：对标大厂 IM + AI（飞书 / 钉钉 / 微信 / ChatGPT Voice 等），补齐「听懂、说出、看懂、会议分析」整层能力，而不是只做一个语音 demo。  
> 原则：**外部云 API 做引擎**；Beaver 自建编排、鉴权、落库、多端同步、Provider 可替换。  
> 状态：方案方向已定；**库表已起草，确认后再写 API / Worker**。  
> 涉及仓库：`beaver-server` / `beaver-agent` / `beaver-desktop` / 移动端。

---

## 一、结论（先看这个）

1. **不要只加一个 `voice` 服务。** 要新增一整层 **Media Intelligence（媒体智能）域**，语音只是其中一块；同一域还要覆盖图片理解、视频理解、会议录制分析等。  
2. **引擎用外部 API**（Azure / Google / OpenAI / 火山 / 阿里等），**客户端禁止直连第三方**；密钥、ACL、缓存、换源、审计都走 Beaver 网关。  
3. **Media = 事实层（听清、看清、转写、合成）**；**Agent = 智能层（对话、推理、纪要文风、工具）**。两者分层，不要把麦克风和 STT 塞进 Agent。  
4. **第一期单进程 `media`（api + 可选 rpc）**，对内按 Speech / Vision / Video / Meeting 分包；量上来再拆微服务，对外 API 前缀保持 `/api/media/v1/...`。  
5. 库：`beaver_media`；端口预留 **22300 / 32300**；表：`media_job` / `media_artifact` / `media_speech_session`（待确认）。

---

## 二、现状与缺口（三仓梳理）

### 2.1 已有能力

| 能力 | 现状 |
|------|------|
| IM 语音消息 | `VoiceMsg` / `AudioFileMsg`，存 `fileUrl` + 时长；桌面主要播放，PC 不录 type=5 |
| 文件上传 | `file_api`（local / 七牛 / MinIO） |
| 实时通话 | `call` + LiveKit（真人通话，非 AI 语音、无转写） |
| AI 对话 | `beaver-agent`：HTTP accepted + MQ → 用户 WS 文本流式（`delta` / `lg` / `done`） |
| 媒体已听 | `markMessageMedia` + datasync |
| Agent 多模态 | `parts` 目前偏 text / image / file；模型表预留 audio，**未落地** |

### 2.2 明确缺口

| 缺口 | 说明 |
|------|------|
| 流式 STT | 手机边说边出字（大厂标配） |
| 批量 STT | 聊天右键「转文字」+ 多端同步 + 可搜索 |
| TTS | AI 用声音回用户，可打断 |
| AI 语音会话 | VAD、endpoint、barge-in、session 恢复 |
| 会议分析 | 录制 → 转写 → 说话人 → 纪要 / ToDo |
| 视频 / 图片理解 | 章节、OCR、描述、给 Agent 引用 |
| 统一任务与 Provider | 可替换供应商、配额、审计 |

---

## 三、为什么不是「只加 voice」

大厂很少把顶层服务叫 Voice。合理切分是：

| 层 | 管什么 | Beaver |
|----|--------|--------|
| **传输** | 实时音视频进房、推流 | 已有 `call` + LiveKit |
| **存储** | 文件 / 录音 / 录像 URL | 已有 `file` |
| **理解与生成** | STT / TTS / OCR / 视频理解 / 会议后处理 | **新建 `media`** |
| **编排** | Agent 推理、工具、多轮对话 | 已有 `agent` |
| **业务承载** | 消息、会话、纪要落库同步 | `chat` / `datasync` / `ws` |

若只建 `voice`：会议分析、视频总结、图片 OCR 会被迫再起一堆碎片服务，或继续往 chat/agent 里塞，边界会烂掉。

**正确做法：一个 Media Intelligence 域，多个能力模块，共用 Job / Provider / Artifact / 权限。**

```
客户端 / 会议 / IM / Agent
            │
            ▼
     ┌──────────────────────────┐
     │  Media Intelligence 域    │  ← 新建
     │  统一：任务、权限、Provider │
     └────────────┬─────────────┘
                  │
     ┌────────────┼────────────┬────────────┐
     ▼            ▼            ▼            ▼
  Speech       Vision       Video        Meeting
  STT/TTS      OCR/描述     抽帧/理解     录制后处理
  流式字幕      图文理解     字幕轨        说话人/纪要
     │            │            │            │
     └────────────┴────────────┴────────────┘
                  │
         外部云 API（主备可切换）
                  │
         结果 → chat / call / agent / 搜索 / datasync
```

---

## 四、外部 API vs 自建：怎么选

| 层 | 谁做 | 为什么 |
|----|------|--------|
| STT / TTS / Vision **模型能力** | 外部 API | 模型、语种、延迟、音色迭代成本极高，非核心竞争力 |
| **业务编排、鉴权、落库、同步、计费、降级** | 自建 `media` | 否则客户端直打第三方，权限/缓存/多端/审计/换源全崩 |

推荐形态：

```
App / Desktop
    │  只调 Beaver
    ▼
gateway → media_api / media realtime
    │
    ├─ 鉴权、幂等、任务队列、会话态
    ├─ Provider Adapter（统一接口）
    │     ├─ OpenAI Realtime / Whisper + TTS
    │     ├─ Azure Speech
    │     ├─ 国内：火山 / 阿里 / 腾讯 等
    │     └─ 备用：失败自动切
    ▼
结果回写 chat / agent / datasync / WS
```

**何时才考虑自建引擎：** 强合规音频不出域、极端定制口音/词库、或日调用量到自建更划算。开源 IM 第一阶段：**100% 外部 API + 自有网关**。

---

## 五、产品能力全景（可以做 / 要对齐大厂的事）

### 5.1 两条语音产品线（不要混一个 API）

| 线 | 用户预期 | 特征 | 实现 |
|----|----------|------|------|
| **A. 聊天语音转文字** | 右键/长按，出字幕，可搜索 | 异步、可缓存、强一致多端 | `media_job` + `media_artifact` |
| **B. 与 AI 语音对话** | 边说边出字，AI 语音回，可打断 | 低延迟、流式、会话态 | `media_speech_session` + Agent |

### 5.2 Speech（语音）

| 能力 | 优先级 | 说明 |
|------|--------|------|
| **流式 STT（边说边出字）** | P0 | 移动端 AI 语音默认路径；partial 灰字、final 定稿 |
| **批量 STT** | P0 | IM 语音消息右键转文字；幂等缓存 |
| **TTS 播报** | P0 | AI 回复语音化；流式优先；可打断（barge-in） |
| 语种检测 / 用户默认语言 | P0 | 自动检测 + 策略 |
| 音色 / 语速 | P1 | `voice_id` 绑定用户或 AI 人设 |
| 降噪 / AEC | P1 | 端侧为主；服务端可做预处理规范 |
| 声纹 | P2 | 可选 |

**移动端 AI 语音 UX（对标大厂）：**

1. 进入语音模式 / 按住说话  
2. 屏幕立刻出 **partial** 字幕（会改）  
3. 静音断句 / 松手 → **final** 定稿  
4. **只用 final**（或用户点发送）进 Agent，partial 不狂打 Agent  
5. Agent 文本流式回 + **并行 TTS 播放**；用户再说话可打断 TTS  

### 5.3 Vision（图）

| 能力 | 优先级 | 说明 |
|------|--------|------|
| OCR | P1 | 截图 / 图片文字 |
| 图片描述 | P1 | 给 Agent 或消息卡片 |
| 与 Agent image part 打通 | P1 | 可先理解再喂 Agent，或 multimodal 直吃 |

### 5.4 Video（视频）

| 能力 | 优先级 | 说明 |
|------|--------|------|
| 音轨转写 | P1 | 本质复用 Speech |
| 关键帧 / 章节 | P1～P2 | 时间戳引用 |
| 视频问答 | P2 | 摘要 + Agent |

### 5.5 Meeting（会议录制分析）

组合产品，不是另起平行服务：

```
LiveKit 录制 → file 存音视频
  → Speech：转写 + 说话人分离
  → Agent / Summary Worker：纪要 / ToDo / 章节
  → artifact 挂 call 资源 → IM 可看、可搜索、多端同步
```

| 能力 | 优先级 |
|------|--------|
| 通话结束后自动/手动生成转写 | P1 |
| 纪要 + 待办 | P1 |
| 实时通话字幕 | P2 |
| 说话人分离 | P1～P2 |
| 企业 DLP / 敏感词 | P2 |

### 5.6 横切（这一域必须有，否则做不成大厂）

- 统一异步任务（job_id、进度、重试、幂等）  
- 统一 Provider 适配与主备降级  
- 统一产物模型（挂 message / file / call / agent_topic）  
- 权限与消息/通话 ACL 一致  
- 转写 / 纪要进搜索索引  
- 配额、审计、失败可重试 UX  
- 会话恢复、弱网重连（realtime）  

### 5.7 与「发语音气泡」的区别

| 场景 | 路径 |
|------|------|
| AI 语音聊 | 实时流式转写 → 文本进 Agent → 可选 TTS |
| IM 发语音气泡 | 仍可录文件发 `VoiceMsg`；需要字幕再 **异步转文字** |
| 两者 | 共用外部 STT，**协议与 API 分开** |

---

## 六、端到端架构

### 6.1 全局挂接

```
call + LiveKit     → 实时通话 / 录制触发
file               → 媒资存储
media              → STT/TTS/OCR/视频理解/会议后处理   ← 新建
agent              → 对话、总结、多轮、工具
chat + ws + datasync → 展示、同步、搜索入口
gateway            → 统一入口与鉴权透传
```

### 6.2 场景 A：右键转文字

```
客户端（messageId + 权限）
  → POST /api/media/v1/speech/transcribe
  → 查是否已有成功 artifact → 有则直接返回
  → 无则入队 media_job → Worker 拉 file → 外部 Batch STT
  → 写 media_artifact（version++）
  → WS 推送 artifact 就绪 + datasync 可拉
  → 气泡下展示转写；可进搜索
```

Agent **不参与**。

### 6.3 场景 B：手机 AI 语音聊（流式）

```
手机麦流 ── media WS/WebRTC ──▶ media（外部 Streaming STT）
                                    │
                          partial / final 字幕回推手机
                                    │
                               final 定稿
                                    ▼
                               agent（当文本 / 可选 audio part）
                                    ▼
                          文本流（现有 MQ→用户 WS）
                                    +
                          media TTS 音频帧 / URL 回推
```

可选方案 B2：LiveKit Bot Participant（与真人通话栈统一）。第一期推荐 **media 自有 realtime session**；LiveKit 继续服务真人通话，后期再评估合并。

### 6.4 文本流与音频流分离

- 文本：继续用现有 `chat_message_stream_receive`（agent delta）  
- 音频 / 实时字幕：走 **media session 通道**（或后续 `audio_delta` 事件）  
- PC 可只显示文字，点播放再拉 TTS `file_key`

### 6.5 Agent 边界

| 该谁做 | 不该谁做 |
|--------|----------|
| Media：听清、看清、转成结构化素材 | Agent 里直接接麦克风、绑死某家 STT |
| Agent：基于素材推理、对话、写纪要 | Agent 变成转码厂 / 媒体存储 |

会议纪要：**转写走 media，总结建议调 agent**（或 media worker 内调 LLM——二选一，推荐前者以统一 Agent 能力）。

---

## 七、服务端改动清单（beaver-server + agent）

### 7.1 新建 `media`

| 项 | 内容 |
|----|------|
| 库 | `beaver_media` |
| 进程 | `media_api`（22300）+ 可选 `media_rpc`（32300） |
| 模块 | Provider Adapter / Job Worker / Realtime Session / ACL / Artifact |
| 部署 | etcd 注册、deploy、gateway 反代 `/api/media/...` |
| MQ | `media_job_topic`；完成通知走 `ws_push_topic` |

第一期 **不拆** speech / video / meeting 三个独立微服务。

### 7.2 改造现有服务

| 服务 | 改动 |
|------|------|
| `gateway_api` | 反代 `/api/media/v1/*`；realtime WS Hijack |
| `common`（WS/MQ） | 新 type：artifact 就绪、会议纪要就绪等；realtime 字幕优先走 media 自有 WS |
| `chat` | 薄：展示转写 / 引用 artifact；**不**内嵌第三方 STT |
| `datasync` | `getSyncMediaArtifacts`（按 version 增量） |
| `call` | 录制结束 / webhook → 投递 `meeting_analyze` job；挂 recording file_key |
| `file` | 基本不动；保证 mime/时长；Opus→16k 等转码可放 media worker |
| `ws_api` | 消费完成事件推在线端 |
| `open_api` | 后期：开放平台也可调转写/TTS |

### 7.3 beaver-agent 薄扩展

| 改动 | 说明 |
|------|------|
| `parts` | 支持 `audio`（url + 可选 transcript）；实时场景也可只收 **纯文本 final** |
| `human_content_from_parts` | 有 transcript 用文本；无则先经 media |
| `SendMessage` 可选字段 | `reply_modality=text\|voice\|both`、`voice_id` |
| 流式 | 文本 delta 不变；TTS 由 media 做或 agent 声明「请 TTS」 |
| 会议总结 | 供 media worker 调用 |
| IM 助手 | 语音消息若已有转写，带文本进图 |

`agent_rpc.proto` 与 beaver-agent 仓保持同步。

### 7.4 明确不改 / 少改

- 不在 `chat_rpc` 写死 STT/TTS  
- 客户端不直连 OpenAI/Azure  
- Provider 配置走 yaml / 环境变量，**第一期不建** provider 配置表  
- auth / friend / group 等业务域基本不动  

---

## 八、客户端改动（桌面 / 移动）

### 8.1 beaver-desktop

| 点 | 改动 |
|----|------|
| 语音消息右键 | `voice.ts` 增加「转文字」→ 调 transcribe |
| 气泡展示 | 转写折叠区；失败重试 |
| AI 侧栏 / 工作台 | 可选麦克风 → realtime session 或本地录完上传 |
| TTS | 播报助手回复；复用现有 AudioPlayer / 独立窗 |
| 流消费 | 如有 `audio_delta` / artifact 事件，接到现有 stream-receiver |

现状：无 MediaRecorder / Speech API；PC 故意不录 type=5——**转文字与 AI 听写仍要做**。

### 8.2 移动端

| 点 | 改动 |
|----|------|
| AI 语音模式 | 实时 session、partial/final UI、打断 |
| 按住说话 / 连续对话 | 产品二选或都支持；协议上 session 统一 |
| 后台音频会话 | iOS/Android 音频焦点、听筒/扬声器 |
| IM 语音气泡 | 可继续发 VoiceMsg；长按转文字 |
| 播放 | 倍速、已听状态（已有 mark） |

---

## 九、音频与媒体规范

| 项 | 建议 |
|----|------|
| 上行（移动） | Opus 优先 |
| 给 ASR | 服务端归一 16k mono WAV/FLAC |
| TTS 输出 | Opus 64k + 备用 AAC；带 duration、checksum |
| 存储 | 仍走 `file_api`；media 只存 metadata / file_key |
| 短语音 vs 长音频 | 同一 STT，不同 max duration / 配额档 |

---

## 十、API 骨架（确认库表后再实现）

### 10.1 P0 Speech

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/media/v1/speech/session/create` | 创建实时会话，返回 sessionId / 连接凭证 |
| WS | `/api/media/v1/speech/session/connect` | 上行音频+控制帧；下行 partial/final/TTS |
| POST | `/api/media/v1/speech/transcribe` | `messageId` / `fileKey` / `force` |
| GET | `/api/media/v1/speech/jobs/:jobId` | 查批量任务 |
| POST | `/api/media/v1/speech/synthesize` | TTS（可另开流式） |

Realtime 控制帧示例：`start` / `audio` / `stop` / `cancel` / `barge_in`。  
下行事件示例：`transcript_partial` / `transcript_final` / `tts_audio` / `error` / `session_closed`。

### 10.2 P1 统一 Job / 会议 / 资源查询

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/media/v1/jobs` | type=meeting_analyze / video_analyze / vision_* |
| GET | `/api/media/v1/jobs/:jobId` | 进度与状态 |
| GET | `/api/media/v1/resources/:resourceType/:resourceId/artifacts` | 某消息/通话上的产物 |

### 10.3 datasync / WS

| 项 | 说明 |
|----|------|
| `POST /api/datasync/v1/getSyncMediaArtifacts` | `since=version` 增量 |
| WS 推送 | artifact 就绪（转写贴到语音消息、纪要贴到通话） |

---

## 十一、数据库设计（待确认）

库名：`beaver_media`。代码：`app/media/media_models/`，已挂入 `database/migrations.go`。

### 11.1 表一览

| 表 | 模型 | 作用 |
|----|------|------|
| `media_job` | `MediaJob` | 异步任务运行态 |
| `media_artifact` | `MediaArtifact` | 客户端可同步产物（真相） |
| `media_speech_session` | `MediaSpeechSession` | 实时语音会话 |

会议 / 视频 / OCR **不另建业务表**：用 `job.type` + `artifact.kind` 扩展。

### 11.2 `media_job`

| 字段 | 说明 |
|------|------|
| `job_id` | 业务唯一 ID |
| `type` | speech_transcribe / speech_synthesize / meeting_analyze / video_analyze / vision_ocr / vision_describe |
| `status` | 1待执行 2执行中 3成功 4失败 5取消 |
| `user_id` | 发起人 |
| `resource_type` + `resource_id` | 挂靠 message / file / call / … |
| `idempotency_key` | `*string`，UNIQUE 允许多 NULL |
| `input_json` | language、file_key、voice_id、force 等 |
| `provider` / `progress` / `error_msg` | 运行态 |
| `started_at` / `finished_at` | 时间 |

### 11.3 `media_artifact`

| 字段 | 说明 |
|------|------|
| `artifact_id` | 业务唯一 ID |
| `job_id` / `session_id` | 来源（可空） |
| `kind` | transcript / transcript_segments / tts_audio / meeting_minutes / meeting_todos / video_chapters / ocr_text / vision_caption |
| `status` | 1就绪 2失败 |
| `user_id` | 同步过滤 |
| `resource_type` + `resource_id` + `kind` | **联合唯一**：同一资源同一类产物一条当前记录 |
| `content_text` / `payload_json` / `file_key` | 文本 / 结构化 / 文件引用 |
| `language` / `provider` | 元数据 |
| `version` | **datasync 游标**，覆盖更新时自增 |

### 11.4 `media_speech_session`

| 字段 | 说明 |
|------|------|
| `session_id` | 业务唯一 ID |
| `user_id` / `device_id` | 用户与设备 |
| `scene` | ai_voice / dictation / live_caption |
| `status` | 1已创建 2进行中 3已关闭 |
| `ref_type` + `ref_id` | 关联 conversation / agent_topic / call（可空） |
| `language` / `voice_id` / `provider` | 识别与合成 |
| `meta_json` | reply_modality、model_id、打断策略等 |
| `closed_at` | 关闭时间 |

### 11.5 待拍板原则

1. 三张表是否够（会议不独立成表）？  
2. artifact `(resource_type, resource_id, kind)` 唯一 + 覆盖更新 `version++`？  
3. 实时 **partial 不落库**，仅 final / 结束后写 artifact？  
4. 是否需要 `org_id`、计费字段、独立说话人表（建议 P1 说话人进 `payload_json`，不够再拆表）？  

**确认前不写 API / Worker 逻辑。**

---

## 十二、落地顺序

| 阶段 | 内容 |
|------|------|
| **0** | 库表确认（当前） |
| **1** | `media_api` 骨架 + gateway + Provider 空实现 |
| **2** | Batch 转写（右键）+ artifact + datasync/ws |
| **3** | Realtime session：流式 STT + TTS + 打断 |
| **4** | Agent 薄对接（final 文本 / reply_modality） |
| **5** | Call 录制 → meeting job → 转写 + 纪要 |
| **6** | Vision / Video job；搜索索引；配额审计 |
| **7** | 量起来再按需拆子服务 |

---

## 十三、设计原则备忘

- STT/TTS **不要**写进 `chat_rpc` / `agent_rpc` 内部  
- 转写是 **强缓存任务**；AI 语音是 **软实时会话**——分开扩容  
- 文本流与音频流分离  
- 幂等：`idempotency_key`；产物 `version++`；语音轮次可用 `turn_id`（可放 session meta）  
- 大文件始终 `file`；media 只存引用与文本事实  

---

## 十四、相关路径索引

| 项 | 路径 |
|----|------|
| 本方案 | `docs/media-intelligence方案.md` |
| 模型 | `app/media/media_models/` |
| 迁移注册 | `database/migrations.go`（库 `beaver_media`） |
| 端口 | `server.md`（22300 / 32300） |
| Agent 契约 | `app/agent/agent_rpc/agent_rpc.proto`（权威在 beaver-agent） |
| 通话 | `app/call/` + LiveKit |
| 文件 | `app/file/` |
