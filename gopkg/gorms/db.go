package gorms

import (
	"strings"
	"sync"

	"github.com/pkg/errors"

	gaussdb "github.com/okyer/gorm4gaussdb"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	databaseOnce   sync.Once
	databaseClient *gorm.DB
)

func NewDatabase(dsn string, config *gorm.Config) (*gorm.DB, error) {
	var err error
	var dialector gorm.Dialector

	if config == nil {
		config = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Error),
			DisableForeignKeyConstraintWhenMigrating: true, // 在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为true
			// SkipDefaultTransaction:                   true, // 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它
			// PrepareStmt:                              true, // 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
		}
	}

	dbDialect := viper.GetString("db.dialect")
	switch strings.ToLower(dbDialect) {
	case "gaussdb":
		dsn = strings.Trim(dsn, "'")
		dialector = gaussdb.Open(dsn)
	case "mysql":
		dialector = mysql.Open(dsn)
	case "postgres":
		dialector = postgres.Open(dsn)
	default:
		return nil, errors.Errorf("db.dialect not found: %s", dbDialect)
	}

	databaseOnce.Do(func() {
		databaseClient, err = gorm.Open(dialector, config)
	})

	return databaseClient, err
}
