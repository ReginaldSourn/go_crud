package domain

import "time"

type Device struct {
	ID        int64      `gorm:"primaryKey;type:bigserial" json:"id"`
	Name      string     `gorm:"size:128;not null" json:"name"`
	TypeID    int64      `gorm:"not null" json:"type_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}
