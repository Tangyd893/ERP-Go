# 部署手册

## 前置条件

- Kubernetes 1.25+ 集群
- Helm 3.x
- kubectl 配置指向目标集群
- PostgreSQL 16 + RabbitMQ 3.12 + Redis 7 + MinIO（可用外部托管或集群内部署）

## 快速部署（Helm）

```bash
# 1. 创建命名空间
kubectl create namespace erp-go

# 2. 配置 Secret
kubectl create secret generic erp-go-secrets -n erp-go \
  --from-literal=database-url="postgres://erp:${DB_PASSWORD}@postgres:5432/erp_go?sslmode=disable" \
  --from-literal=rabbitmq-url="amqp://erp:${RMQ_PASSWORD}@rabbitmq:5672/" \
  --from-literal=jwt-secret="${JWT_SECRET}" \
  --from-literal=sonar-token="${SONAR_TOKEN}"

# 3. 安装
helm install erp-go ./docker/helm/erp-go \
  --namespace erp-go \
  --set image.tag=v0.9.0 \
  --set gateway.replicas=3 \
  --set ingress.host=erp-go.example.com

# 4. 验证
kubectl get pods -n erp-go
kubectl get svc -n erp-go
```

## 手动部署（kubectl）

```bash
# 部署 Gateway
kubectl apply -f docker/k8s/gateway-rbac.yaml
kubectl apply -f docker/k8s/gateway-deployment.yaml

# 检查
kubectl rollout status deployment/api-gateway -n erp-go
curl http://api-gateway.erp-go.svc.cluster.local:8080/health
```

## 回滚

```bash
# Helm 回滚
helm rollback erp-go -n erp-go

# kubectl 回滚
kubectl rollout undo deployment/api-gateway -n erp-go
```

## 环境变量

| 变量 | 说明 | 默认 |
|------|------|------|
| `ENVIRONMENT` | 环境标识 | `production` |
| `DATABASE_URL` | PostgreSQL 连接串 | — |
| `RABBITMQ_URL` | RabbitMQ 连接串 | — |
| `SERVICE_TARGET_*` | 服务间调用地址 | `http://<svc>:<port>` |

## 健康检查

所有服务暴露 `/health` 端点：

```bash
curl http://localhost:8080/health  # {"status":"ok","service":"gateway","db":true}
curl http://localhost:8081/health  # iam-service
curl http://localhost:8085/health  # order-service
curl http://localhost:8086/health  # inventory-service
curl http://localhost:8087/health  # warehouse-service
```
