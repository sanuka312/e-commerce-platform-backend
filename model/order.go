package model

type Order struct {
	OrderId     uint   `gorm:"PrimaryKey" json:"order_id"`
	UserId      uint   `gorm:"not null" json:"user_id"`
	ProductId   uint   `gorm:"not null"`
	Quantity    uint   `gorm:"not null" json:"qty"`
	AddressId   *uint  `json:"address_id"`
	OrderStatus string `gorm:"size:50;default:'pending'" json:"order_status"`

	//Relationships
	Product Product  `gorm:"foreignKey:ProductId" json:"product"`
	Address *Address `gorm:"foreignKey:AddressId" json:"address"`
	User    *User    `gorm:"foreignKey:UserId"`
}
