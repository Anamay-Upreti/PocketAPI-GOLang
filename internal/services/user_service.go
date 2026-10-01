package services

import (
	"pocketapi/internal/models"
	"pocketapi/internal/repositories"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService() *UserService {

	return &UserService{
		repository: repositories.NewUserRepository(),
	}
}

func (service *UserService) GetProfile(userID string) (*models.User, error) {

	user, err := service.repository.FindByID(userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (service *UserService) UpdateAvatar(
	userID string,
	avatarPath string,
) error {

	return service.repository.UpdateAvatar(userID, avatarPath)
}