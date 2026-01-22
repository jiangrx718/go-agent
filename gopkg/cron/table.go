package cron

import (
	rxCron "go-agent/gopkg/cron/base"
	"go-agent/gopkg/log"
)

type TableStatus struct {
}

func NewTableStatus() rxCron.Cron {
	return &TableStatus{}
}

func (ts *TableStatus) Spec() string {
	return "* * * * *"
}

func (ts *TableStatus) Run() {
	log.Sugar().Info("每分钟执行任务")
	// 执行处理业务逻辑
}
