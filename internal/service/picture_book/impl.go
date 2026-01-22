package picture_book

import (
	"go-agent/internal/dao"
	"go-agent/internal/dao/picture_book"
)

type Service struct {
	pictureBookDao dao.PictureBook
}

func NewService() *Service {
	return &Service{
		pictureBookDao: picture_book.NewDao(),
	}
}
