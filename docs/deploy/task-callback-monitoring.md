# 任务回调监控（Seedance callback_url）

监控配置不在本仓库，而在部署包（`deliver-v2` / `ktapi-deploy-<ver>.tar.gz`）里。本文是部署包修改的参考副本，升级到 ktapi 0.0.3 时按下面顺序操作。

## 升级顺序

1. 先升级 ktapi 到 0.0.3。启动时 AutoMigrate 会创建 `task_callback_subscriptions` 和 `task_callback_events`。
2. 重新执行 `node-agent/create-monitor-role.sql`。这个脚本可以重复执行，会补上新表的列级授权。
3. 再更新 `node-agent/sql_exporter/ktapi_realtime.collector.yml` 并重启 `agent-sql-exporter`，然后更新 Prometheus 规则和 Grafana 看板。

如果顺序反了（表还不存在或没有授权），sql_exporter 的查询会失败，并触发 `BusinessMetricsMissing`。

## 监控账号授权（列级，只读）

```sql
GRANT SELECT (id, state, status, attempts, created_at, updated_at, lease_until)
  ON task_callback_events TO ktapi_monitor;
```

不授予 `payload`（客户可见的任务快照）、`lease_owner`、`last_error`，也不授予 `task_callback_subscriptions` 的任何列（里面有客户回调 URL）。

## 指标（sql_exporter，collector `ktapi_realtime`）

| 指标 | 含义 |
|---|---|
| `ktapi_task_callback_backlog{state="pending"\|"delivering"}` | 待投递 / 投递中的事件数 |
| `ktapi_task_callback_oldest_pending_age_seconds` | 最早一条 pending 事件已等待的秒数（无积压为 0） |
| `ktapi_task_callback_finished_15m{state="succeeded"\|"failed"}` | 最近 15 分钟送达 / 放弃的事件数 |
| `ktapi_task_callback_lease_stuck` | 租约过期 60 秒以上仍为 delivering 的事件数 |

查询已在 PostgreSQL 9.6 上验证。

## 告警（`monitor/prometheus/rules/business.rules.yml`）

| 告警 | 条件 | 级别 |
|---|---|---|
| `TaskCallbackBacklog` | `ktapi_task_callback_oldest_pending_age_seconds > 300`，持续 5 分钟 | warning |
| `TaskCallbackFailureRateHigh` | 15 分钟内放弃占比 > 50% 且至少 5 条，持续 10 分钟 | warning |
| `TaskCallbackLeaseStuck` | `ktapi_task_callback_lease_stuck > 0`，持续 10 分钟 | warning |

## 日志看板（`monitor/grafana/dashboards/ktapi-logs.json`）

- **任务回调投递（按结果）**：用 Loki 统计 `task_callback event=delivery_attempt` 日志里的 `result=` 字段，取值有 succeeded / timeout / connect / dns / tls / http_4xx / http_5xx / ssrf_blocked 等。
- **任务回调日志（task_callback）**：显示所有 `task_callback event=` 结构化日志，包括入队、每次投递尝试、上游快路径提示、租约回收和清理。

日志只带回调地址的 `scheme://host` 和 URL 哈希前缀，不含路径和查询串。
