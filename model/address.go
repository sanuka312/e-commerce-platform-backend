package model

type Address struct {
	AddressId  uint   `gorm:"primaryKey;autoIncrement" json:"address_id"`
	UserId     uint   `gorm:"type:uuid;not null" json:"user_id"`
	Line1      string `gorm:"size:200;not null" json:"line1"`
	Line2      string `gorm:"size:200;not null" json:"line2"`
	City       string `gorm:"size:100;not null" json:"city"`
	PostalCode string `gorm:"size:100;not null" json:"postal_code"`
	Country    string `gorm:"size:100;not null" json:"country"`
}
