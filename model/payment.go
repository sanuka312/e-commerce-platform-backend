package model

type Payment struct {
	PaymentId     uint    `gorm:"PrimaryKey" json:"payment_id"`
	OrderId       uint    `gorm:"not null" json:"order_id"`
	UserId        string  `gorm:"type:uuid;not null" json:"user_id"`
	PaymentMethod string  `gorm:"size:100;not null" json:"payment_method"`
	PaymentAmount float64 `gorm:"type:decimal(10,2);not null" json:"payment_amount"`
	Status        string  `gorm:"size:50;not null;default:'pending'" json:"status"`

	//Relationships
	Order Order `gorm:"foreignKey:OrderId" json:"order"`
	User  User  `gorm:"foreignKey:UserId;references:KeyCloakUserId" json:"user"`
}
