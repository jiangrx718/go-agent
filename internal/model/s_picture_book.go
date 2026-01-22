package model

import (
	"time"
)

// 绘本表
type SPictureBook struct {
	Id         uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	BookId     string    `gorm:"column:book_id;type:char(32);default:'';comment:绘本id;NOT NULL" json:"book_id"`
	Title      string    `gorm:"column:title;type:varchar(1024);default:'';comment:绘本标题;NOT NULL" json:"title"`
	Icon       string    `gorm:"column:icon;type:varchar(1024);default:'';comment:绘本封面;NOT NULL" json:"icon"`
	CategoryId string    `gorm:"column:category_id;type:char(36);default:'';comment:绘本所属栏目;NOT NULL" json:"category_id"`
	Type       int       `gorm:"column:type;type:int(11);default:0;comment:1中文绘本,2英文绘本,3古诗绘本,4英语词汇;NOT NULL" json:"type"`
	Status     string    `gorm:"column:status;type:varchar(20);default:on;comment:状态,on启用,off禁用;NOT NULL" json:"status"`
	Position   int       `gorm:"column:position;type:int(11);default:0;comment:排序位置;NOT NULL" json:"position"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:添加时间;NOT NULL" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}

func (m *SPictureBook) TableName() string {
	return "s_picture_book"
}
