package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"user_id"`
	Email     string    `gorm:"unique;not null" json:"user_email"`
	Name      string    `gorm:"not null" json:"user_name"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`

	Products []Product `gorm:"many2many:user_products;" json:"products"`
}

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"product_id"`
	URL       string    `gorm:"not null" json:"product_url"`
	LastPrice float64   `gorm:"not null" json:"product_last_price"`
	ImageURL  string    `json:"image_url"`
	Name      string    `gorm:"not null" json:"product_name"`
	UpdatedAt time.Time `json:"updated_at"`

	Users []User `gorm:"many2many:user_products;" json:"users"`
}
