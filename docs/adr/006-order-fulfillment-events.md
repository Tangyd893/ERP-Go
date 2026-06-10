# ADR-006: 订单履约事件路由与采购入库

## 状态

已采纳（2026-06-09）

## 背景

ADR-005 确立了订单履约的三通道架构（Outbox 轮询 / HTTP 回调 / RabbitMQ Consumer）。采购入库（`inbound.received`）是履约闭环的最后一块——采购单入库后需驱动库存增加，此事件的路由方案需明确。

## 决策

### 事件路由总表

| 事件 | 通道 | 处理器位置 | 幂等键 |
|------|------|-----------|--------|
| `order.approved` | Outbox 轮询 | `order-service` → `P4OutboundFlowCoordinator.HandleOrderApproved` | `msg-{messageID}` |
| `order.cancelled` | Outbox 轮询 | `order-service` → release stock event | `msg-{messageID}` |
| `outbound.shipped` | RabbitMQ Consumer（主）/ HTTP 回调（备） | `order-service` 队列 `order.fulfillment` | `ship-{outbound_id}` |
| `inbound.received` | RabbitMQ Consumer | `warehouse-service` 队列 `warehouse.inbound` | `inbound-{inbound_id}` |

### `inbound.received` 路由方案

**选择：RabbitMQ Consumer 在 `warehouse-service`**

理由：
1. **就近处理**：采购入库的库存增加操作由 `warehouse-service` 的 `P4OutboundFlowCoordinator.HandleInboundReceived` 统一处理，无需跨服务 HTTP 调用
2. **幂等已有**：Inbox 以 `inbound-{inbound_id}` 为幂等键，防止重复入库
3. **实现就绪**：`HandleInboundReceived` 已在 `shared/workflows/p4_outbound_flow.go:344` 实现，并通过 `TestP4HandleInboundReceived` 验证（含幂等重放）
4. **消费端已接线**：`warehouse-service/cmd/server/main.go:79` 在 RabbitMQ 就绪时自动启动 `warehouse.inbound` 队列消费

备选方案（已否决）：
- `purchase-service` 消费：增加跨服务调用链（purchase → warehouse），且 purchase-service 当前无 inbound 协调逻辑
- HTTP 回调：增加网络依赖和重试复杂度，不如 Consumer 天然支持重试+死信

## 后果

- `warehouse-service` 需同时处理出库履约（`outbound.shipped` 回调给 Order）和入库履约（`inbound.received` 直接处理）
- 采购入库的消息发布方（`purchase-service` 或上游 ERP）需向 `warehouse.inbound` 队列发送 `inbound.received` 事件
- 死信/重试策略沿用 RabbitMQ Consumer 的现有机制（见 `shared/outbox/rabbitmq_consumer.go`）
