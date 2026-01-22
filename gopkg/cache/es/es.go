package es

import (
	"errors"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
)

type Config struct {
	Addresses []string `json:"addresses" mapstructure:"addresses"` // 服务地址集合
	Username  string   `json:"username" mapstructure:"username"`   // 用户名
	Password  string   `json:"password" mapstructure:"password"`   // 密码
}

// Initialize 根据配置文件加载ES实例
func Initialize() error {
	// 解析ES配置
	var cfgMap map[string]Config
	if err := viper.UnmarshalKey("elasticsearch", &cfgMap); err != nil {
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

// ClientManager 管理数据库客户端Client
type ClientManager struct {
	mu      sync.RWMutex
	clients map[string]*elasticsearch.Client
}

// Add 添加Client到Manager
func (cm *ClientManager) Add(name string, client *elasticsearch.Client) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.clients[name] = client
}

// Get 从Manage获取Client
func (cm *ClientManager) Get(name string) *elasticsearch.Client {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.clients[name]
}

// NewClientManager 实例化ClientManager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*elasticsearch.Client),
	}
}

// clientManager
var clientManager = NewClientManager()

// Get 从ClientManager获取Client
func Get(name string) (*elasticsearch.Client, error) {
	client := clientManager.Get(name)
	if client == nil {
		return client, errors.New("client: [" + name + "] not exists")
	}
	return client, nil
}

// NewClient 创建es客户端连接
func NewClient(cfg Config) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
}
