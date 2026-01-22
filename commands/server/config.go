package server

import (
	"go-agent/gopkg/cache/es"
	"go-agent/gopkg/cron"
	"go-agent/gopkg/gorms"
	"go-agent/gopkg/log"
	"go-agent/gopkg/viper"
	"go-agent/internal/dao"

	"github.com/urfave/cli/v2"
)

func InitConfig(ctx *cli.Context) error {
	return InitConfigFromConfigPath(getConfigPath(ctx), getEnvPath(ctx))
}

func InitConfigFromConfigPath(configPath, envPath string) error {
	// 初始化配置文件
	if err := viper.Init(configPath, envPath); err != nil {
		return err
	}
	// 初始化日志
	if err := log.InitFromViper(); err != nil {
		return err
	}
	// 初始化orm
	if err := gorms.InitGenFromViper(dao.SetDefault); err != nil {
		return err
	}
	// 初始化Redis
	//if err := rxRedis.InitFromViper(); err != nil {
	//	return err
	//}
	//初始化cron定时任务
	if err := cron.DoCron(); err != nil {
		return err
	}
	// ES配置
	if err := es.Initialize(); err != nil {
		return err
	}
	return nil
}

func getConfigPath(ctx *cli.Context) string {
	if configFile := ctx.String("config"); configFile != "" {
		return configFile
	}

	return "config/config.yml"
}

func getEnvPath(ctx *cli.Context) string {
	if envPath := ctx.String("env"); envPath != "" {
		return envPath
	}

	return ""
}
