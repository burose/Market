package models

type Product struct {
	ProductID   uint   `json:"product_id" gorm:"column:product_id;primary_key;AUTO_INCREMENT"`
	Name        string `json:"name" gorm:"column:name;"`
	Price       uint   `json:"price" gorm:"column:price;"`
	Picture     string `json:"picture" gorm:"column:picture;"`
	Description string `json:"description" gorm:"column:description;"`
	Number      uint   `json:"number" gorm:"column:number;"`
}
