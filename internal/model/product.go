package model

import "time"

type Product struct {
	ID          uint64    `gorm:"primaryKey" redis:"id"`
	Name        string    `gorm:"not null" redis:"name"`
	Description string    `gorm:"type:text" redis:"description"`
	Price       float64   `gorm:"not null" redis:"price"`
	Quantity    uint64    `gorm:"not null" redis:"quantity"`
	CreatedAt   time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt   time.Time `gorm:"not null;default:current_timestamp"`
}

func (p *Product) GetID() uint64 {
	return p.ID
}

func (p *Product) GetPrimaryKey() string {
	return "id"
}

func (p *Product) TableName() string {
	return "products"
}
