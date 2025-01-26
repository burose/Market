package models

type Cart struct {
	CartId uint `json:"cart_id" gorm:"column:cart_id;primary_key;"`
	UserId uint `json:"user_id" gorm:"column:user_id;"`

	ProductID uint `json:"product_id" gorm:"column:product_id"`
	Quantity  uint `json:"quantity" gorm:"column:quantity"`
	Price     uint `json:"price" gorm:"column:price"`
}
