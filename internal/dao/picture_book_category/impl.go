package picture_book_category

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
