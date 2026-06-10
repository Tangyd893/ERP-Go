# ADR-005: 订单履约事件架构（首期）

## 状态

已采纳（2026-05-25）

## 背景

Order 服务已实现进程内 Outbox 轮询处理 `order.approved` / `order.cancelled`；WMS 出库完成事件需跨服务通知 Order 完成扣库存与发货。

## 决策

### 通道 1：Order 进程内 Outbox（已有）

- 触发：本服务写 Outbox → `StartPolling` 消费
- 范围：`order.approved` → 锁库 + 建出库单；`order.cancelled` → 释放库存事件

### 通道 2：WMS → Order HTTP 回调（首期跨服务）

- 触发：`POST /api/v1/warehouse/outbounds/:id/ship`
- WMS 更新出库状态后 HTTP 调用 `POST /api/v1/order/fulfillment/outbound-shipped`
- Order 服务同步执行 `HandleOutboundShipped`（扣库存 + 更新订单状态 + 写 Outbox）

### 通道 3：RabbitMQ Consumer（已接入 Order 服务）

- 队列：`order.fulfillment`，绑定 `warehouse.outbound.shipped`
- Order 服务 DB 就绪且配置 `RABBITMQ_URL` 时自动启动
- 与 HTTP 回调互为补充；Inbox 幂等避免双通道重复处理

## 后果

- 首期不依赖独立 worker 进程即可跑通履约尾段
- WMS 与 Order 需网络可达（`ORDER_SERVICE_URL`）
- 幂等由 Inbox `messageID = ship-{outbound_id}` 保证
