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
