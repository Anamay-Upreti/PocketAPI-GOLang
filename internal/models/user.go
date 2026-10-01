package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {

	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name string `gorm:"size:100;not null"`

	Email string `gorm:"size:150;uniqueIndex;not null"`

	Password string `gorm:"not null" json:"-"`

	Role string `gorm:"default:user"`

	Avatar string `gorm:"default:''" json:"avatar"`

	CreatedAt time.Time

	UpdatedAt time.Time
}


func (user *User) BeforeCreate(tx *gorm.DB) error {

	user.ID = uuid.New()

	return nil
}