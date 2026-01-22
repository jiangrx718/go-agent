package main

import (
	"go-agent/commands"
	"go-agent/commands/server"
	"go-agent/gopkg/log"
	"go-agent/internal/worker"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// 创建后台工作进程
	backgroundWorker := worker.NewBackgroundWorker(10, 100)

	// 在程序退出时停止后台工作进程
	defer backgroundWorker.Stop()

	// 启动后台工作进程
	backgroundWorker.Start()

	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		// 启动 account 服务
		return server.Run(c)
	}
	app.Before = server.InitConfig
	app.After = server.Flush
	app.Commands = commands.All()
	app.Flags = server.Flags()

	if err := app.Run(os.Args); err != nil {
		log.Sugar().Fatal(err)
	}
}
