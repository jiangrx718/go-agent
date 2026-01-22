package commands

import (
	"go-agent/commands/generate"
	"go-agent/commands/gorm"
	"go-agent/commands/migrate"
	"go-agent/commands/worker"

	"github.com/urfave/cli/v2"
)

func All() []*cli.Command {
	commands := []*cli.Command{
		migrate.Command(),
		generate.Command(),
		gorm.Command(),
		worker.Command(),
	}
	return commands
}
