// Package server 启动 report-service（无状态聚合服务）
//
// report-service 不持有数据库，不创建业务表。
// 当前模式：
//   - DEMO_MODE=true 或 ENVIRONMENT 为空 → 返回示例数据
//   - 生产环境（DEMO_MODE != true）→ 返回空数据集
//
// 后续计划（T-501+）：
//   通过跨服务 API 调用聚合各业务服务数据，实现真实报表。
//   依赖：各业务服务需暴露内部查询 API 或共享数据库读取权限。
package main
