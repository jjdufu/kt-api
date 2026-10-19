# Seedance 素材与视频支持实现计划

基线：`jjdufu/kt-api`（fork of QuantumNous/new-api @ `058c823`），当前 `main` 版本 `0.0.1`。  
供应商文档：`/workspace/uploads/apidoc.txt`（14 个接口：2 视频 + 12 素材）。  
工作分支：`feat/seedance-assets`。

## 目标

在 kt-api 上原生支持 Seedance（易速 / yisu-api）的：

1. **12 个素材管理 Action**（同步 HTTP 转发，带用户所有权隔离、渠道 pin、计费日志）
2. **视频任务路径别名** `/api/v3/contents/generations/tasks`（仍走 Doubao JS 插件，渠道类型 45/54）
3. **`asset://` 引用** 时强制 pin 到素材创建时的同一渠道（及多 key 时的 key 指纹）
4. **`callback_url` 处理**：剥离或由 kt-api 侧代收，避免供应商直连客户回调导致渠道不一致等问题

## 非目标（明确不做）

- **CNY / 人民币定价与币种换算**（用户已延期）
- 管理后台 UI（除非原生定价配置绝对必要；本阶段优先用已有 model price / 虚拟模型机制）
- 部署到 db51
- 新增 channel type 60；Seedance 视频继续用 **45（VolcEngine）/ 54（DoubaoVideo）+ yisu base URL**
- 用 JS 插件注册 `/api/*`（`pkg/jsplugin/routing.go` 的 `reservedRouteNamespaces` 禁止）

## 14 个接口映射

| # | 方法 | 路径 / Action | 实现路径 | 虚拟模型 / 计费 |
|---|------|----------------|----------|-----------------|
| 1 | POST | `/api/v3/contents/generations/tasks` | 视频别名 → Doubao 插件 submit | 既有 Seedance 模型价 |
| 2 | GET | `/api/v3/contents/generations/tasks/:id` | 视频别名 → Doubao 插件 query | 查询不计费 |
| 3 | POST | `/api/assets/?Action=CreateAssetGroup&Version=2024-01-01` | 原生 asset relay | `seedance-asset`（价 0，仍记日志） |
| 4 | POST | `...Action=CreateAsset` | 原生（JSON + multipart） | `seedance-asset` |
| 5 | POST | `...Action=ListAssetGroups` | 原生；**强制** `Filter.GroupIds` = 调用方已有 ID | `seedance-asset` |
| 6 | POST | `...Action=GetAssetGroup` | 原生；所有权校验 | `seedance-asset` |
| 7 | POST | `...Action=ListAssets` | 原生；组归属校验 | `seedance-asset` |
| 8 | POST | `...Action=GetAsset` | 原生（JSON + multipart 兼容） | `seedance-asset` |
| 9 | POST | `...Action=UpdateAssetGroup` | 原生 | `seedance-asset` |
| 10 | POST | `...Action=UpdateAsset` | 原生 | `seedance-asset` |
| 11 | POST | `...Action=DeleteAsset` | 原生；上游失败则保留本地行 | `seedance-asset` |
| 12 | POST | `...Action=DeleteAssetGroup` | 原生；级联本地素材行 | `seedance-asset` |
| 13 | POST | `...Action=CreateVisualValidateSession` | 原生；会话创建时预扣/按次 | `seedance-visual-validate`（按次） |
| 14 | POST | `...Action=GetVisualValidateResult` | 原生；hash 匹配；成功后登记 LivenessFace 组 | `seedance-visual-validate-result`（价 0，仍记日志） |

供应商认证：`Authorization: Bearer <渠道 Key>`，Base URL 取自被 pin 的渠道（典型 `https://yisu-api.gigaai.cc`）。

## 架构

```
客户端
  │  TokenAuth
  ▼
┌─ 视频别名 ─────────────────────────────────────────────┐
│  POST/GET /api/v3/contents/generations/tasks[/:id]      │
│    → 改写/挂到 Doubao 插件路由等价路径                    │
│    → PrepareTaskPlugin + asset:// pin + Distribute      │
│    → RelayTask / RelayTaskFetch                         │
└─────────────────────────────────────────────────────────┘

┌─ 素材同步路径（非 task 框架）────────────────────────────┐
│  POST /api/assets/                                      │
│    → PrepareAssetRoute（解析 Action/Version 白名单）     │
│    → 所有权/pin（读本地表 → ChannelPin + key fingerprint）│
│    → ModelPriceHelperPerCall + NewBillingSession         │
│    → AssetHandler 转发供应商（透传状态码与信封）          │
│    → 成功后写/删本地 ownership；Settle / 失败 Refund     │
│    → RecordConsumeLog（零价也记）                        │
└─────────────────────────────────────────────────────────┘
```

要点：

- **插件不能挂 `/api/*`** → 素材必须原生 Go。
- **Pin 先例**：`middleware/task_plugin.go` 的 `applyOriginTaskIntent` → `dto.ChannelPin` → `Distribute`。素材与 `asset://` 视频沿用同一约束结构；必要时扩展 `PinSourceAsset` 与 key 指纹字段。
- **计费**：`ModelPriceHelperPerCall`、`NewBillingSession` / Settle / Refund、`RecordConsumeLog`。真人认证在 **CreateVisualValidateSession** 时按次计费；回调查询串里的 `reqMeasureInfoValue` **仅给客户 CallbackURL 使用**，不作为 kt-api 二次扣费依据。
- **BytedToken**：本地只存 SHA-256；30 分钟过期；查询成功后删除会话摘要（单次使用）。
- **ListAssetGroups**：无论客户端是否传 `Filter.GroupIds`，服务端注入/收紧为当前用户拥有的 ID 集合（空集则不应打到供应商列出他人组）。
- **渠道**：推荐单 key；多 key 时在组/会话上记录 key fingerprint，后续请求 pin 同一 channel + 同一 key。

## 数据表

| 表 | 用途 |
|----|------|
| `asset_groups` | 用户 ↔ 上游 `group-*`；`channel_id`、key 指纹、`group_type`（AIGC / LivenessFace）、`project_name` |
| `assets` | 用户 ↔ 上游 `asset-*`；归属 `group_id`、同一 channel/key |
| `visual_validate_sessions` | `byted_token_hash`、user、channel/key、`project_name`、`expires_at`；不存明文 token |

`model/main.go` `AutoMigrate` 注册上述模型。本地以 SQLite 验证迁移幂等；MySQL/Postgres 尽量在 box 跑，否则写测试并在 PR 中标明依赖 Actions 矩阵。

## 计划修改 / 新增文件

| 区域 | 文件（拟） |
|------|------------|
| 计划 | `docs/seedance-assets-plan.md`（本文，用户明确要求） |
| 模型 | `model/asset.go`、`model/asset_test.go`；`model/main.go` AutoMigrate |
| 路由 | `router/asset-router.go`；`router/main.go` 挂载；视频别名挂在 `router/video-router.go` 或同级 |
| 中间件 | `middleware/asset_route.go`（PrepareAssetRoute、所有权、pin） |
| 控制器 | `controller/asset.go` |
| 转发 | `relay/asset_handler.go`（或 `service/asset_relay.go`） |
| DTO | `dto/channel_constraints.go` 扩展 pin source / key；必要时 `dto/asset.go` |
| 计费/价格 | `setting/ratio_setting/model_ratio.go` 默认虚拟模型价；消费日志路径 |
| 插件 | `plugins/tasks/doubao/plugin.js`：`asset://` intent、callback_url 处理 |
| 测试 | 集中一个 `*_test.go`（httptest mock 供应商），覆盖跨用户拒绝、pin、计费、token hash、List 注入、multipart、视频别名 |
| 变更说明 | 按 AGENTS：本 fork 用 release-please 维护根 `CHANGELOG.md`；**不**在 `plugins/` 下加文档。用户可见说明写入 PR 正文；如需独立说明可用本计划文件进度区 |

## 实施步骤与验收

### Step 0 — 计划文档（当前）

- [x] 撰写本计划（简体中文，无时间估算）
- [x] 单独 commit：`docs: add Seedance assets implementation plan`

**验收**：文件存在；commit 仅含计划。

### Step 1 — 数据模型 + AutoMigrate + 所有权 helper 单测

- [ ] `AssetGroup` / `Asset` / `VisualValidateSession` 模型与索引（user_id + 上游 ID 唯一）
- [ ] AutoMigrate 注册；SQLite 起库两次证明幂等
- [ ] Helper：按 user 查组/素材、登记、删除、List 用 ID 列表、token hash/过期/单次使用
- [ ] 单测覆盖跨用户 miss、过期 token、重复消费拒绝

**验收**：`go test ./model/ -count=1 -run Asset`（或等价）通过；`go build ./...` 通过。

### Step 2 — 路由 + PrepareAssetRoute + controller/handler 骨架

- [ ] `POST /api/assets/` 注册（TokenAuth + rate limit + Prepare + Distribute + handler）
- [ ] Action/Version 白名单；非法 Action → 火山信封错误
- [ ] Handler 能把请求转发到渠道 BaseURL（无业务侧写库亦可先通）

**验收**：本地用 mock 上游打通 CreateAssetGroup 往返；`go vet` 干净。

### Step 3 — 12 Action 的所有权 / pin / multipart

- [ ] 写操作成功后落本地行；读/改/删前校验归属
- [ ] ListAssetGroups 注入 GroupIds
- [ ] CreateAsset / GetAsset multipart 兼容
- [ ] 删除：上游失败保留本地；删组级联本地 assets
- [ ] GetVisualValidateResult 成功 → 登记 LivenessFace 组并删会话
- [ ] pin：`ChannelPin` + 多 key fingerprint

**验收**：跨用户 ID 返回 record not found / 等价；同用户 pin 到创建渠道。

### Step 4 — 虚拟模型计费 + 消费日志

- [ ] 默认价：`seedance-asset`=0，`seedance-visual-validate`=按次（可配置），`seedance-visual-validate-result`=0
- [ ] 全路径 `NewBillingSession` → 成功 Settle / 失败 Refund
- [ ] 零价仍 `RecordConsumeLog`
- [ ] 遵守 `.agents/rules/billing.md`（无负向扣费、预扣后结算）

**验收**：测试断言 settle/refund 与日志 model 名；真人认证仅在 create session 扣费。

### Step 5 — 视频别名 + plugin `asset://` intent + callback_url

- [ ] 别名路由挂到与 Doubao 插件等价的 submit/query
- [ ] plugin 从 `content[].*url` 收集 `asset://` → intent → middleware pin（同渠道/key）
- [ ] `callback_url`：剥离供应商回调或改为 kt-api 中转后再通知客户（二选一，实现时选更安全且少改插件协议的方案并在 commit 说明）
- [ ] 错误信封尽量与供应商文档对齐（`error.code/type/request_id`）

**验收**：带 `asset://` 的提交命中 pin；别名路径可编译并由测试覆盖路由存在。

### Step 6 — httptest mock 供应商 e2e 测试

集中测试文件，覆盖：

- [ ] 跨用户 deny
- [ ] pin 同渠道（及 key fingerprint）
- [ ] billing settle / refund
- [ ] BytedToken hash / 过期 / 单次
- [ ] ListAssetGroups 注入
- [ ] multipart CreateAsset
- [ ] 视频别名路由

**验收**：`go test` 目标包通过；不引入脆弱的随机/sleep 压测。

### Step 7 — 收尾：文档进度、PR、CI、release-please

- [ ] 本计划全部 checkbox 勾选
- [ ] `go build` / `go vet` / 相关测试在 box 通过（前端交 Actions，防 OOM）
- [ ] Push `feat/seedance-assets`；开 PR → 等 CI → 必要时 approve Actions → merge
- [ ] merge 后等 release-please 开版本 PR（`feat:` 提交触发）→ approve smoke → merge
- [ ] 确认 Docker Hub `dex2all/ktapi:<ver>`（当前 main 为 `0.0.1`，期望 bump 到 `0.0.2` 或按 release-please 规则）

**验收**：报告 PR 号、CI 状态、镜像 tag 或仍打开的 release PR URL。

## 测试计划（摘要）

- 单元：ownership helpers、token hash、List 注入纯函数
- 集成：httptest 模拟 yisu；内存/SQLite DB；伪造 user/token/channel
- 不做：真实手机 H5 真人认证、真实供应商密钥（需用户提供时再补手工验收）
- DB：优先 SQLite；若 box 有 MySQL/Postgres 容器则跑迁移冒烟，否则注明 Actions 覆盖

## 风险

| 风险 | 缓解 |
|------|------|
| 多租户共用同一供应商 Key | 本地 ownership 强制隔离；List 注入 GroupIds；禁止无过滤上游 List |
| 多 key 渠道导致素材与视频渠道不一致 | 组/会话存 key fingerprint；pin 时锁定；文档建议 Seedance 渠道单 key |
| 渠道被禁用后素材不可用 | pin 时检查 `ChannelStatusEnabled`；返回明确错误（对齐 origin_task_channel_disabled 风格） |
| `/api` 被插件误注册 | 保持 reserved；素材走原生路由 |
| RAM 不足导致本地全量前端 CI OOM | box 只跑 Go；前端交给 GitHub Actions |
| callback_url 直连客户暴露内部细节或绕过计费 | 剥离或 kt-api 中转；计费以 session create 为准 |
| BytedToken 明文泄露 | 只存 SHA-256；日志脱敏 |

## 进度日志

| 日期（Asia/Shanghai） | 步骤 | 说明 |
|----------------------|------|------|
| 2026-10-19 | Step 0 | 创建本计划并单独提交 |
