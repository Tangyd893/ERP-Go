-- ============================================================
-- Notification 通知服务 - 数据库初始化
-- ============================================================
CREATE SCHEMA IF NOT EXISTS notification;
SET search_path TO notification, public;

CREATE TABLE IF NOT EXISTS notifications (
    id         VARCHAR(36) PRIMARY KEY,
    tenant_id  VARCHAR(36) NOT NULL,
    user_id    VARCHAR(36) NOT NULL,
    title      VARCHAR(255) NOT NULL,
    content    TEXT DEFAULT '',
    type       VARCHAR(16) NOT NULL DEFAULT 'info',
    read       BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notif_user ON notifications(user_id, read);
CREATE INDEX idx_notif_created_at ON notifications(created_at DESC);
