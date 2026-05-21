-- ============================================================
-- File 文件服务 - 数据库初始化
-- ============================================================

CREATE SCHEMA IF NOT EXISTS file;
SET search_path TO file, public;

CREATE TABLE IF NOT EXISTS files (
    id          VARCHAR(36) PRIMARY KEY,
    tenant_id   VARCHAR(36)  NOT NULL,
    bucket      VARCHAR(64)  NOT NULL,
    object_key  VARCHAR(255) NOT NULL,
    file_name   VARCHAR(255) NOT NULL,
    file_size   BIGINT       DEFAULT 0,
    mime_type   VARCHAR(128) DEFAULT '',
    source_type VARCHAR(32)  DEFAULT '',
    source_id   VARCHAR(36)  DEFAULT '',
    created_by  VARCHAR(36)  DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_files_tenant_id ON files(tenant_id);
CREATE INDEX idx_files_source ON files(source_type, source_id);
