package models

type Product struct {
    ID          uint    `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string  `gorm:"required" json:"name,omitempty"`
}
