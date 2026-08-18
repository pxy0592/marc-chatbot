# Data Model: 微信群聊机器人管理

## BotInstance

| Field | Type | Rules |
|---|---|---|
| id | string | 固定单实例标识，主键 |
| public_account_id | string? | 登录后可见的公开账号标识 |
| display_name | string? | 登录后可见名称 |
| status | enum | `unbound`, `awaiting_scan`, `connecting`, `online`, `reconnecting`, `offline`, `error` |
| status_reason | string? | 非敏感、可展示原因 |
| qr_code | string? | 仅绑定期间保存；完成、取消或过期后清空 |
| qr_expires_at | timestamp? | 必须晚于生成时间 |
| last_seen_at | timestamp? | 最近在线或事件时间 |
| updated_at | timestamp | 每次状态变化更新 |

### State transitions

```text
unbound -> awaiting_scan -> connecting -> online
awaiting_scan -> unbound        (cancelled/expired)
connecting -> error/offline
online -> reconnecting/offline/error
reconnecting -> online/offline/error
offline/error -> awaiting_scan   (new bind) or connecting (provider reconnect)
any -> unbound                   (explicit unbind)
```

## BindingSession

| Field | Type | Rules |
|---|---|---|
| id | string | 主键 |
| status | enum | `pending`, `scanned`, `confirmed`, `expired`, `cancelled`, `failed` |
| qr_code | string? | 仅 pending/scanned 可返回 |
| created_at | timestamp | 必填 |
| expires_at | timestamp | 必须大于 created_at |
| completed_at | timestamp? | 终态时设置 |
| failure_reason | string? | 不含令牌和凭证 |

同一 bot 同时最多一个非终态 BindingSession。

## Group

| Field | Type | Rules |
|---|---|---|
| id | string | 内部主键 |
| provider_group_id | string | provider 稳定标识，唯一 |
| name | string | 可变化，不作为唯一键 |
| member_count | integer? | 非负，未知时为空 |
| enabled | boolean | 默认 false |
| available | boolean | 当前机器人是否仍可访问 |
| discovered_at | timestamp | 首次发现时间 |
| updated_at | timestamp | 最后同步时间 |

## GroupMessage

| Field | Type | Rules |
|---|---|---|
| id | string | 内部主键 |
| provider_message_id | string | provider 消息标识，唯一去重 |
| group_id | string | 外键 -> Group |
| sender_id | string | provider 联系人标识 |
| sender_name | string | 接收时快照 |
| direction | enum | `inbound`, `outbound` |
| message_type | enum | `text`, `unsupported`, `system` |
| content | string | 文本或安全摘要；最大 4000 字符 |
| occurred_at | timestamp | provider 消息时间 |
| received_at | timestamp | 服务接收时间 |
| processing_status | enum | `received`, `ignored`, `sent`, `failed` |
| self_message | boolean | 机器人自身消息为 true |

## OutboundMessageRequest

| Field | Type | Rules |
|---|---|---|
| id | string | 主键 |
| idempotency_key | string | 唯一，1-128 字符 |
| group_id | string | 外键 -> Group |
| content | string | trim 后 1-2000 字符 |
| requested_by | string | 管理员标识 |
| status | enum | `pending`, `sending`, `succeeded`, `failed` |
| provider_message_id | string? | 成功后保存 |
| failure_reason | string? | 失败后保存非敏感原因 |
| created_at | timestamp | 必填 |
| updated_at | timestamp | 每次状态变化更新 |

### State transitions

```text
pending -> sending -> succeeded
pending -> failed
sending -> failed
```

终态请求不可再次发送；相同 idempotency_key 返回已有请求。

## AuditEvent

| Field | Type | Rules |
|---|---|---|
| id | string | 主键 |
| actor | string | 管理员标识或 `system` |
| action | enum | `binding.start`, `binding.cancel`, `bot.login`, `bot.logout`, `group.enable`, `group.disable`, `message.send` |
| target_type | string | `bot`, `binding`, `group`, `message` |
| target_id | string | 目标稳定标识 |
| result | enum | `success`, `failure` |
| detail | string? | 安全、非敏感摘要 |
| created_at | timestamp | 必填，按时间倒序查询 |

## Relationships

```text
BotInstance 1 --- * BindingSession
BotInstance 1 --- * Group (运行时发现)
Group       1 --- * GroupMessage
Group       1 --- * OutboundMessageRequest
OutboundMessageRequest 0..1 --- 1 GroupMessage (成功发送后)
AuditEvent 以 target_type + target_id 关联任意业务对象
```

## Retention

- GroupMessage 和 AuditEvent 默认保留 30 天。
- BindingSession 的二维码内容在进入终态后立即清除，记录本身可保留用于审计。
- BotInstance 不持久化 Puppet token 或微信登录凭证。
