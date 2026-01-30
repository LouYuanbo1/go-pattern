package model

import "time"

type Order struct {
	ID        uint64    `gorm:"primaryKey" redis:"id"`
	UserID    uint64    `gorm:"not null" redis:"user_id"`
	ProductID uint64    `gorm:"not null" redis:"product_id"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp"`
}

func (o *Order) GetID() uint64 {
	return o.ID
}

func (o *Order) GetPrimaryKey() string {
	return "id"
}

func (o *Order) TableName() string {
	return "orders"
}
