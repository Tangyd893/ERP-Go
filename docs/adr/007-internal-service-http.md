# ADR-007: 集群内服务间使用 HTTP 明文通信

## 状态

已接受（Accepted）

## 日期

2026-06-05

## 背景

Kubernetes 集群内各微服务（gateway → iam / order / inventory / warehouse）之间通过 Service DNS 名称互相调用。当前服务间调用使用 `http://` 协议，SonarCloud 安全热点规则 `kubernetes:S5332` 标记了此配置。

## 决策

**集群内 Service 间通信继续使用 HTTP 明文。**

## 理由

1. **网络边界明确**：所有服务间流量均发生在 Kubernetes 集群内部网络（CNI 隔离），不经过公网。
2. **无外部可达性**：`*.svc.cluster.local` DNS 仅集群内可解析，外部无法访问。
3. **Sidecar / Service Mesh 可后续引入**：如需加密，可在不修改应用代码的情况下通过 Istio/Linkerd 注入 mTLS sidecar。
4. **运维复杂度可控**：引入内部 TLS 需管理集群内证书轮换，当前阶段 ROI 不高。

## 后果

- SonarCloud 该热点标记为 **Safe（已审查）**，不再计入 QG 阻断。
- 若未来集群网络策略变化或需要合规加密，应引入 Service Mesh mTLS 而非逐服务配置 TLS。

## 相关

- `docker/k8s/gateway-deployment.yaml` — `SERVICE_TARGET_IAM` 环境变量
- SonarCloud Hotspot: `kubernetes:S5332`
