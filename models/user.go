package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email"`
	CartID   uint   `json:"cart_id" gorm:"foreignKey:CartID"`
	OrderID  uint   `json:"order_id" gorm:"foreignKey:OrderID"`
}
