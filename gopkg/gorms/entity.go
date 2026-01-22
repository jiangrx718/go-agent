package gorms

import (
	"gorm.io/gorm"
)

type Entity struct {
	ID        uint           `json:"-" gorm:"column:id;type:uint;not null;primaryKey;autoIncrement;comment:主键;"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间"`
}
