# ADR-007: 事件驱动双通道职责与扩展计划

日期：2026-05-25
状态：采纳（替代 ADR-005）

## 背景

ERP-Go 当前已落地三条事件通道：
1. **Order 进程内 Outbox 轮询**（处理 `order.approved` / `order.cancelled`）
2. **HTTP 回调**（WMS → Order，通过 HTTP 适配器调用）
3. **RabbitMQ Consumer**（Order 服务消费 `outbound.shipped`）

随着财务、采购、通知等跨服务场景增多，需要在已有工程基础上明确各通道职责边界、避免混淆，并规划采购入库等新场景的处理位置。

## 决策

### 事件通道职责矩阵

| 事件类型 | 推荐通道 | 处理位置 | 理由 |
|---|---|---|---|
| `order.approved` | Outbox 轮询 | `order-service` `StartPolling` | Order 是事件源，仅需通知 inventory/warehouse 执行同步操作 |
| `order.cancelled` | Outbox 轮询 | `order-service` `StartPolling` | 同上 |
| `outbound.shipped` | RabbitMQ Consumer | `order-service` 队列 `order.fulfillment` | 跨服务异步解耦；Inbox 幂等可防重复 |
| `inbound.received` | Outbox 轮询 | `warehouse-service`（计划） | 入库单一服务域内操作；流程短 |
| 财务结算事件 | RabbitMQ（计划） | `finance-service` | 结算异步、高延迟场景 |
| 通知事件 | RabbitMQ（计划） | `notification-service` | 邮件/短信异步发送 |

### 双通道共存策略

- **Outbox 轮询**：仅在事件源服务进程内使用，适用于"通知 1-2 个下游同步操作"的短链路场景。
- **RabbitMQ**：用于跨服务异步解耦场景，配合 Inbox 幂等，支持独立 Consumer 进程扩缩。
- **HTTP 回调**：仅作为 RabbitMQ 不可用时的**降级通道**，不推荐作为主通道。

### Inbox 幂等约定

所有跨服务事件处理必须通过 Inbox 机制保证幂等：
- `messageID` 格式：`{event-type}-{aggregate-id}`（如 `outbound.shipped-OB20250501`）
- Inbox 表唯一索引 `(message_id)` 阻止重复处理
- Consumer handler 与 HTTP 回调 handler 共享 Inbox 去重

### 死信与重试

| 层级 | 重试策略 | 死信处理 |
|---|---|---|
| Outbox 轮询 | 固定间隔 3 次 → failed | 查询 `GET /api/v1/order/outbox/failed` |
| RabbitMQ | exponential backoff → DLQ | DLQ 人工补偿 UI（计划） |

## 采购入库接线

`HandleInboundReceived` 现已在 `shared/workflows` 中有实现和测试（`TestP4HandleInboundReceived`）。计划挂载到 `warehouse-service`：

1. `warehouse-service` 启动时注册 `inbound.received` 事件 handler
2. Handler 调用 `coordinator.HandleInboundReceived`
3. 通过 HTTP 适配器通知 `inventory-service` 增加库存
4. 复用 Outbox 轮询机制（在 warehouse 进程内）

## 后果

- 团队对事件通道选择有明确矩阵，不再每个新场景临时讨论
- HTTP 回调降级为备选，减少对网络直连的依赖
- 采购入库接线已有明确方案，减少实施摩擦
- 死信可查询是后续可观测性的前置条件
