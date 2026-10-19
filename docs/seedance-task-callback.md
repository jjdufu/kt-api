# Seedance 任务回调（callback_url）

kt-api 的 Seedance 视频生成别名接口 `POST /api/v3/contents/generations/tasks` 支持 Ark 文档里的 `callback_url`。任务状态变化后，kt-api **自己**向客户的 `callback_url` 推送任务快照。上游供应商不会直接调用客户地址。

## 客户约定

### 提交

```json
{
  "model": "doubao-seedance-…",
  "content": [ … ],
  "callback_url": "https://example.com/hooks/seedance",
  "execution_expires_after": 172800
}
```

- `callback_url` 可选。不合法时，提交直接返回 400（错误信封与别名接口其它参数错误相同），**不扣费、不建任务**。不合法包括：不是字符串、协议不允许、带 userinfo、端口不允许、非规范数字 IP、解析到私网 / 回环 / 链路本地 / 元数据地址、域名无法解析等。
- `callback_url` 不会发给上游，也不出现在 GET 返回里。
- `execution_expires_after` 取值 3600–259200 秒，默认 172800。

### 回调请求

- `POST <callback_url>`，`Content-Type: application/json`，`User-Agent: kt-api-callback/<version>`。
- **body 与同一时刻 `GET /api/v3/contents/generations/tasks/{id}` 的返回字节完全相同**（Ark 任务对象：`id`、`model`、`status`、`content`、`usage`、`error`、`created_at`、`updated_at` 等），不含上游任务 ID。
- 会推送的 `status`：`queued`、`running`、`succeeded`、`failed`、`expired`。
  - `failed` 采用 Ark 形状：`status: "failed"` 加上 `error: {code, message}`。
  - kt-api 本地超时（超过 `execution_expires_after` + 30 分钟宽限仍未完成）报 `expired`。
  - **`cancelled` 不推送回调**，最终状态以 GET 为准。
- 同一任务的事件按顺序串行投递。如果更新的状态已入队，尚未送达的 `queued` / `running` 事件会被跳过（supersede），终态一定会投递。
- 视频地址 `content.video_url` 原样透传上游签名 URL（kt-api 不代理），会过期，请及时下载。

### 送达与重试

- 收到任意 **2xx** 即视为送达，响应体会被忽略。单次请求超时 5 秒。
- `succeeded` / `failed`：最多 **4 次**尝试（首次 + 3 次重试），间隔 **5 秒**（调度粒度 1 秒）。
- `queued` / `running` / `expired`：只发 **1 次**，不重试。
- 3xx 不跟随，按失败计。投递时目标解析到被禁止的地址（`ssrf_blocked`）会立即放弃，不重试。
- 可选的扩展退避默认关闭（开启后，终态事件在 24 小时窗口内继续按指数退避重试）。
- 所有尝试都失败后事件标记为 failed。管理员可在“系统设置 → 运维设置 → 任务回调”里重投（`event_id` 不变）。

### 去重

默认开启附加请求头（`task_callback_setting.extra_headers_enabled`）：

| 头 | 含义 |
|---|---|
| `X-KtApi-Event-Id` | 事件唯一 ID，重试和重投时**不变**，请据此去重 |
| `X-KtApi-Delivery-Attempt` | 第几次尝试（从 1 开始） |
| `X-KtApi-Task-Id` | 任务 ID（与 body.id 相同） |

至少一次投递：同一事件可能重复到达（例如超时后其实已经收到），接收端必须幂等。

### 签名校验（可选，默认关闭）

开启 `task_callback_setting.signing_enabled` 并配置主密钥 `signing_secret` 后，每个回调会额外带上：

- `X-KtApi-Timestamp`：Unix 秒
- `X-KtApi-Signature: v1=<hex>`，其中 `<hex> = HMAC-SHA256(用户密钥, timestamp + "." + event_id + "." + 原始 body)`

用户密钥由主密钥按用户派生，不落库。用户登录控制台后可通过 `GET /api/user/task_callback/signing` 查看。更换主密钥会让所有用户的密钥同时失效（不支持双签名过渡）。签名需要 `X-KtApi-Event-Id`，因此请同时保持附加请求头开启。

Python 校验示例：

```python
import hmac, hashlib, time

def verify(secret: str, headers: dict, raw_body: bytes, tolerance=300) -> bool:
    ts = headers["X-KtApi-Timestamp"]
    event_id = headers["X-KtApi-Event-Id"]
    if abs(time.time() - int(ts)) > tolerance:
        return False
    mac = hmac.new(secret.encode(), f"{ts}.{event_id}.".encode() + raw_body, hashlib.sha256)
    expected = "v1=" + mac.hexdigest()
    return any(hmac.compare_digest(expected, s.strip())
               for s in headers["X-KtApi-Signature"].split(","))
```

Go 版参考实现：`service.VerifyTaskCallbackSignature`。测试向量：`secret="k"`、`ts=1760000000`、`event_id="evt_1"`、`body="{}"` 得到 `v1=216affa7e5e4e57a43d51be7126cf056a210fe3e94808c738254c392d19f4d6a`。

## 部署与配置

### 环境变量

| 变量 | 建议 |
|---|---|
| `TASK_TIMEOUT_MINUTES` | **≥ 4350**（72h 最大 `execution_expires_after` + 30 分钟宽限）。Seedance 任务按各自截止时间判超时；该值过小时，日志会提示把它调大 |

### 选项（系统设置 → 运维设置 → 任务回调；存储为 `task_callback_setting.*`）

| 选项 | 默认 | 说明 |
|---|---|---|
| `enabled` | true | 接受 `callback_url` 并运行投递 worker。关闭后新提交的 `callback_url` 被忽略，已入队事件保持 pending，重新开启后继续投递 |
| `extra_headers_enabled` | true | `X-KtApi-*` 头 |
| `https_only` | false | 只允许 https |
| `allow_private_ip` | false | 仅限测试 |
| `domain_allowlist_mode` / `domain_list` | false / [] | 域名白名单（支持 `*.example.com`）|
| `allowed_ports` | 80, 443, 8080, 8443 | 支持区间 `8000-9000` |
| `max_url_length` | 2048 | |
| `reject_unresolvable` | true | 提交时域名解析失败即拒绝 |
| `timeout_seconds` / `retry_interval_seconds` / `terminal_max_attempts` | 5 / 5 / 4 | |
| `extended_retry` | false | 扩展退避 |
| `workers` / `per_host_concurrency` | 16 / 4 | 每实例并发 / 每目标主机并发 |
| `retention_days` | 7 | 已结束事件的保留天数 |
| `signing_enabled` / `signing_secret` | false / 空 | HMAC 签名 |
| `upstream_public_base` / `upstream_inject_mode` | 空 / only_when_customer_callback | 上游回调快路径 |

### 上游回调快路径（默认关闭）

需要同时满足：配置了公网可达的 `upstream_public_base`，并且渠道打开了“上游任务回调快路径”（`upstream_task_callback_enabled`）。开启后，提交时会把 `<base>/api/task-callback/upstream/<一次性 token>` 作为上游 `callback_url`。上游回调只用来触发一次即时轮询，**不信任 body**。db51 / 内网部署无公网入口，保持关闭，只靠轮询。

### 数据表与监控

0.0.3 启动时自动创建 `task_callback_subscriptions`（存储明文回调 URL）和 `task_callback_events`。监控指标、授权和告警见 [deploy/task-callback-monitoring.md](deploy/task-callback-monitoring.md)。

### 排查

- 任务详情 → “回调”区块：每个事件的状态、尝试次数、HTTP 状态和错误（只显示主机名和 URL 哈希前缀）。
- 日志搜索 `task_callback event=`，例如 `event=delivery_attempt ... result=http_5xx`。
