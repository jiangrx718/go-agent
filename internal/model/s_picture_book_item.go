package model

import (
	"time"
)

// 绘本详情
type SPictureBookItem struct {
	Id        uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BookId    string    `gorm:"column:book_id;type:char(32);default:'';comment:绘本id;NOT NULL" json:"book_id"`
	Pic       string    `gorm:"column:pic;type:varchar(1024);default:'';comment:绘本详情图;NOT NULL" json:"pic"`
	BPic      string    `gorm:"column:b_pic;type:varchar(1024);default:'';comment:绘本详情大图;NOT NULL" json:"b_pic"`
	Audio     string    `gorm:"column:audio;type:varchar(1024);comment:绘本详音频;NOT NULL" json:"audio"`
	Position  int       `gorm:"column:position;type:int(11);default:0;comment:排序位置;NOT NULL" json:"position"`
	Content   string    `gorm:"column:content;type:varchar(1024);default:'';comment:内容;NOT NULL" json:"content"`
	Status    string    `gorm:"column:status;type:varchar(20);default:on;comment:状态,on启用,off禁用;NOT NULL" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:添加时间;NOT NULL" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}

func (m *SPictureBookItem) TableName() string {
	return "s_picture_book_item"
}
