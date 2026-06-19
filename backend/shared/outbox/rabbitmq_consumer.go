package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
)

const (
	retryExchange   = "erp.events.retry"
	deadLetterQueue = "erp.dlq"
	maxRetries      = 3
)

// ConsumeHandler 消费处理回调，返回 error 时消息进入重试
type ConsumeHandler func(ctx context.Context, eventType string, messageID string, payload []byte) error

// RabbitMQConsumer RabbitMQ 消费者
type RabbitMQConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	url       string
	queueName string
	bindKeys  []string
	handler   ConsumeHandler
	inbox     InboxStore
	log       logger.Logger
	prefetch  int
	mu        sync.Mutex
	closed    bool
	consumerTag string
	cancel     context.CancelFunc
}

// NewRabbitMQConsumer 创建 RabbitMQ 消费者
func NewRabbitMQConsumer(ctx context.Context, url, queueName string, bindKeys []string, handler ConsumeHandler, inbox InboxStore, log logger.Logger) (*RabbitMQConsumer, error) {
	cancelCtx, cancel := context.WithCancel(ctx)
	_ = cancelCtx // kept for graceful shutdown via cancel
	c := &RabbitMQConsumer{
		url:       url,
		queueName: queueName,
		bindKeys:  bindKeys,
		handler:   handler,
		inbox:     inbox,
		log:       log,
		prefetch:  10,
		cancel:    cancel,
	}
	if err := c.connect(); err != nil {
		cancel()
		return nil, err
	}
	return c, nil
}

func (c *RabbitMQConsumer) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("RabbitMQ 消费者连接失败: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("RabbitMQ 消费者通道创建失败: %w", err)
	}
	if err := ch.Qos(c.prefetch, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("QoS 设置失败: %w", err)
	}

	args := amqp.Table{
		"x-dead-letter-exchange": retryExchange,
	}
	queue, err := ch.QueueDeclare(c.queueName, true, false, false, false, args)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("队列声明失败: %w", err)
	}
	queueName := queue.Name

	if err := ch.ExchangeDeclare(retryExchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("重试 exchange 声明失败: %w", err)
	}

	deadLetterArgs := amqp.Table{
		"x-message-ttl": int32(60000),
	}
	_, err = ch.QueueDeclare(deadLetterQueue, true, false, false, false, deadLetterArgs)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("死信队列声明失败: %w", err)
	}
	if err := ch.QueueBind(deadLetterQueue, "#", retryExchange, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("死信队列绑定失败: %w", err)
	}

	for _, key := range c.bindKeys {
		if err := ch.QueueBind(queueName, key, defaultExchange, false, nil); err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("队列绑定失败 %s: %w", key, err)
		}
		c.log.Infof("绑定队列 %s -> exchange=%s, key=%s", queueName, defaultExchange, key)
	}

	c.mu.Lock()
	c.conn = conn
	c.channel = ch
	c.mu.Unlock()

	c.consumerTag = fmt.Sprintf("%s-consumer", queueName)
	deliveries, err := ch.Consume(queueName, c.consumerTag, false, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("开始消费失败: %w", err)
	}

	go c.processDeliveries(deliveries)
	c.log.Infof("RabbitMQ 消费者已启动: queue=%s", queueName)
	return nil
}

func (c *RabbitMQConsumer) processDeliveries(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		c.handleDelivery(&d)
	}
}

func (c *RabbitMQConsumer) handleDelivery(d *amqp.Delivery) {
	eventType := d.RoutingKey
	messageID := d.MessageId
	if messageID == "" {
		messageID = fmt.Sprintf("%s-%s", eventType, d.Timestamp.String())
	}

	duplicate, err := c.inbox.IsDuplicate(context.Background(), messageID)
	if err != nil {
		c.log.Errorf("Inbox 幂等检查失败: %v", err)
		d.Nack(false, true)
		return
	}
	if duplicate {
		d.Ack(false)
		return
	}

	if err := c.handler(context.Background(), eventType, messageID, d.Body); err != nil {
		c.log.Errorf("事件处理失败: type=%s, err=%v", eventType, err)
		retryCount := c.getRetryCount(d)
		if retryCount >= maxRetries {
			d.Nack(false, false)
			return
		}
		d.Nack(false, true)
		return
	}

	// 处理器可能已写入 Inbox（如 P4 协调器），避免重复保存
	if dup, _ := c.inbox.IsDuplicate(context.Background(), messageID); dup {
		d.Ack(false)
		return
	}

	inboxMsg := &InboxMessage{
		MessageID:   messageID,
		EventType:   eventType,
		Payload:     d.Body,
		ProcessedAt: d.Timestamp,
	}
	if err := c.inbox.Save(context.Background(), inboxMsg); err != nil {
		c.log.Errorf("Inbox 保存失败: %v", err)
		d.Nack(false, true)
		return
	}
	d.Ack(false)
}

func (c *RabbitMQConsumer) getRetryCount(d *amqp.Delivery) int {
	if d.Headers == nil {
		return 0
	}
	countRaw, ok := d.Headers["x-retry-count"]
	if !ok {
		return 0
	}
	switch v := countRaw.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		n, _ := v.Int64()
		return int(n)
	}
	return 0
}

// Close 关闭消费者
func (c *RabbitMQConsumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	c.cancel()
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
