-- Outbox 发件箱表（事务性消息发送）
CREATE TABLE IF NOT EXISTS outbox_messages (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id VARCHAR(64) NOT NULL,
    aggregate_type VARCHAR(64) NOT NULL,
    event_type VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    published_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_status ON outbox_messages(status);
CREATE INDEX IF NOT EXISTS idx_outbox_aggregate ON outbox_messages(aggregate_id);
CREATE INDEX IF NOT EXISTS idx_outbox_event_type ON outbox_messages(event_type);

-- Inbox 收件箱表（幂等消费）
CREATE TABLE IF NOT EXISTS inbox_messages (
    id BIGSERIAL PRIMARY KEY,
    message_id VARCHAR(128) NOT NULL UNIQUE,
    event_type VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_inbox_message_id ON inbox_messages(message_id);
