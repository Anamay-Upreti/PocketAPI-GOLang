package handlers

import (
	"github.com/gofiber/fiber/v2"

	"pocketapi/internal/services"
	"pocketapi/internal/utils"
	"pocketapi/internal/validators"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: services.NewAuthService(),
	}
}

// =========================
// REGISTER
// =========================

type RegisterBody struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (handler *AuthHandler) Register(c *fiber.Ctx) error {

	var body RegisterBody

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(
			c,
			400,
			"Invalid JSON body",
		)
	}

	if err := validators.Validate.Struct(body); err != nil {
		return utils.Error(
			c,
			400,
			validators.FormatValidationError(err),
		)
	}

	user, err := handler.service.Register(
		services.RegisterRequest{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		},
	)

	if err != nil {
		return utils.Error(
			c,
			400,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		201,
		"User registered successfully",
		user,
	)
}

// =========================
// LOGIN
// =========================

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (handler *AuthHandler) Login(c *fiber.Ctx) error {

	var body LoginBody

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(
			c,
			400,
			"Invalid JSON body",
		)
	}

	// Validate request body
	if err := validators.Validate.Struct(body); err != nil {
		return utils.Error(
			c,
			400,
			validators.FormatValidationError(err),
		)
	}

	// Login and create access + refresh tokens
	accessToken, refreshToken, user, err := handler.service.Login(
		services.LoginRequest{
			Email:    body.Email,
			Password: body.Password,
		},
		c.Get("User-Agent"),
		c.IP(),
	)

	if err != nil {
		return utils.Error(
			c,
			401,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		200,
		"Login successful",
		fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

// =========================
// REFRESH TOKEN
// =========================

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (handler *AuthHandler) Refresh(c *fiber.Ctx) error {

	var body RefreshRequest

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(
			c,
			400,
			"Invalid JSON body",
		)
	}

	// Validate refresh token
	if err := validators.Validate.Struct(body); err != nil {
		return utils.Error(
			c,
			400,
			validators.FormatValidationError(err),
		)
	}

	accessToken, refreshToken, err := handler.service.Refresh(
		body.RefreshToken,
		c.Get("User-Agent"),
		c.IP(),
	)

	if err != nil {
		return utils.Error(
			c,
			401,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		200,
		"Token refreshed successfully",
		fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

// =========================
// LOGOUT
// =========================

func (handler *AuthHandler) Logout(c *fiber.Ctx) error {

	var body RefreshRequest

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(
			c,
			400,
			"Invalid JSON body",
		)
	}

	if err := validators.Validate.Struct(body); err != nil {
		return utils.Error(
			c,
			400,
			validators.FormatValidationError(err),
		)
	}

	err := handler.service.Logout(body.RefreshToken)

	if err != nil {
		return utils.Error(
			c,
			401,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		200,
		"Logged out successfully",
		nil,
	)
}

// =========================
// LOGOUT ALL DEVICES
// =========================

func (handler *AuthHandler) LogoutAll(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	err := handler.service.LogoutAll(userID)

	if err != nil {
		return utils.Error(
			c,
			500,
			"Failed to logout from all devices",
		)
	}

	return utils.Success(
		c,
		200,
		"Logged out from all devices",
		nil,
	)
}
