# Seedance 视频任务回调（callback_url）设计计划

> 本文只写计划，不含代码实现、不含时间估算。状态：**已评审，实施中**（分支 `feat/seedance-callback`）。

- 基线：`jjdufu/kt-api` `main` @ `0863808`（`0.0.2`，fork of new-api）。前序计划：`docs/seedance-assets-plan.md`。
- 供应商文档：`/workspace/uploads/apidoc.txt`（`apidoc.pdf` 的文本版，“整理日期：2026-09-18”，服务地址 `https://yisu-api.gigaai.cc`）。
- 方舟官方文档（本地缓存）：`/workspace/create.md`、`/workspace/d1520757.md`（创建视频生成任务）、`/workspace/d1521309.md`（查询视频生成任务）；在线版 `https://ark.volcengine.com/region:cn-beijing/docs/create-video-generation-task-api`（2026-10-19 复核内容一致）。
- 目标版本：下一个 release-please 版本（预计 `0.0.3`，以 release-please 实际计算为准）。

---

## 1. 文档原文与来源标注

### 1.1 供应商文档原文（逐字）

`/workspace/uploads/apidoc.txt` 第 202 行（请求字段表）：

```
 callback_url               string    否      任务状态变更后的回调地址。
```

`/workspace/uploads/apidoc.txt` 第 427–444 行（字段详解）：

```
callback_url
可选字符串。任务状态改变时，平台向该地址发送 POST 请求，回调结构与查询任务接口的响应体一致。可能出现的
status：

 状态                               说明

 queued                           排队中

 running                          运行中

 succeeded                        已成功

 failed                           已失败

 expired                          超过任务过期阈值


成功或失败回调在发送未获确认时最多重试三次，每次等待约 5 秒。
```

与回调直接相关的上下文（同文件）：

- 第 479–481 行：`execution_expires_after 可选整数，默认 172800 秒（48 小时），范围为 3600 至 259200。时间从任务的 created_at 开始计算；超过阈值后，排队或运行中的任务会转为 expired。`
- 第 796–806 行（查询接口“任务状态”表）：`queued / running / cancelled（任务已取消；仅排队中的任务可被取消。）/ succeeded / failed / expired`。
- 第 856–873 行：`任务已失败时，查询返回 HTTP 400，并在同一裸响应中给出错误`，示例 `{"error":{"code":"task_not_exist","message":"task_not_exist","type":"yisu_api_error","request_id":"..."}}`。
- 第 876 行：`建议间隔数秒轮询，直到 status 为 succeeded、failed、cancelled 或 expired。`
- 第 707 行：查询接口 `该接口会同步视频任务的最新状态。仅能查询当前认证用户创建的任务。`
- 文档全文**没有**取消任务接口（14 个接口中无 DELETE / Cancel）。

### 1.2 方舟官方文档原文（逐字）

`/workspace/create.md` 第 647–659 行（与 `d1520757.md` 第 1036 行起、在线版一致）：

```
**callback_url** `string` | 回调地址

填写本次生成任务结果的回调通知地址。当视频生成任务有状态变化时，方舟将向此地址推送 POST 请求，请求内容结构与 [查询任务 API](/docs/get-video-generation-task-api) 的返回体一致。

> **说明**
>
> 回调返回的 `status` 包括以下状态：
>
> * `queued`：排队中
> * `running`：任务运行中
> * `succeeded`：任务成功（如发送失败，即 5 秒内未接收到成功发送的信息，会回调三次）
> * `failed`：任务失败（如发送失败，即 5 秒内未接收到成功发送的信息，会回调三次）
> * `expired`：任务超时，即任务处于「运行中或排队中」状态超过过期时间，可通过 `execution_expires_after` 字段设置过期时间
```

`/workspace/d1521309.md`（查询接口）：失败时返回体含 `error` 对象（`任务成功时返回 null，任务失败时返回错误数据`，含 `error.code`、`error.message`）；`cancelled` 仅 `queued` 可取消，24h 自动删除。

### 1.3 逐项来源矩阵

| 项目 | 供应商文档 | 方舟文档 | 本计划的取值 | 来源 |
|---|---|---|---|---|
| 请求字段 | `callback_url` string 可选 | 同 | 同 | 供应商 |
| 触发时机 | “任务状态改变时” | “有状态变化时” | 每次**对外可见** `status` 变化触发一次 | 供应商 |
| 可能的回调 status | queued/running/succeeded/failed/expired | 同 | 同；`cancelled` 默认**不**回调（见 Q3） | 供应商 |
| HTTP 方法 | POST | POST | POST | 供应商 |
| 负载结构 | “与查询任务接口的响应体一致” | 同 | 与 kt-api 自身 GET 响应体**逐字节一致** | 供应商 + 我方（字节级） |
| Content-Type | 未写 | 未写 | `application/json` | 我方 |
| 请求头 / 签名 | 未写（无签名） | 未写（无签名） | 默认不签名；可选 HMAC 扩展（默认关） | 我方 |
| 重试范围 | 仅“成功或失败回调” | 仅 succeeded/failed | 仅 succeeded/failed 重试；queued/running/expired 单次 | 供应商 |
| 重试次数 | “最多重试三次” | “会回调三次” | 首次 + 最多 3 次重试（共 ≤4 次），见 Q1 | 供应商（与方舟措辞有差异） |
| 重试间隔 | “每次等待约 5 秒” | 未写间隔 | 5 秒（±0.5s 抖动） | 供应商 |
| 单次超时 | 未写 | “5 秒内未接收到成功发送的信息” | 5 秒 | 方舟 |
| “确认”的定义 | 未写 | “成功发送的信息” | 5 秒内收到 HTTP 2xx（不看响应体） | 我方 |
| 重定向 | 未写 | 未写 | 不跟随，3xx 视为失败 | 我方（SSRF） |
| 顺序保证 | 未写 | 未写 | 同一任务按序、终态不被旧状态覆盖；跨任务无序 | 我方 |
| 失败任务的负载 | 查询返回 HTTP 400 + 裸 `error`（无 status） | 200 + `status:"failed"` + `error{code,message}` | 见 Q2，默认采用方舟形状 | 方舟 + 我方 |

---

## 2. 现状分析（代码事实）

1. **剥离点**：`plugins/tasks/doubao/plugin.js` 第 844–846 行在 `buildSubmitRequest` 中 `delete body.callback_url`；但 `native.createTask`（约 748 行）把整个请求体放进 `requestBody.metadata`，需确认 `callback_url` 是否随 metadata 落入 `tasks.properties` / 日志，实现时要在持久化前剔除。
2. **别名路由**：`middleware/seedance_video_alias.go` 把 `/api/v3/contents/generations/tasks[/:task_id]` 映射到 Doubao 插件声明的 `/doubao/api/v3/...`（`plugin.js` 第 318–324 行 `taskCreated` / `taskStatus` 渲染器）。
3. **GET 响应来源**：`middleware/task_plugin.go` `renderTaskPluginQuery`（1159 行起）**只读数据库**，不实时请求上游；`service/task_plugin_view.go` `BuildTaskPluginView` 用 `replacePrivateTaskID` 把 `id/task_id/taskId` 中的上游 ID 换成公开 ID；插件 `taskStatus` 直接返回 `task.data`（上游最近一次查询体）并覆盖 `id`，最后 `c.JSON(200, result)`。
4. **状态推进**：`service/task_polling.go` 每 15 秒由 system task `async_task_poll`（带 per-type 租约，`service/system_task.go`）跑一轮；`updateVideoSingleTask` 写 `task.Data = redactVideoResponseBody(body)`，终态用 `UpdateWithStatus` CAS，CAS 赢家调用 `finalizeTerminalTask` 结算/退款。
5. **内部状态丢信息**：插件 `parseTaskResult` 把 `expired/cancelled` 都映射为 `FAILURE`；因此“对外 status”只能从 `task.data.status`（渲染结果）得到，不能用 `tasks.status`。
6. **提交后的 data 没有 status**：提交成功时 `task.data` 是创建响应 `{"id": ...}`，在第一次轮询前 GET 返回体没有 `status` 字段。
7. **本地判定的终态与 GET 不一致（前置缺陷）**：`sweepTimedOutTasks`、`failTaskFromPoll`（上游 404、连续失败放弃）只改 `tasks.status/fail_reason`，**不改 `task.data`**，于是 GET 仍显示旧的 `running`。回调若要与 GET 一致，必须先修。
8. **供应商失败任务查询返回 HTTP 400（前置缺陷）**：`classifyPollHTTP(400)=other_client` → 插件 `parseTaskResult` 拿到无 `status` 的 `{"error":...}` 返回 `UNKNOWN` → `recordPollFailure(unrecognized)`，直到 `TASK_POLL_MAX_FAILURES` 才以固定文案失败，真实失败原因丢失、失败识别延迟多个轮询周期。
9. **超时冲突**：db51 `TASK_TIMEOUT_MINUTES=1440`（24h）小于文档默认 `execution_expires_after=172800`（48h，最大 72h），kt-api 会先于上游判失败。
10. **SSRF 基建已有**：`common/ssrf_protection.go`（私网/回环/链路本地/CGNAT 100.64/10（含阿里云元数据 100.100.100.200）/ULA fc00::/7（含 AWS fd00:ec2::254）/IPv4-mapped 等）+ `service/protected_fetch_client.go`（**拨号时**逐个校验解析结果再按 IP 拨号，天然防 DNS rebinding）。但它绑定全局 `fetch_setting`，管理员可能为其他用途放宽，回调需要独立配置。
11. **监控栈**：db51 compose 项目 `ktapi`；Prometheus 业务指标来自节点上的 `sql_exporter`（`node-agent/sql_exporter/ktapi_realtime.collector.yml`，`job=ktapi_sql`，:9399），只读账号 `ktapi_monitor` 按列授权（`node-agent/create-monitor-role.sql`）；告警规则在 `monitor/prometheus/rules/*.rules.yml`；日志经 Alloy 进 Loki（ktapi 日志格式 `[LEVEL] 时间 | request_id | message`，仅 `level` 做标签）。kt-api 进程本身**没有** `/metrics` 端点。
12. **db51 可达性**：db51 是内网 `192.168.50.51`，域名 `kt.ggok.win` 使用 `tls internal`，供应商无法回调到 db51。生产节点走 `CADDY_TLS=acme` 公网证书，可达。

---

## 3. 关键设计决策

### 3.1 对客户的契约

- **请求**：`POST /api/v3/contents/generations/tasks` 接受可选 `callback_url`（字符串）。不传时行为与现在完全一致。
- **回调请求**：`POST <callback_url>`，`Content-Type: application/json`，`User-Agent: kt-api-callback/<version>`。
- **回调 body = 该次状态变化时刻 kt-api 自己 `GET /api/v3/contents/generations/tasks/{id}` 会返回的字节**（同一个渲染函数、同一份 `task.data`、同一个公开 ID 替换、同一 JSON 编码器）。实现上抽出一个 `RenderSeedanceTaskStatus(task) ([]byte, error)`，GET 路由与回调共用；回调在**入队时**生成快照存入 outbox，重试时发送同一快照（不重新渲染），保证每个事件都是“那次变化”的状态。
- **ID**：`id` 为 kt-api 公开任务 ID（`task_xxx`）。不出现上游任务 ID、渠道 base URL、渠道 Key、`private_data` 中任何字段（沿用 `replacePrivateTaskID` 并补测试：任何层级的 `id/task_id/taskId` 等于上游 ID 都被替换；body 中不得包含上游 ID 子串，否则单测失败）。
- **`content.video_url` / `last_frame_url`**：与 GET 一致（目前为上游签名直链）。是否改为 kt-api 代理地址属于 GET 契约变化，不在本计划内（见 Q6）。
- **触发**：对外 `status` 从 A 变为 B（B ∈ {queued, running, succeeded, failed, expired}）时触发一次；同一 status 只触发一次；`status` 字段缺失（提交后首次轮询前）不触发。
- **不补造中间态**：若一次观测直接从 `queued` 变为 `succeeded`，不伪造 `running` 事件（伪造的 body 无法与当时的 GET 一致）。上游回调快路径开启时会显著降低漏态概率。
- **默认附加头（kt-api 扩展，不改变 body）**：`X-KtApi-Event-Id`（全局唯一，同一事件重试不变）、`X-KtApi-Delivery-Attempt`（从 1 开始）、`X-KtApi-Task-Id`。供应商/方舟未规定任何头，额外头不影响兼容；可在后台关闭（见 Q5）。
- **客户确认**：5 秒内返回任意 2xx 即视为确认；响应体忽略（最多读取 64 KiB 后丢弃）。

### 3.2 kt-api 如何获知状态：快路径 + 兜底

| 方案 | 优点 | 缺点 |
|---|---|---|
| A：向上游传 kt-api 自己的回调地址 | 秒级感知，中间态更完整 | 需要公网可达（db51 不满足）；上游回调无签名，body 不可信；供应商仅对终态重试 3 次，丢了就没了 |
| B：沿用 15 秒轮询 | 已有、可靠、带 CAS/计费不变量、多实例租约 | 延迟最多一个周期 + 排队；可能漏掉短暂的中间态 |

**决策：A 作为可选快路径 + B 作为永远开启的兜底。**

- 轮询仍是唯一的“状态写入者”。上游回调**只作为触发信号**：收到合法回调后，对该任务立即执行一次与轮询完全相同的单任务刷新（`FetchTask` → `ParseTaskResult` → `UpdateWithStatus` CAS → 结算），**不信任回调 body**。理由：
  1. 文档无签名，body 可伪造；
  2. 供应商可能把 `callback_url` 原样交给方舟，方舟回调的 `id` 可能是方舟 ID（`cgt-...`）而非供应商 ID，body 形状不可依赖（见 Q4）；
  3. 计费只在 CAS 赢家上发生一次，复用同一路径即可保持 `.agents/rules/billing.md` 的不变量。
- 开关与前提：
  - 全局选项 `TaskCallbackUpstreamPublicBase`（默认空 = 关闭快路径）。db51 保持空，仅靠轮询。
  - 渠道级开关（channel setting）`upstream_callback_enabled`，默认关；只有开启的渠道在 `buildSubmitRequest` 时被注入 `callback_url`。
  - 注入策略选项：`only_when_customer_callback`（默认）/ `always` / `never`。
- 上游回调端点：`POST /api/task-callback/upstream/:token`（原生 Go 路由，不在 `/api/v3` 下，不经 TokenAuth）。
  - `token`：每个任务独立 32 字节随机数（base64url），只在库中存 SHA-256；有效期到任务终态后 1 小时。
  - 处理：限流（按来源 IP + 全局）→ body 上限 1 MiB → 查 token 哈希 → 找不到返回 404 → 找到则把任务 ID 投入“立即刷新”队列（per-task singleflight，同一任务最小间隔 2 秒）→ **立即**返回 `200 {}`（必须在供应商的 5 秒窗口内）。
  - 可选的弱校验：若 body 能解析出 `id`/`status`，仅用于日志与指标（`hint_status`），不参与状态判定。
- 多实例：刷新动作本身依赖 `UpdateWithStatus` CAS，多个实例同时刷新同一任务也只有一个赢家结算；回调入队见 3.4 的唯一键去重。

### 3.3 客户 `callback_url` 校验

提交时（在预扣费之前）同步校验，失败返回 HTTP 400，错误信封与现有别名路由一致：`{"error":{"code":"InvalidParameter","message":"callback_url ...","param":"callback_url","type":"invalid_request_error"}}`（错误码取值见 Q7）。

1. 类型为字符串；去首尾空白后为空视为未传。
2. 长度 ≤ 2048 字节（可配置 `TaskCallbackMaxURLLength`）。
3. `url.Parse` 成功；scheme 仅 `http` / `https`（可配置“仅允许 https”，默认关以兼容文档）；必须有 host；禁止 userinfo（`user:pass@`）；禁止 fragment 以外的控制字符。
4. 端口：默认 80、443、8080、8443（独立配置，不复用 `fetch_setting.AllowedPorts`）。
5. host 为 IP 字面量（含 `0x7f.1`、十进制整数、IPv4-mapped IPv6 等变体，统一用 `net.ParseIP` / `netip` 规范化后判断）：落入私网 / 回环 / 链路本地 / CGNAT / ULA / 组播 / 保留 / 元数据地址 → 拒绝。
6. host 为域名：提交时做一次解析（超时 2 秒）；**全部**解析结果都合法才接受；解析失败按配置决定拒绝（默认拒绝）或放行等投递时再判（测试用）。
7. **投递时再校验（权威）**：使用 `protectedFetchDialer` 同款拨号器，在拨号前对每个解析 IP 校验后按 IP 直连，杜绝 DNS rebinding；不使用环境代理（避免绕过拨号检查）；不跟随重定向；TLS 正常校验证书。
8. 独立配置项（后台“任务回调”分组）：
   - `TaskCallbackEnabled`（总开关，默认开）
   - `TaskCallbackSSRFAllowPrivate`（默认关；仅测试环境，如 db51 用内网 mock 接收端）
   - `TaskCallbackDomainAllowlistMode` + `TaskCallbackDomainList`（默认关；打开后只允许列表内域名，支持 `*.example.com`）
   - `TaskCallbackAllowedPorts`
   - `TaskCallbackHTTPSOnly`
9. 校验通过后：从发往上游的请求和持久化的 metadata 中剔除客户 `callback_url`；客户 URL 只存到订阅表；日志中只打印 scheme + host + URL 的 SHA-256 前 12 位，不打印路径与查询串（可能含客户密钥）。

### 3.4 投递：持久化 outbox

#### 3.4.1 表结构（GORM `AutoMigrate`，SQLite / MySQL ≥5.7.8 / PostgreSQL ≥9.6 通用）

`task_callback_subscriptions`（每个带回调的任务一行）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | bigint PK 自增 | |
| `task_id` | varchar(191) | 公开任务 ID，**唯一索引** |
| `user_id` | int | 索引 |
| `callback_url` | text | 客户 URL（需要时可用 `CRYPTO_SECRET` 加密存储，见 Q8） |
| `callback_url_hash` | char(64) | SHA-256，日志/后台展示用 |
| `upstream_token_hash` | varchar(64) 可空 | 上游回调 token 的 SHA-256，**唯一索引**（空值不参与） |
| `upstream_token_expires_at` | bigint | Unix 秒 |
| `last_status` | varchar(20) | 最近一次入队的对外 status |
| `last_seq` | int | 该任务事件序号 |
| `created_at` / `updated_at` | bigint | Unix 秒 |

`task_callback_events`（outbox，每个状态变化一行）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | bigint PK 自增 | |
| `event_id` | varchar(64) | 全局唯一（UUIDv7 或 `evt_` + 随机串），**唯一索引**，即 `X-KtApi-Event-Id` |
| `subscription_id` | bigint | 索引 |
| `task_id` | varchar(191) | 公开任务 ID |
| `user_id` | int | |
| `seq` | int | 同任务内递增 |
| `status` | varchar(20) | 对外 status |
| `payload` | text（MySQL 为 `mediumtext`，通过 GORM `type` 按方言或统一用 `text` + 长度上限校验） | 快照字节 |
| `payload_sha256` | char(64) | |
| `state` | varchar(16) | 见状态机 |
| `attempts` | int | 已尝试次数 |
| `max_attempts` | int | 入队时按 status 计算（终态 4，其余 1） |
| `next_attempt_at` | bigint | |
| `lease_owner` | varchar(64) | 节点 ID + 进程随机串 |
| `lease_until` | bigint | |
| `last_http_status` | int | |
| `last_error` | varchar(512) | 截断、脱敏 |
| `last_latency_ms` | int | |
| `delivered_at` | bigint | |
| `created_at` / `updated_at` | bigint | |

索引：
- 唯一 `(task_id, status)` —— 同一任务同一 status 只入队一次（幂等核心）。
- `(state, next_attempt_at)` —— 调度扫描。
- `(state, lease_until)` —— 过期租约回收。
- `(user_id, created_at)` —— 后台/用户查询。
- `(task_id, seq)`。

迁移约束（遵循 `AGENTS.md`）：
- 不用 JSON 列类型，payload 用 TEXT；不用布尔默认值；索引字符串列 ≤191；时间统一 Unix 秒 bigint。
- 只新增表，不改 `tasks` 表结构；新增 `model/main.go` AutoMigrate 注册。
- 必须在真实 SQLite、MySQL、PostgreSQL 实例上验证：首次迁移、重复启动不产生 ALTER、唯一键冲突行为、租约 CAS 的 `RowsAffected` 语义。
- 保留期清理：已终结事件默认保留 7 天（可配置），由 system task 每小时分批删除（每批 ≤1000 行）。

#### 3.4.2 事件状态机

```
           入队
            │
            ▼
        pending ──(被租约领取)──► delivering ──2xx──► succeeded
          ▲  │                       │
          │  │                       ├─失败且 attempts < max ──► pending（next_attempt_at = now + 5s）
          │  │                       └─失败且 attempts = max ──► failed（终结）
          │  └─(同任务更新的非终态事件到达，且本事件为非终态)──► superseded（终结，不再发送）
          └─(租约过期回收)─ delivering
  另外：订阅被取消 / 总开关关闭 / 用户被删除 ──► cancelled（终结）
```

- `superseded` 只作用于 `queued` / `running` 事件；`succeeded` / `failed` / `expired` 永不被跳过。
- `failed` 后后台可手动“重新投递”：新建一次投递尝试（attempts 清零、保留同一 `event_id`，客户可去重）。

#### 3.4.3 入队（只在一个地方发生）

- 新增 `service.EnqueueTaskCallbackIfChanged(tx, task)`：渲染对外 body → 取 `status` → 与订阅的 `last_status` 比较 → 不同且属于回调集合 → 在**同一事务**中：插入事件（`ON CONFLICT DO NOTHING`，GORM `clause.OnConflict{DoNothing: true}` 三库通用）+ 更新订阅 `last_status/last_seq`。
- 新增 `model.Task.UpdateWithStatusTx(tx, from)`，把“任务 CAS 更新 + 入队”放在一个事务里；只有 CAS `RowsAffected=1` 才入队。
- 所有写任务状态/数据的路径统一走它：`updateVideoSingleTask`、`updateBatchTasks`、`sweepTimedOutTasks`、`failTaskFromPoll`、`recordPollFailure`（放弃时）、上游回调触发的单任务刷新、提交成功（若创建响应已带 status）。
- 非终态的 `data` 变化（如 running 期间字段变化）只更新任务，不入队（status 未变）。
- 对账兜底：调度器每 60 秒扫描“有订阅、`last_status` 不等于当前渲染 status、且任务 `updated_at` 超过 30 秒”的任务补入队，覆盖事务外异常与历史数据。

#### 3.4.4 调度与投递

- 每个实例运行一个 dispatcher 协程（受 `TaskCallbackEnabled` 控制）：
  - 唤醒：本实例入队后通过 channel 立即唤醒；否则每 1 秒轮询一次。
  - 领取：`SELECT id FROM task_callback_events WHERE state='pending' AND next_attempt_at<=now ORDER BY next_attempt_at, id LIMIT n`，然后逐行 `UPDATE ... SET state='delivering', lease_owner=?, lease_until=now+30, attempts=attempts+1 WHERE id=? AND state='pending'`，`RowsAffected=1` 才算领取成功（三库通用的乐观 CAS，SQLite 也可用；不依赖 `SKIP LOCKED`）。
  - 同任务顺序：领取前检查同 `task_id` 是否存在 `seq` 更小且未终结的事件；存在则：若那个事件是非终态 → 把它标为 `superseded`；若是终态（理论上不会出现更大 seq）→ 等待。即同一任务同一时刻最多一个在途投递。
  - 并发：全局 worker 数 `TaskCallbackWorkers`（默认 16）；每个目标 host 并发上限 `TaskCallbackPerHostConcurrency`（默认 4），避免单个慢客户拖垮全部。
  - 单次超时 5 秒（连接 + TLS + 发送 + 等待响应头总计，`context.WithTimeout`），超时视为“未获确认”。
  - 结果写回：`UPDATE ... WHERE id=? AND lease_owner=? AND state='delivering'`，防止租约丢失后误写。
- 重试（严格按文档）：
  - `succeeded` / `failed`：首次失败后最多再试 3 次，每次等待 5 秒（±0.5 秒抖动）。总计 ≤4 次，约 20 秒内结束。
  - `queued` / `running` / `expired`：只发一次，失败即 `failed`。
  - kt-api 扩展（默认关，`TaskCallbackExtendedRetry`）：耗尽文档重试后再按指数退避（30s、2m、10m、30m、1h，最多 24h）继续尝试，适用于终态事件。开启后客户可能在文档规定之外收到回调，所以默认关闭（见 Q1）。
- 失败分类（写入 `last_error` 与指标标签）：`timeout`、`dns`、`ssrf_blocked`、`connect`、`tls`、`http_3xx`、`http_4xx`、`http_5xx`、`body_too_large`。`ssrf_blocked` 不重试（结果不会改变）。
- 多实例安全：领取与写回均为条件 UPDATE；租约 30 秒（大于单次超时 5 秒的数倍）；进程退出时尽力释放本实例租约。
- 重启恢复：启动时与每分钟回收 `state='delivering' AND lease_until<now` 的行 → 回到 `pending`（已计入的 attempts 保留，若已达上限则判 `failed`）。由于超时只有 5 秒、租约 30 秒，重启最多让一个事件被重复发送一次——这是“至少一次”语义，客户用 `X-KtApi-Event-Id` 或 `id+status` 去重（文档化）。

### 3.5 签名（kt-api 扩展，默认关）

- 供应商与方舟文档均无签名，默认不签名以保持兼容。
- 可选 `TaskCallbackSigningEnabled`：开启后增加
  - `X-KtApi-Timestamp: <unix 秒>`
  - `X-KtApi-Signature: v1=<hex(HMAC-SHA256(secret, timestamp + "." + event_id + "." + body))>`
- 密钥：按用户生成（用户在个人设置中查看/轮换），存储时用 `CRYPTO_SECRET` 加密；轮换期间同时发送新旧两个签名（`v1=new,v1=old`）一段时间。
- 文档给出客户侧校验示例（时间窗 5 分钟、常数时间比较）。

### 3.6 计费、日志、后台

- 回调免费：不预扣、不结算、不写消费日志；不改变任务计费路径。上游回调触发的刷新复用轮询的 CAS + `finalizeTerminalTask`，不会重复结算。
- 投递日志：每次尝试写一行英文结构化日志（`common.LogText`，并补 `common/log_text.zh-CN.json` 中文），字段见 3.7；不写入 `logs` 表（避免污染消费统计与 `ErrorLogsSpike` 告警）。
- 后台可见性：
  - 管理员：任务详情页新增“回调”页签（订阅 URL 的 host + 哈希、事件列表、每次尝试的状态码/耗时/错误、手动重投按钮）；新增“任务回调”列表页（按状态/用户/时间筛选）。
  - 普通用户：在自己的任务详情中只读查看事件与结果（不显示其他用户数据）。
  - 新增 API：`GET /api/task/:id/callbacks`（用户）、`GET /api/admin/task-callbacks`、`POST /api/admin/task-callbacks/:event_id/redeliver`（管理员，记审计日志）。
- 配置：所有 `TaskCallback*` 选项通过 `config.GlobalConfig.Register("task_callback_setting", ...)` 注册，后台“系统设置 → 运维设置”新增“任务回调”分组；部署文档列出对应选项与 db51 / 生产推荐值。前端只在 GitHub Actions 构建（box 内存不足）。

### 3.7 监控（适配现有 `ktapi` 监控栈）

kt-api 没有 `/metrics`，保持现有做法：**用 `sql_exporter` 从 outbox 表出指标**，不新开端口。

- `node-agent/sql_exporter/ktapi_realtime.collector.yml` 新增查询：
  - `ktapi_task_callback_events_pending{status}`：`state IN ('pending','delivering')` 数量。
  - `ktapi_task_callback_oldest_pending_seconds`：最老待投递事件的等待秒数。
  - `ktapi_task_callback_finished_15m{status,state}`：最近 15 分钟终结事件数（succeeded/failed/superseded）。
  - `ktapi_task_callback_attempt_latency_avg_15m`：平均 `last_latency_ms`。
  - `ktapi_task_callback_terminal_undelivered_15m`：终态事件最终 `failed` 数。
- `node-agent/create-monitor-role.sql` 按列授权：`GRANT SELECT (id, task_id, status, state, attempts, next_attempt_at, lease_until, last_http_status, last_latency_ms, created_at, updated_at, delivered_at) ON task_callback_events TO ktapi_monitor;`（不授 `payload`、`last_error`、订阅表的 `callback_url`）。
- `monitor/prometheus/rules/business.rules.yml` 新增：
  - `TaskCallbackBacklog`：`ktapi_task_callback_oldest_pending_seconds > 120`，`for: 5m`，warning（调度器卡住/全部 worker 被慢客户占满）。
  - `TaskCallbackFailureRateHigh`：最近 15 分钟终态事件失败率 > 50% 且样本 ≥ 10，`for: 10m`，warning（多半是客户接收端问题，提示看 Loki）。
  - `TaskCallbackLeaseStuck`：`delivering` 且 `lease_until` 已过期的行数 > 0 持续 5 分钟，critical（回收逻辑失效）。
- Loki：日志行保持 ktapi 现有格式，message 中用固定 `key=value`：`task_callback event=delivery_attempt event_id=... task_id=... status=succeeded attempt=2 result=http_5xx http_status=502 latency_ms=812 host=example.com url_hash=ab12cd34ef56`；上游回调：`task_callback event=upstream_hint task_id=... hint_status=running accepted=true`。不新增 Loki 标签（保持低基数），Grafana 用 `|~ "task_callback"` + `| logfmt` 查询。Grafana 日志看板新增一个“任务回调”面板。

---

## 4. 前置修正（回调正确性依赖）

1. **本地终态写回 data**：kt-api 自行判定终态时（超时扫描、连续失败放弃、上游 404），把 `task.data` 改写为方舟形状的失败体：保留上一份 data 的 `model/created_at` 等字段，设 `status`（超时 → 见下）、`error:{code,message}`、`updated_at=now`。保证 GET 与回调一致。
2. **超时语义**：
   - 默认 kt-api 超时扫描改为每任务截止时间 = `created_at + execution_expires_after（未传则 172800）+ 宽限（默认 30 分钟）`，并以 `TASK_TIMEOUT_MINUTES` 作为全局上限兜底；部署文档把 db51/生产的 `TASK_TIMEOUT_MINUTES` 调到 ≥ 4350（72h + 30min）。
   - 超过截止时间仍未收到上游终态时，对外 status 用 `expired`（与文档“超过任务过期阈值”语义一致）还是 `failed` 见 Q9，默认 `expired`。
3. **供应商失败查询（HTTP 400 + 裸 error）**：插件 `parseTaskResult` 在 `resp.status==400` 且 body 只有 `error` 对象时判为 `FAILURE`（reason 取 `error.message`），并把 `task.data` 合成为方舟形状：`{id, model（取上一份 data）, status:"failed", error:{code, message}, created_at, updated_at}`；轮询层对“400 + 可识别的失败体”不再计入 unrecognized。需在真实供应商上确认其他 400 场景不会被误判（见 Q2）。
4. **`cancelled`**：供应商无取消接口，但查询可能返回 `cancelled`。默认按文档回调列表**不发送** `cancelled` 回调（GET 仍如实返回），见 Q3。

---

## 5. 测试计划

### 5.1 单元测试

- `callback_url` 校验表驱动：合法 http/https；非法 scheme（`ftp`、`file`、`gopher`、`javascript`）；userinfo；超长；端口；IP 字面量全集（`127.0.0.1`、`0.0.0.0`、`10.x`、`172.16.x`、`192.168.x`、`100.64.x`、`100.100.100.200`、`169.254.169.254`、`[::1]`、`[fe80::1]`、`[fd00:ec2::254]`、`[::ffff:127.0.0.1]`、`0x7f000001`、`2130706433`、`127.1`）；域名解析到私网；allowlist 模式；测试放行开关。
- 拨号器：stub resolver 首次返回公网、拨号时返回私网（rebinding）→ 拒绝；多个 A 记录含一个私网 → 整体拒绝。
- 渲染一致性：同一任务，GET 路由响应字节 == 回调快照字节；上游 ID 不出现在 body 任何位置。
- 入队：各 status 变化入队一次；重复同 status 不入队；唯一键冲突静默；CAS 失败不入队；`superseded` 规则；状态缺失不入队。
- 重试调度：终态 4 次、间隔 5 秒；非终态 1 次；`ssrf_blocked` 不重试；扩展重试开关。
- 签名：固定向量；开关关闭时无签名头。
- 上游回调端点：token 正确/错误/过期；body 超限；返回在 5 秒内；singleflight 合并。

### 5.2 httptest 端到端（集中在一个测试文件）

- 组件：mock 供应商（可编程状态序列：queued → running → succeeded / failed（HTTP 400 裸 error）/ expired / cancelled）、mock 客户接收端（记录 body、头、次数；可编程返回 2xx/5xx/超时/3xx）、kt-api 路由 + SQLite。
- 覆盖：
  1. 每一种状态转换都恰好收到一次回调，body 与同时刻 GET 字节一致，`id` 为公开 ID。
  2. 发往 mock 供应商的创建请求中不含客户 `callback_url`；快路径开启时含 kt-api 自己的 `/api/task-callback/upstream/<token>`。
  3. 快路径：mock 供应商回调 kt-api → kt-api 立即刷新 → 客户秒级收到；伪造 body（错误 status）不影响结果。
  4. 客户接收端前 3 次 500、第 4 次 200 → 终态共 4 次尝试；连续 4 次失败 → `failed`；非终态失败只 1 次。
  5. 客户端超时（6 秒才响应）→ 视为失败并重试。
  6. 3xx 不跟随。
  7. SSRF：提交时拒绝的 URL 返回 400 且不扣费；提交时合法、投递时解析到私网 → `ssrf_blocked`。
  8. 重启恢复：投递中途“崩溃”（取消 dispatcher context、保留 `delivering` 行）→ 新 dispatcher 回收后完成投递，最多重复一次，`event_id` 不变。
  9. 多实例：两个 dispatcher 共享一个库并发领取 → 每个事件只被一个领取、无重复写回。
  10. 本地超时/上游 404/放弃轮询 → GET 与回调都显示 `expired`/`failed` 且带 `error`。
  11. 计费：回调全过程不产生消费/退款日志变化；快路径与轮询同时触发终态时只结算一次。
- 不使用 `time.Sleep` 等待：重试间隔、超时通过可注入时钟/配置缩短。

### 5.3 多数据库迁移验证

- 在 box 用 Docker 起 MySQL 5.7、MySQL 8、PostgreSQL 9.6、PostgreSQL 15 与 SQLite，跑：首次迁移、二次启动无 ALTER、唯一键冲突 DoNothing、领取 CAS `RowsAffected`、清理任务。box 起不来的版本交给 GitHub Actions 矩阵并在 PR 中注明。

### 5.4 真实供应商验证（db51，待用户提供测试 Key）

- db51 内网不可被供应商回调：快路径保持关闭，只验证轮询路径；客户接收端用 db51 本机/内网 mock，临时打开 `TaskCallbackSSRFAllowPrivate`，验证完关闭。
- 验证清单：真实 `queued → running → succeeded` 序列与 body；构造失败任务确认 HTTP 400 失败体的真实形状与我们的合成体；`expired` 用 `execution_expires_after=3600` 观测（耗时较长，可选）；确认供应商对 `callback_url` 的真实行为（若有可公网访问的临时环境，再验证快路径与 Q4）。
- 记录结果到本文进度区。

---

## 6. 实施步骤（每步一个 commit，Conventional Commits）

1. `docs: add Seedance task callback plan` —— 本文件（等待用户评审后再提交）。
2. `fix(task): keep task data consistent with locally decided terminal states` —— 前置修正 4.1，含单测。
3. `fix(doubao): treat supplier HTTP 400 failed-task body as failure` —— 前置修正 4.3，插件 `parseTaskResult` + 轮询分类调整 + 单测。
4. `feat(task): honor execution_expires_after for task timeout` —— 前置修正 4.2，含部署文档变量说明。
5. `feat(task): add task callback settings` —— `task_callback_setting` 注册、默认值、选项读写 API（不含 UI）。
6. `feat(task): add task callback subscription and outbox models` —— 两张表、索引、AutoMigrate、保留期清理 system task、三库迁移测试。
7. `feat(task): validate customer callback_url with SSRF protection` —— 独立校验器 + 拨号器复用 + 单测。
8. `feat(seedance): accept callback_url on video task submit` —— 别名路由读取并校验 `callback_url`（预扣费前）、提交成功后建订阅、从上游请求与持久化 metadata 中剔除；插件注释更新。
9. `refactor(task): share Seedance task status renderer between GET and callbacks` —— 抽出渲染函数，GET 行为不变（快照测试锁定字节）。
10. `feat(task): enqueue task callback events on status transitions` —— `UpdateWithStatusTx` + `EnqueueTaskCallbackIfChanged`，接入全部状态写入路径 + 对账扫描。
11. `feat(task): deliver task callbacks with retries and leases` —— dispatcher、领取/写回 CAS、重试策略、超时、并发、重启回收、结构化日志。
12. `feat(task): add optional HMAC signing for task callbacks` —— 默认关闭。
13. `feat(task): add upstream callback fast path` —— 公网回调端点、token、单任务刷新、singleflight、渠道开关、注入策略。
14. `test(task): add end-to-end callback tests with mock supplier and receiver` —— 5.2 全部场景。
15. `feat(web): show task callback deliveries and settings` —— 后台设置分组、任务详情页签、列表页、重投按钮、i18n。
16. `feat(monitor): add task callback metrics and alerts` —— 若监控配置不在本仓库，则在部署包（`deliver-v2`）中同步修改 sql_exporter 查询、监控账号授权、告警规则、Grafana 面板，并在 PR 中说明。
17. `docs: document Seedance task callbacks` —— 面向客户的回调说明（契约、重试、去重、签名校验示例）、部署文档变量、本计划进度区勾选。
18. 发布：推送 `feat/seedance-callback` 分支 → PR → CI（Go 测试 + 前端 Actions）→ 合并 → release-please 生成版本 PR（配置 `bump-patch-for-minor-pre-major: true`，`feat:` 在 1.0 前只升 patch，预计 `0.0.3`）→ 合并 → 确认 `CHANGELOG.md`、`VERSION`、Docker Hub `dex2all/ktapi:<ver>` → 更新部署包 `KTAPI_TAG`。（本步需用户明确授权后执行，当前不推送。）

每步验收：`go build ./...`、`go vet`、相关包 `go test -count=1` 通过；涉及数据库的步骤附三库验证结果。

---

## 7. 风险与缓解

| 风险 | 缓解 |
|---|---|
| 客户 URL 指向内网（SSRF） | 提交时 + 拨号时双重校验、不跟随重定向、不走代理、独立配置 |
| 慢/挂死的客户接收端拖垮投递 | 5 秒超时、per-host 并发上限、全局 worker 上限、非终态单次 |
| 回调与 GET 不一致 | 共用渲染函数、入队时快照、字节级测试 |
| 上游回调伪造 | 每任务随机 token + 不信任 body，只触发重新查询 |
| 重复回调 | 唯一键 `(task_id,status)` 保证只入队一次；投递是“至少一次”，提供 `X-KtApi-Event-Id` |
| 多实例重复投递 | 条件 UPDATE 领取 + 租约 + 带 owner 的写回 |
| 上游 ID / 渠道信息泄露 | ID 替换 + 子串扫描测试；日志不打印 URL 路径/查询串 |
| outbox 无限增长 | 保留期清理任务；监控积压 |
| 快路径在内网环境不可用 | 默认关闭，轮询始终兜底 |

---

## 8. 已确认的决策（2026-10-19 用户评审）

- **Q1 重试次数**：“最多重试三次” = 首次 + 3 次重试，共 4 次，间隔约 5 秒，单次超时 5 秒，仅 `succeeded` / `failed`；`queued` / `running` / `expired` 只发一次。长退避“扩展重试”保留实现但默认关闭。
- **Q2 失败任务**：GET 与回调都用方舟形状，HTTP 200，`status:"failed"` + `error` 对象。
- **Q3 `cancelled`**：不触发回调（GET 仍如实返回）。
- **Q4 / Q10 / Q11 上游快路径**：实现，但默认关闭；db51 只用轮询。
- **Q5 额外请求头**：`X-KtApi-*` 默认开启，可在设置中关闭。
- **Q6 视频 URL**：不代理，保持现状（不在本次范围）。
- **Q7 校验错误形状**：沿用 Seedance 别名路由现有的参数错误信封（插件 `error` 渲染器 + 宿主状态码），保持一致。
- **Q8 存储**：测试阶段明文存储客户回调 URL（日志只打印 host + 哈希）。
- **Q9 本地超时**：kt-api 自身超时对外报 `expired`；按任务 `execution_expires_after`（默认 172800，最大 259200）+ 宽限计算截止时间；db51 与部署包 `TASK_TIMEOUT_MINUTES` 调为 4350。
- **Q12 HTTPS-only**：默认允许 http，提供可选的 https-only 开关。

## 9. 进度日志

| 日期（Asia/Shanghai） | 步骤 | 说明 |
|---|---|---|
| 2026-10-19 | Step 1 | 计划评审通过，记录决策并提交 |
| 2026-10-19 | Step 2 | 本地裁定的终态（超时/查询失败/未找到）写入 Ark 形态 data（status+error），GET 与回调同源 |
| 2026-10-19 | Step 3 | doubao 插件将供应商 HTTP 400 失败任务裸 error 体识别为 FAILURE，持久化为 Ark failed 形态 |
| 2026-10-19 | Step 4 | 支持 execution_expires_after（3600–259200，默认 172800，+1800 宽限），独立分页超时扫描，超时报 expired |
| 2026-10-19 | Step 5 | 新增 task_callback_setting 配置组（默认关闭 fast path、http 允许、X-KtApi-* 头开启） |
| 2026-10-19 | Step 6 | 新增 task_callback_subscriptions / task_callback_events 表与模型、清理系统任务；SQLite 3.50、MySQL 5.7.44/8.0、MariaDB 11.8、PostgreSQL 9.6/17 迁移与并发认领验证通过 |
| 2026-10-19 | Step 7 | 独立 callback_url 校验器（scheme/https-only/userinfo/端口/非规范数字 IP/私网/解析全部合法）与投递客户端（不走代理、不跟随重定向、拨号时逐 IP 校验防 DNS rebinding） |
| 2026-10-19 | Step 8 | 别名路由在 Distribute/预扣费前读取并校验 callback_url（非字符串/不合法 → 400 插件错误信封），任务落库后建订阅；callback_url 不发往上游、不入任务 metadata |
| 2026-10-19 | Step 9 | 抽出 RenderTaskPluginBody / RenderSeedanceTaskStatus，GET 直接写渲染字节；快照测试锁定 GET 字节且与回调渲染逐字节一致、不含上游 ID |
| 2026-10-19 | Step 10 | UpdateWithStatus 在 CAS 赢家事务内（savepoint）调用观察者入队，失败不阻塞任务写；cancelled 只记 last_status；每分钟对账系统任务（有未终结订阅时才运行）补入队；渠道缺失批量失败改写 Ark failed data |
