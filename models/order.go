package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	OrderId uint       `json:"order_id" gorm:"primary_key" gorm:"column:order_id;"`
	UserId  uint       `json:"user_id" gorm:"column:user_id;"`
	CartId  uint       `json:"cart_id" gorm:"column:cart_id;"`
	Total   uint       `json:"total" gorm:"column:total;"`
	Time    *time.Time `json:"time" gorm:"type:datetime(3)" gorm:"column:time;"`
	Status  string     `json:"status" gorm:"default:unpaid" gorm:"column:status;"`
}

// BeforeCreate 钩子函数，在创建订单时设置创建时间 在reate函数的时候自动调用
func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now() // 获取当前时间
	o.Time = &now     // 设置时间为当前时间
	return nil
}
