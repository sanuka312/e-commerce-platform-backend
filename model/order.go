package model

import "time"

type Order struct {
	OrderId        uint   `gorm:"PrimaryKey" json:"order_id"`
	KeycloakUserID string `gorm:"not null;index" json:"keycloak_user_id"`
	ProductId      uint   `gorm:"not null"`
	PaymentId      uint   `gorm:"not null"`

	ProductPrice float64 `json:"product_price"`
	Quantity     uint    `gorm:"not null" json:"qty"`

	TotalPrice  float64   `json:"total_price"`
	AddressId   *uint     `json:"address_id"`
	OrderStatus string    `gorm:"size:50;default:'pending'" json:"order_status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	//Relationships
	Product Product  `gorm:"foreignKey:ProductId" json:"product"`
	Address *Address `gorm:"foreignKey:AddressId" json:"address"`
	Payment Payment  `gorm:"foreignKey:PaymentId" json:"payment"`
	Price   Product  `gorm:"foreignKey:ProductPrice" json:"price"`
}
