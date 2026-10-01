package repositories

import (
	"pocketapi/internal/database"
	"pocketapi/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Find user using email.
func (repo *UserRepository) FindByEmail(email string) (*models.User, error) {

	var user models.User

	err := database.DB.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create user.
func (repo *UserRepository) Create(user *models.User) error {

	return database.DB.Create(user).Error
}

func (repo *UserRepository) FindByID(id string) (*models.User, error) {

	var user models.User

	err := database.DB.
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) UpdateAvatar(userID string, avatar string) error {

	return database.DB.
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar", avatar).Error
}