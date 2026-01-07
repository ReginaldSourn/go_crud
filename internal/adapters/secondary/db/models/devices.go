package models

import (
	"time"
)

type Device struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Status    string    `gorm:"type:varchar(50);not null" json:"status"`
	Telemetry string    `gorm:"type:text" json:"telemetry"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Device) TableName() string {
	return "devices"
}
