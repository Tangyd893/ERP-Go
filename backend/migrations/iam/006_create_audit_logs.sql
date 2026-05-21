-- 操作审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL DEFAULT '',
    user_id VARCHAR(64) DEFAULT '',
    username VARCHAR(128) DEFAULT '',
    action VARCHAR(64) NOT NULL,
    resource_type VARCHAR(64) DEFAULT '',
    resource_id VARCHAR(64) DEFAULT '',
    detail TEXT DEFAULT '',
    ip VARCHAR(64) DEFAULT '',
    user_agent TEXT DEFAULT '',
    request_id VARCHAR(64) DEFAULT '',
    trace_id VARCHAR(64) DEFAULT '',
    result VARCHAR(16) DEFAULT '',
    result_detail TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
