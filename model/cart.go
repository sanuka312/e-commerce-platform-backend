package model

type Cart struct {
	CartID uint   `gorm:"primaryKey" json:"cart_id"`
	UserID string `gorm:"type:uuid;not null" json:"user_id"`

	User User `gorm:"foreignKey:UserID;references:KeyCloakUserId" json:"user"`
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
	Cart Cart `gorm:"foreignKey:CartID" json:"cart"`
}
