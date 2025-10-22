package model

type Cart struct {
	CartID     uint `gorm:"primaryKey" json:"cart_id"`
	UserID     uint `gorm:"not null" json:"user_id"`
	ProductID  uint `gorm:"not null" json:"product_id"`
	Quantity   int  `gorm:"not null" json:"quantity"`
	IsSelected bool `gorm:"not null" json:"is_selected"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

type CartItem struct {
	ID           uint    `json:"id"`
	CartID       uint    `gorm:"not null"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	Quantity     int     `json:"quantity"`
	IsSelected   bool    `json:"is_selected"`

	//Foreign Key to Cart ID
	Cart Cart `gorm:"foreignKey:CartID" json:"cart"`
}
