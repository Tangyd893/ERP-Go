package outbox

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
)

const (
	defaultExchange    = "erp.events"
	defaultContentType = "application/json"
	reconnectDelay     = 3 * time.Second
)

// RabbitMQPublisher RabbitMQ 事件发布器，实现 EventPublisher 接口
type RabbitMQPublisher struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	url        string
	exchange   string
	log        logger.Logger
	mu         sync.Mutex
	closed     bool
	notifyConn chan *amqp.Error
	notifyChan chan *amqp.Error
}

// NewRabbitMQPublisher 创建 RabbitMQ 发布器
func NewRabbitMQPublisher(url, exchange string, log logger.Logger) (*RabbitMQPublisher, error) {
	if exchange == "" {
		exchange = defaultExchange
	}
	p := &RabbitMQPublisher{
		url:      url,
		exchange: exchange,
		log:      log,
	}
	if err := p.connect(); err != nil {
		return nil, err
	}
	go p.handleReconnect()
	return p, nil
}

func (p *RabbitMQPublisher) connect() error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return fmt.Errorf("RabbitMQ 连接失败: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("RabbitMQ 通道创建失败: %w", err)
	}
	if err := ch.ExchangeDeclare(p.exchange, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("RabbitMQ exchange 声明失败: %w", err)
	}
	p.mu.Lock()
	p.conn = conn
	p.channel = ch
	p.notifyConn = conn.NotifyClose(make(chan *amqp.Error))
	p.notifyChan = ch.NotifyClose(make(chan *amqp.Error))
	p.mu.Unlock()
	p.log.Infof("RabbitMQ 发布器已连接: %s, exchange=%s", p.url, p.exchange)
	return nil
}

func (p *RabbitMQPublisher) handleReconnect() {
	for {
		select {
		case err, ok := <-p.notifyConn:
			if !ok || err == nil {
				return
			}
		case err, ok := <-p.notifyChan:
			if !ok || err == nil {
				return
			}
		}
		p.mu.Lock()
		if p.closed {
			p.mu.Unlock()
			return
		}
		p.mu.Unlock()
		p.log.Warnf("RabbitMQ 连接断开，%d 秒后重连...", int(reconnectDelay.Seconds()))
		time.Sleep(reconnectDelay)
		for {
			var connErr error
			if connErr = p.connect(); connErr == nil {
				break
			}
			p.log.Errorf("RabbitMQ 重连失败: %v", connErr)
			time.Sleep(reconnectDelay)
		}
	}
}

// Publish 发布事件到 RabbitMQ exchange
func (p *RabbitMQPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed || p.channel == nil {
		return fmt.Errorf("RabbitMQ 发布器已关闭")
	}
	return p.channel.PublishWithContext(ctx, p.exchange, eventType, false, false, amqp.Publishing{
		ContentType:  defaultContentType,
		DeliveryMode: amqp.Persistent,
		Body:         payload,
		Timestamp:    time.Now(),
	})
}

// Close 关闭发布器连接
func (p *RabbitMQPublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// 确保 RabbitMQPublisher 实现 EventPublisher 接口
var _ EventPublisher = (*RabbitMQPublisher)(nil)
