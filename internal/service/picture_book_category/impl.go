package picture_book_category

import (
	"go-agent/internal/dao"
	"go-agent/internal/dao/picture_book_category"
)

type Service struct {
	pictureBookCategoryDao dao.PictureBookCategory
}

func NewService() *Service {
	return &Service{
		pictureBookCategoryDao: picture_book_category.NewDao(),
	}
}
