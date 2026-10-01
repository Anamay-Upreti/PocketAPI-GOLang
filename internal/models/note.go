package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	Title string `gorm:"size:200;not null" json:"title"`

	Content string `gorm:"type:text" json:"content"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`

	User User `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}