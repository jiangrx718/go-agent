package picture_book

import (
	"go-agent/gopkg/gorms"
)

type Dao struct {
	*gorms.BaseDao
}

func NewDao() *Dao {
	return &Dao{
		BaseDao: gorms.NewBaseDao(),
	}
}
