package model

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey" redis:"id"`
	Name      string    `gorm:"not null" redis:"name"`
	Email     string    `gorm:"not null;unique" redis:"email"`
	CreatedAt time.Time `gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"not null;default:current_timestamp"`
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) GetPrimaryKey() string {
	return "id"
}

func (u *User) TableName() string {
	return "users"
}
