package models

import (
	"time"
)

type Order struct {
	ID           uint   `gorm:"primary key"`
	CustomerName string `gorm:"not null; type: varchar(191)"`
	Items        []Item
	OrderedAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
