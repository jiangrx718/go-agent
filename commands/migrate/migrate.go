package migrate

import (
	"go-agent/gopkg/gorms"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "数据库迁移",
		Subcommands: []*cli.Command{
			{
				Name:        "up",
				Usage:       "自动迁移数据库",
				Description: "自动迁移数据库",
				Action: func(ctx *cli.Context) error {
					tx := gorms.Client()
					tx.DisableForeignKeyConstraintWhenMigrating = true
					tables := []any{}
					return tx.AutoMigrate(tables...)
				},
			},
		},
	}
}
