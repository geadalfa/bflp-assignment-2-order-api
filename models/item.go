package models

import (
	"time"
)

type Item struct {
	ID          uint   `gorm:"primaryKey"`
	ItemCode    string `gorm:"not null;type:varchar(255)"`
	Description string
	Quantity    int `gorm:"not null; type: int"`
	OrderID     uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
