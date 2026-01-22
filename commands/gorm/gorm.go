package gorm

import (
	"go-agent/gopkg/gorms"
	"go-agent/gopkg/log"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

// go run main.go mysql insert
func Command() *cli.Command {
	return &cli.Command{
		Name:  "mysql",
		Usage: "mysql数据库迁移",
		Subcommands: []*cli.Command{
			{
				Name:        "insert",
				Usage:       "导入MySQL数据",
				Description: "导入MySQL数据",
				Action: func(ctx *cli.Context) error {
					logPrefix := "mysql-insert"

					pwd, _ := os.Getwd()
					err := filepath.Walk(pwd+"/data/mysql", func(path string, info os.FileInfo, err error) error {
						if info.IsDir() == false {
							if strings.Contains(path, "mysql") {
								file, err := os.Open(path)
								if err != nil {
									log.Sugar().Error(ctx.Context, logPrefix, zap.Any("os.Open file error", err))
									return err
								}

								var delim byte = '\n'
								if strings.Contains(path, "document_parsing_result") {
									delim = ';'
								}

								reader := bufio.NewReader(file)
								for {
									sql, err := reader.ReadString(delim)
									if err != nil {
										if errors.Is(err, io.EOF) {
											break
										}
										log.Sugar().Error(ctx.Context, logPrefix, zap.Any("reader.ReadString file error", err), zap.Any("path", path))
										return err
									}

									sql = strings.TrimSpace(sql)
									if len(sql) < 1 {
										continue
									}

									tx := gorms.Client().Exec(sql)
									if tx.Error != nil {
										log.Sugar().Error(ctx.Context, logPrefix, zap.Any("gorms.Client().Exec error", tx.Error), zap.Any("path", path))
										return tx.Error
									}
								}
								fmt.Println("insert completed ", path)
							}
						}

						return nil
					})
					if err != nil {
						log.Sugar().Error(ctx.Context, logPrefix, zap.Any("mysql-insert error", err))
						return err
					}

					return nil
				},
			},
		},
	}
}
