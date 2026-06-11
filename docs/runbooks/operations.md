# 运维手册

## 日常巡检

```bash
# Pod 状态
kubectl get pods -n erp-go -o wide

# 资源使用
kubectl top pods -n erp-go

# 最近事件
kubectl get events -n erp-go --sort-by=.lastTimestamp | tail -20

# 健康检查
for svc in 8080 8081 8082 8085 8086 8087 8088 8091; do
  curl -s http://localhost:$svc/health | jq -c '{port, status} + .'
done
```

## 日志查看

```bash
# Gateway
kubectl logs -n erp-go deployment/api-gateway --tail=100 -f

# 特定服务（本地 dev）
Get-Content .cache/logs/order.log -Tail 50
Get-Content .cache/logs/warehouse.log -Tail 50
```

## 备份恢复

### PostgreSQL

```bash
# 备份
pg_dump $DATABASE_URL > backup_$(date +%Y%m%d).sql

# 恢复
psql $DATABASE_URL < backup_20260101.sql
```

### RabbitMQ

```bash
# 导出定义
rabbitmqadmin export rabbit.definitions.json

# 导入定义
rabbitmqadmin import rabbit.definitions.json
```

## 扩容

```bash
# Gateway
kubectl scale deployment/api-gateway -n erp-go --replicas=5

# 通过 Helm
helm upgrade erp-go ./docker/helm/erp-go -n erp-go --set gateway.replicas=5
```

## 死信处理

```bash
# 查看死信队列
curl http://localhost:8085/api/v1/order/outbox/dead-letter

# 重试
curl -X POST http://localhost:8085/api/v1/order/outbox/retry
```

## 告警指标

| 指标 | 阈值 | 检查方式 |
|------|------|----------|
| Pod 重启次数 | > 3 / 15min | `kubectl get pods` |
| 内存使用 | > 80% limit | `kubectl top pods` |
| /health 失败 | 任意 | 健康检查脚本 |
| 死信堆积 | > 50 | Outbox 查询 API |
| DB 连接池耗尽 | — | 应用日志 |

## 应急处理

| 现象 | 处理 |
|------|------|
| 服务不可用 | `kubectl rollout restart deployment/<svc> -n erp-go` |
| DB 连接失败 | 检查 Secret 中 DATABASE_URL，验证网络 |
| RabbitMQ 断连 | 检查 RabbitMQ Pod，应用自动重连 |
| 库存超卖 | 事务+行锁防并发；若发生，查 Inventory Journal 对账 |
