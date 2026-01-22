package gorms

import (
	"errors"
	"gorm.io/gorm"
)

type BaseDao struct {
}

func NewBaseDao() *BaseDao {
	return &BaseDao{}
}

func (b *BaseDao) ConvertError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}
