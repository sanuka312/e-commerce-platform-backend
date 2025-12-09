package model

type Cart struct {
	CartID         uint   `gorm:"primaryKey" json:"cart_id"`
	KeycloakUserID string `gorm:"not null;index" json:"keycloak_user_id"`

	// Relationships

	Items []CartItem `gorm:"foreignKey:CartID" json:"cart_items"`
}

type CartItem struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	CartID    uint `gorm:"not null"`
	ProductID uint `json:"product_id"`

	UnitPrice  float64 `json:"unit_price"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	IsSelected bool    `json:"is_selected"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
