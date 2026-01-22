package queue

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type Config struct {
	Addresses string `json:"addresses" mapstructure:"addresses"`   // 服务地址集合
	Username  string `json:"username" mapstructure:"username"`     // 用户名
	Password  string `json:"password" mapstructure:"password"`     // 密码
	QueueName string `json:"queue-name" mapstructure:"queue-name"` //队列名
}

type ClientManager struct {
	mu      sync.RWMutex
	clients map[string]*RabbitMQ
}

func (cm *ClientManager) Add(name string, client *RabbitMQ) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[name] = client
}

func (cm *ClientManager) Get(name string) *RabbitMQ {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.clients[name]
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*RabbitMQ),
	}
}

var clientManager = NewClientManager()

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func Initialize() error {
	// 解析ES配置
	var cfgMap map[string]Config
	if err := viper.UnmarshalKey("rabbitmq", &cfgMap); err != nil {
		return err
	}
	// 数据库客户端链接添加到ClientManager
	for name, cfg := range cfgMap {
		// 实例化数据库链接组
		client, err := NewClient(cfg)
		if err != nil {
			return err
		}
		// 添加到GroupManager
		clientManager.Add(name, client)
	}
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewClient(cfg Config) (client *RabbitMQ, err error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s", cfg.Username, cfg.Password, cfg.Addresses))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		cfg.QueueName, // 队列名称
		false,         // 消息是否持久化
		false,         // 未使用队列是否自动删除
		false,         // 是否为排他队列
		false,         // 队列是否阻塞
		nil,           // 额外的属性
	)
	if err != nil {
		return nil, err
	}
	client = &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   q,
	}

	return
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

func (r *RabbitMQ) Publish(body string) {
	err := r.channel.Publish(
		"",           // 交换机名称
		r.queue.Name, // 路由键
		false,        // 如果无法路由到队列则返回消息
		false,        // 如果交换到队列的路由无匹配则丢弃消息
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Printf(" [x] Sent %s\n", body)
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.queue.Name, // 队列名称
		"",           // 消费者标签
		true,         // 自动确认消息
		false,        // 是否为独占消费者
		false,        // 是否等待服务器确认
		false,        // 是否阻塞
		nil,          // 额外的属性
	)
	failOnError(err, "Failed to register a consumer")
	return msgs, err
}
