package generate

import (
	"go-agent/gopkg/gorms"
	"go-agent/internal/model"

	"github.com/urfave/cli/v2"
	"gorm.io/gen"
)

func Command() *cli.Command {
	return &cli.Command{
		Name: "generate",
		Action: func(ctx *cli.Context) error {
			g := gen.NewGenerator(gen.Config{
				OutPath: "internal/dao",
				Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
			})
			g.UseDB(gorms.Client())
			g.ApplyBasic(
				model.SPictureBook{},
				model.SPictureBookItem{},
				model.SPictureBookCategory{},
			)
			g.Execute()
			return nil
		},
	}
}
