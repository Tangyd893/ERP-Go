# 数据库与消息队列备份恢复 Runbook

## 备份策略

| 组件 | 频率 | 保留 | 工具 |
|---|---|---|---|
| PostgreSQL | 每日 02:00 | 7 天 | `pg_dump` / `pg_basebackup` |
| RabbitMQ | 每日 03:00 | 7 天 | 定义导出 + 消息持久化 |
| MinIO | 每日 04:00 | 30 天 | `mc mirror` |

## 1. PostgreSQL 备份

### 1.1 创建备份

```bash
# 全量逻辑备份
pg_dump -h $PGHOST -p $PGPORT -U $PGUSER -d erp_go \
  --format=custom \
  --file=/backups/erp_go_$(date +%Y%m%d_%H%M).dump

# 仅结构备份
pg_dump -h $PGHOST -p $PGPORT -U $PGUSER -d erp_go \
  --schema-only \
  --file=/backups/erp_go_schema_$(date +%Y%m%d).sql
```

### 1.2 Kubernetes CronJob 示例

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: pg-backup
  namespace: erp-go
spec:
  schedule: "0 2 * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: backup
              image: postgres:16-alpine
              command:
                - /bin/sh
                - -c
                - |
                  pg_dump -h $PGHOST -U $PGUSER -d erp_go \
                    --format=custom -f /backups/erp_go_$(date +%Y%m%d_%H%M).dump
              env:
                - name: PGHOST
                  value: "postgres-service"
                - name: PGUSER
                  valueFrom:
                    secretKeyRef:
                      name: erp-go-secrets
                      key: db-user
                - name: PGPASSWORD
                  valueFrom:
                    secretKeyRef:
                      name: erp-go-secrets
                      key: db-password
              volumeMounts:
                - name: backup-storage
                  mountPath: /backups
          volumes:
            - name: backup-storage
              persistentVolumeClaim:
                claimName: backup-pvc
          restartPolicy: OnFailure
```

### 1.3 恢复

```bash
# 全量恢复
pg_restore -h $PGHOST -p $PGPORT -U $PGUSER -d erp_go \
  --clean --if-exists \
  /backups/erp_go_20260525_0200.dump

# 恢复后立即跑迁移（确保 schema 版本一致）
./scripts/migrate.sh
```

## 2. RabbitMQ 备份

### 2.1 导出定义

```bash
# 导出所有定义（exchanges, queues, bindings）
rabbitmqadmin export /backups/rabbit_definitions_$(date +%Y%m%d).json

# 或通过 HTTP API
curl -u admin:admin http://localhost:15672/api/definitions \
  -o /backups/rabbit_definitions.json
```

### 2.2 恢复

```bash
rabbitmqadmin import /backups/rabbit_definitions_20260525.json
```

### 2.3 消息持久化验证

确认所有关键队列已设置为 `durable: true`：
- `order.fulfillment`
- 各业务 dead-letter-queue

## 3. MinIO 备份

```bash
# 镜像整个 bucket 到本地
mc mirror minio/erp-go /backups/minio/erp-go/

# 恢复
mc mirror /backups/minio/erp-go/ minio/erp-go/
```

## 4. 定期演练检查项

| 项目 | 验证方法 | 期望 |
|---|---|---|
| pg_dump 可执行 | `docker exec erp-postgres pg_dump --version` | 版本号输出 |
| restore 可恢复 | 在 staging 环境执行 pg_restore | 无错误，`\dt` 可见表 |
| RabbitMQ 定义可导入 | `rabbitmqadmin import` 到本地环境 | 队列/exchange 恢复 |
| MinIO 文件可读 | `mc ls minio/erp-go/` | 文件列表正确 |
