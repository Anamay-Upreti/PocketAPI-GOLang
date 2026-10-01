package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"pocketapi/internal/database"
	"pocketapi/internal/models"
)

type RefreshTokenRepository struct{}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{}
}

// Create a refresh token session.
func (repo *RefreshTokenRepository) Create(
	token *models.RefreshToken,
) error {

	return database.DB.Create(token).Error
}

// Find active refresh token using its hash.
func (repo *RefreshTokenRepository) FindActiveByHash(
	tokenHash string,
) (*models.RefreshToken, error) {

	var token models.RefreshToken

	err := database.DB.
		Preload("User").
		Where(
			"token_hash = ? AND revoked_at IS NULL AND expires_at > ?",
			tokenHash,
			time.Now(),
		).
		First(&token).Error

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Revoke one refresh token.
func (repo *RefreshTokenRepository) Revoke(
	tokenID uuid.UUID,
) error {

	now := time.Now()

	return database.DB.
		Model(&models.RefreshToken{}).
		Where("id = ?", tokenID).
		Update("revoked_at", now).Error
}

// Revoke all sessions for a user.
func (repo *RefreshTokenRepository) RevokeAll(
	userID string,
) error {

	now := time.Now()

	return database.DB.
		Model(&models.RefreshToken{}).
		Where(
			"user_id = ? AND revoked_at IS NULL",
			userID,
		).
		Update("revoked_at", now).Error
}

// Rotate refresh token atomically.
func (repo *RefreshTokenRepository) Rotate(
	oldTokenID uuid.UUID,
	newToken *models.RefreshToken,
) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		now := time.Now()

		result := tx.
			Model(&models.RefreshToken{}).
			Where(
				"id = ? AND revoked_at IS NULL",
				oldTokenID,
			).
			Update("revoked_at", now)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}

		if err := tx.Create(newToken).Error; err != nil {
			return err
		}

		return nil
	})
}