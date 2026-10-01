package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`

	User User `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	TokenHash string `gorm:"uniqueIndex;not null" json:"-"`

	ExpiresAt time.Time `gorm:"not null;index" json:"-"`

	RevokedAt *time.Time `json:"-"`

	UserAgent string `gorm:"size:500" json:"-"`

	IPAddress string `gorm:"size:100" json:"-"`

	CreatedAt time.Time `json:"created_at"`
}

func (token *RefreshToken) BeforeCreate() error {
	token.ID = uuid.New()
	return nil
}