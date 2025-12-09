package model

type Payment struct {
	PaymentId      uint    `gorm:"PrimaryKey" json:"payment_id"`
	OrderId        *uint   `gorm:"index" json:"order_id"`
	KeycloakUserID string  `gorm:"not null;index" json:"keycloak_user_id"`
	PaymentMethod  string  `gorm:"size:100;not null" json:"payment_method"`
	PaymentAmount  float64 `gorm:"type:decimal(10,2);not null" json:"payment_amount"`
	Status         string  `gorm:"size:50;not null;default:'pending'" json:"status"`
}
