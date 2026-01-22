package cron

import (
	rxCron "go-agent/gopkg/cron/base"

	"github.com/spf13/viper"
)

func DoCron() error {
	if !viper.GetBool("cron.switch") {
		return nil
	}

	tableStatusList := make([]rxCron.Cron, 0)
	tableStatusList = append(tableStatusList, NewTableStatus())
	if err := rxCron.InitFromMinute(tableStatusList); err != nil {
		return err
	}
	return nil
}
