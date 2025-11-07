package model

type Cart struct {
	CartID uint `gorm:"primaryKey" json:"cart_id"`
	UserID uint `gorm:"not null" json:"user_id"` // Internal numeric user_id

	// Relationships
	User  User       `gorm:"foreignKey:UserID;references:UserId" json:"user"`
	Items []CartItem `gorm:"foreignKey:CartID" json:"cart_items"`
}

type CartItem struct {
	ID          uint    `json:"id"`
	CartID      uint    `gorm:"not null"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	IsSelected  bool    `json:"is_selected"`

	//Foreign Key to Cart ID
	Cart    Cart    `gorm:"foreignKey:CartID" json:"cart"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
