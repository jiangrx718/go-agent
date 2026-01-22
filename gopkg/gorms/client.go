package gorms

import (
	"time"

	"github.com/spf13/viper"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitFormViper() (*gorm.DB, error) {
	viper.SetDefault("DB_DSN", "")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 20)

	var err error
	if db, err = NewDatabase(viper.GetString("db.dsn"), &gorm.Config{}); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(viper.GetInt("DB_MAX_OPEN_CONNS"))
	sqlDB.SetMaxIdleConns(viper.GetInt("DB_MAX_IDLE_CONNS"))
	sqlDB.SetConnMaxIdleTime(time.Second * 5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, err
}

func InitGenFromViper(setDefault func(db *gorm.DB, opts ...gen.DOOption)) error {
	var err error
	if db, err = InitFormViper(); err != nil {
		return err
	}
	if viper.GetBool("debug") {
		db = db.Debug()
	}
	if viper.GetBool("local") {
		db.Config.Logger = logger.Default.LogMode(logger.Silent)
	}
	setDefault(db)
	return nil
}

func Client() *gorm.DB {
	return db
}
