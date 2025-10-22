package model

type Product struct {
	ProductID    uint    `gorm:"primaryKey" json:"product_id"`
	ProductName  string  `gorm:"size:250; not null" json:"product_name"`
	ProductPrice float64 `gorm:"not null" json:"price"`
	ProductStock int     `gorm:"not null" json:"stock_qty"`

	CategoryID   uint   `gorm:"not null" json:"category_id"`
	CategoryName string `gorm:"size:100; not null" json:"category_name"`

	ProductImages []ProductImage `gorm:"foreignKey:ProductID" json:"product_images"`
}
