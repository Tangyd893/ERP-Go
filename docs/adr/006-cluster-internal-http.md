# ADR-006 集群内服务通信协议

日期：2026-05-25  
状态：采纳

## 背景

Gateway 到后端微服务的通信在 Kubernetes 集群内通过 Service DNS 进行。Sonar 热点 `kubernetes:S5332` 标记了明文 HTTP 协议的使用。

## 决策

集群内部服务间通信使用明文 HTTP（非 TLS），理由如下：

1. **信任边界**：服务仅通过 Kubernetes 内部 Service 网络通信，不直接暴露到公网。Kubernetes CNI 提供网络隔离，外部流量无法直接访问 Pod IP 或 ClusterIP Service。
2. **性能**：集群内 TLS 终止增加延迟，在微服务调用的高频场景下影响不明显但无安全收益。
3. **TLS 终止边界**：TLS 在 Ingress Controller / Gateway 入口统一终止，集群内部为可信域。

## 约束

- Gateway 入口必须启用 TLS（通过 Ingress 或外部负载均衡器）。
- 禁止将后端 Service 通过 NodePort / LoadBalancer 直接暴露到公网。
- 生产环境部署前须由安全团队评审网络策略与 Pod 安全上下文。

## 影响

- `docker/k8s/gateway-deployment.yaml` 中 `SERVICE_TARGET_*` 使用 `http://` 集群内 FQDN。
- Sonar 热点 `kubernetes:S5332` 标记为 Reviewed（已评审）。
