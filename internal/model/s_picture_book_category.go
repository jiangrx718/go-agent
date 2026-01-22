package model

import (
	"time"
)

// 绘本类型表
type SPictureBookCategory struct {
	Id           uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键id" json:"id"`
	CategoryId   string    `gorm:"column:category_id;type:char(36);default:'';comment:栏目ID;NOT NULL" json:"category_id"`
	CategoryName string    `gorm:"column:category_name;type:varchar(256);default:'';comment:栏目名称;NOT NULL" json:"category_name"`
	Sort         int       `gorm:"column:sort;type:int(11);default:0;comment:排序,倒序;NOT NULL" json:"sort"`
	Type         int       `gorm:"column:type;type:int(11);default:0;comment:1中文绘本,2英文绘本,3古诗绘本,4英语词汇;NOT NULL" json:"type"`
	Status       string    `gorm:"column:status;type:varchar(20);default:on;comment:状态,on启用,off禁用;NOT NULL" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:添加时间;NOT NULL" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"updated_at"`
}

func (m *SPictureBookCategory) TableName() string {
	return "s_picture_book_category"
}
