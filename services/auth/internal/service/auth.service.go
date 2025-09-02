package service

import (
	"context"
	"errors"
	"time"

	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/internal/events"
	"github.com/imdinnesh/openfinstack/services/auth/internal/repository"
	"github.com/imdinnesh/openfinstack/services/auth/models"
	"github.com/imdinnesh/openfinstack/services/auth/redis"
	"github.com/imdinnesh/openfinstack/services/auth/utils"
)

type AuthService interface {
	RegisterUser(email, password string) (*models.User, error)
	LoginUser(email, password string) (string, string, error)
	RefreshToken(oldRefreshToken string) (string, string, error)
	RevokeToken(token string) error
	Profile(userID uint) (*models.User, error)
}

type authService struct {
	userRepo      repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	cfg           *config.Config
	redis         *redis.Client
	publisher     *events.UserEventPublisher
}

func NewAuthService(repo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository, cfg *config.Config, rds *redis.Client, publisher *events.UserEventPublisher) AuthService {
	return &authService{
		userRepo:      repo,
		refreshTokenRepo: refreshTokenRepo,
		cfg:           cfg,
		redis:        rds,
		publisher:    publisher,
	}
}

// RegisterUser hashes password and creates user
func (s *authService) RegisterUser(email, password string) (*models.User, error) {

	userExists, _ := s.userRepo.FindByEmail(email)
	if userExists != nil {
		return nil, errors.New("user already exists")
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashed,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	if err := s.publisher.PublishUserCreated(context.Background(), user.ID, user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser checks credentials and generates tokens
func (s *authService) LoginUser(email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", err
	}

	if err := utils.CheckPassword(password, user.PasswordHash); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := s.generateJWT(user.ID, 15*time.Minute,user.Role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateJWT(user.ID, 7*24*time.Hour, user.Role)
	if err != nil {
		return "", "", err
	}

	// Save the refresh token in the database
	if err := s.refreshTokenRepo.Create(&models.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken issues new tokens if valid
func (s *authService) RefreshToken(oldRefreshToken string) (string, string, error) {
	claims, err := s.parseJWT(oldRefreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}


	userID := uint(claims["user_id"].(float64))

	isExpired,err:=s.refreshTokenRepo.IsExpired(userID);

	if err!=nil{
		return "","",errors.New("expired refresh token");
	}

	if isExpired{
		return "","",errors.New("expired refresh token");
	}

	accessToken, err := s.generateJWT(userID, 15*time.Minute,"user")
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateJWT(userID, 7*24*time.Hour,"user")
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RevokeToken adds token to Redis blacklist
func (s *authService) RevokeToken(token string) error {
	return s.redis.BlacklistToken(token, 15*time.Minute)
}

func (s *authService) Profile(userID uint) (*models.User, error) {
	user, err := s.userRepo.FindById(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
