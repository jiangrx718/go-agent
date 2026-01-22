package redis

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var defaultName string
var clients map[string]redis.UniversalClient
var clientDefault map[string]*redis.Client

type Config struct {
	Type     string   `json:"type" mapstructure:"type"` // cluster, failover,single-node , default is single-node
	Addr     []string `json:"addr" mapstructure:"addr"`
	Password string   `json:"password" mapstructure:"password"`
	DB       int      `json:"db" mapstructure:"db"`
}

func InitFromViper() error {
	var cfg map[string]Config
	err := viper.UnmarshalKey("redis.client", &cfg)
	if err != nil {
		return err
	}
	clients = make(map[string]redis.UniversalClient)
	for k := range cfg {
		if clients[k], err = NewClient(cfg[k]); err != nil {
			return err
		}
	}

	hlog.Infof("redis client connect successful，环境地址: %s", strings.Join(cfg["account"].Addr, " "))

	return nil
}

func InitFromViperDefault() error {
	defaultName = viper.GetString("redis.default")
	var cfg map[string]Config
	err := viper.UnmarshalKey("redis.client", &cfg)
	if err != nil {
		return err
	}
	clientDefault = make(map[string]*redis.Client)
	for k := range cfg {
		if clientDefault[k], err = NewClientDefault(cfg[k]); err != nil {
			return err
		}
	}
	return nil
}
