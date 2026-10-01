package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"

	"pocketapi/internal/services"
	"pocketapi/internal/utils"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler() *UserHandler {

	return &UserHandler{
		service: services.NewUserService(),
	}
}

func (handler *UserHandler) UploadAvatar(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	file, err := c.FormFile("avatar")

	if err != nil {
		return utils.Error(c, 400, "Avatar file is required")
	}

	if file.Size > utils.MaxAvatarSize {
		return utils.Error(c, 400, "Maximum avatar size is 2MB")
	}

	if err := utils.ValidateImage(file.Filename); err != nil {
		return utils.Error(c, 400, err.Error())
	}

	err = os.MkdirAll("uploads", os.ModePerm)

	if err != nil {
		return utils.Error(c, 500, "Unable to create upload directory")
	}

	filename := utils.GenerateAvatarFilename(file.Filename)

	path := filepath.Join("uploads", filename)

	err = c.SaveFile(file, path)

	if err != nil {
		return utils.Error(c, 500, "Failed to save avatar")
	}

	avatarURL := fmt.Sprintf("/uploads/%s", filename)

	err = handler.service.UpdateAvatar(userID, avatarURL)

	if err != nil {
		return utils.Error(c, 500, err.Error())
	}

	return utils.Success(
		c,
		200,
		"Avatar uploaded successfully",
		fiber.Map{
			"avatar_url": avatarURL,
		},
	)
}

func (handler *UserHandler) GetProfile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	user, err := handler.service.GetProfile(userID)
	if err != nil {
		return utils.Error(c, 404, "User not found")
	}

	return utils.Success(
		c,
		200,
		"Profile fetched successfully",
		user,
	)
}