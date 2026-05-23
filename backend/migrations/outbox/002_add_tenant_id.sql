-- 为 outbox_messages 表添加 tenant_id 列，支持多租户事件隔离
ALTER TABLE outbox_messages ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(64) NOT NULL DEFAULT '';
CREATE INDEX IF NOT EXISTS idx_outbox_tenant_id ON outbox_messages(tenant_id);
