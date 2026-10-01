package services

import (
	"time"

	"github.com/google/uuid"

	"pocketapi/internal/models"
	"pocketapi/internal/repositories"
	"pocketapi/internal/utils"
)

type AuthService struct {
	repository        *repositories.UserRepository
	refreshRepository *repositories.RefreshTokenRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		repository:        repositories.NewUserRepository(),
		refreshRepository: repositories.NewRefreshTokenRepository(),
	}
}

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

// =========================
// REGISTER
// =========================

func (service *AuthService) Register(
	req RegisterRequest,
) (*models.User, error) {

	// Check whether email already exists
	_, err := service.repository.FindByEmail(req.Email)

	if err == nil {
		return nil, utils.NewBadRequest("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, utils.NewInternal(
			"failed to hash password",
			err,
		)
	}

	// Create user
	user := models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	err = service.repository.Create(&user)
	if err != nil {
		return nil, utils.NewInternal(
			"failed to create user",
			err,
		)
	}

	// Never return password
	user.Password = ""

	return &user, nil
}

// =========================
// LOGIN
// =========================

func (service *AuthService) Login(
	req LoginRequest,
	userAgent string,
	ipAddress string,
) (string, string, *models.User, error) {

	// Find user
	user, err := service.repository.FindByEmail(req.Email)

	if err != nil {
		return "", "", nil,
			utils.NewUnauthorized(
				"invalid email or password",
			)
	}

	// Check password
	err = utils.CheckPassword(
		user.Password,
		req.Password,
	)

	if err != nil {
		return "", "", nil,
			utils.NewUnauthorized(
				"invalid email or password",
			)
	}

	// Generate access token
	accessToken, err := utils.GenerateToken(
		user.ID.String(),
		user.Email,
		user.Role,
	)

	if err != nil {
		return "", "", nil,
			utils.NewInternal(
				"failed to generate access token",
				err,
			)
	}

	// Create refresh-token session
	refreshToken, err := service.CreateRefreshSession(
		user.ID,
		userAgent,
		ipAddress,
	)

	if err != nil {
		return "", "", nil,
			utils.NewInternal(
				"failed to create refresh session",
				err,
			)
	}

	// Never return password
	user.Password = ""

	return accessToken, refreshToken, user, nil
}

// =========================
// CREATE REFRESH SESSION
// =========================

func (service *AuthService) CreateRefreshSession(
	userID uuid.UUID,
	userAgent string,
	ipAddress string,
) (string, error) {

	// Generate random refresh token
	rawToken, err := utils.GenerateRefreshToken()

	if err != nil {
		return "", err
	}

	// Store only hash in database
	tokenHash := utils.HashRefreshToken(rawToken)

	refreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	// Save session
	err = service.refreshRepository.Create(&refreshToken)

	if err != nil {
		return "", err
	}

	// Return raw token to client
	return rawToken, nil
}

// =========================
// REFRESH ACCESS TOKEN
// =========================

func (service *AuthService) Refresh(
	rawRefreshToken string,
	userAgent string,
	ipAddress string,
) (string, string, error) {

	// Hash incoming token
	tokenHash := utils.HashRefreshToken(
		rawRefreshToken,
	)

	// Find active refresh token
	oldToken, err := service.refreshRepository.FindActiveByHash(
		tokenHash,
	)

	if err != nil {
		return "", "",
			utils.NewUnauthorized(
				"invalid or expired refresh token",
			)
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(
		oldToken.User.ID.String(),
		oldToken.User.Email,
		oldToken.User.Role,
	)

	if err != nil {
		return "", "",
			utils.NewInternal(
				"failed to generate access token",
				err,
			)
	}

	// Generate new refresh token
	newRawToken, err := utils.GenerateRefreshToken()

	if err != nil {
		return "", "",
			utils.NewInternal(
				"failed to generate refresh token",
				err,
			)
	}

	// Create new refresh-token record
	newRefreshToken := models.RefreshToken{
		ID:        uuid.New(),
		UserID:    oldToken.UserID,
		TokenHash: utils.HashRefreshToken(newRawToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}

	// Rotate old token -> revoke old + create new
	err = service.refreshRepository.Rotate(
		oldToken.ID,
		&newRefreshToken,
	)

	if err != nil {
		return "", "",
			utils.NewInternal(
				"refresh token rotation failed",
				err,
			)
	}

	return accessToken, newRawToken, nil
}

// =========================
// LOGOUT
// =========================

func (service *AuthService) Logout(
	rawRefreshToken string,
) error {

	// Hash incoming token
	tokenHash := utils.HashRefreshToken(
		rawRefreshToken,
	)

	// Find active token
	token, err := service.refreshRepository.FindActiveByHash(
		tokenHash,
	)

	if err != nil {
		return utils.NewUnauthorized(
			"invalid or expired refresh token",
		)
	}

	// Revoke token
	err = service.refreshRepository.Revoke(
		token.ID,
	)

	if err != nil {
		return utils.NewInternal(
			"failed to logout",
			err,
		)
	}

	return nil
}

// =========================
// LOGOUT ALL DEVICES
// =========================

func (service *AuthService) LogoutAll(
	userID string,
) error {

	err := service.refreshRepository.RevokeAll(
		userID,
	)

	if err != nil {
		return utils.NewInternal(
			"failed to logout from all devices",
			err,
		)
	}

	return nil
}
