# Redis 鉴权缓存键升级说明（2026-11-07 上游同步）

上游 `2de51a19e35a602a0041d380b29be7c9bcd1e0e5` 修复 Redis Cluster 的跨槽 Lua/MGET：

| 用途 | 旧键 | 新键 |
| --- | --- | --- |
| 用户缓存 | `user:<id>` | 不变 |
| 令牌缓存 | `token:<HMAC>` | 不变 |
| 用户更新围栏 | `auth:user:fence:<id>` | `auth:user:fence:{user:<id>}` |
| 用户已提交版本 | `auth:user:version:<id>` | `auth:user:version:{user:<id>}` |
| 令牌围栏 | `token:fence:<HMAC>` | `token:fence:{token:<HMAC>}` |

花括号内的内容与原缓存键完全一致，因此相关键在同一 Redis 槽位。缓存 hash 名不变，配额计数不拆分；但旧、新版本不共享围栏/版本键，**不支持新旧节点长时间混跑的滚动升级**。

## 部署边界

- 单实例：先拉取镜像，再由 Compose 停掉旧容器并替换，避免同时处理请求；不得用第二实例流量灰度替代原容器。
- 多实例：先摘流并排空所有旧节点，再协调切换；不要在新旧节点并存期间修改权限、禁用用户、撤销令牌或重置缓存。
- 如停止前有未完成的安全更新或缓存回填失败，先排查数据库鉴权状态与缓存一致性；旧围栏不可当作新版本仍能读取的安全保障。需要清理/迁移缓存时先单独确认方案并备份，禁止直接 FLUSHDB/FLUSHALL（会波及会话等数据）。
- 本次自动化只升级 ktapi 应用，不重启或清空 Redis/PostgreSQL，不删除鉴权缓存，也不清日志。
- 回滚也应全体协调，不能将旧镜像与新镜像混用；数据库备份及升级前 `.env` 必须保留。

## 回归

`TestAuthCacheKeysShareClusterHashTag` 校验两类缓存的 fence/version hash tag 与缓存键一致。缓存完整性、配额保留、并发失效及权限版本测试与后端全量测试共同执行。
