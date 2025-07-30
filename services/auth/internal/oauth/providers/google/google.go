package google

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/models"
	"github.com/imdinnesh/openfinstack/services/auth/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var googleOAuthConfig = &oauth2.Config{
	RedirectURL:  config.Load().GoogleRedirectURL,
	ClientID:     config.Load().GoogleClientID,
	ClientSecret: config.Load().GoogleClientSecret,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

type GoogleProvider struct {
	db *gorm.DB
}

func NewGoogleProvider(db *gorm.DB) *GoogleProvider {
	return &GoogleProvider{db: db}
}

type GoogleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	ID    string `json:"id"`
}

type OAuthResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (g *GoogleProvider) GetLoginURL(deviceID string) string {
	return googleOAuthConfig.AuthCodeURL(deviceID)
}

func (g *GoogleProvider) HandleCallback(code, deviceID string) (*OAuthResponse, error) {
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %w", err)
	}
	defer resp.Body.Close()

	var gUser GoogleUser
	err = json.NewDecoder(resp.Body).Decode(&gUser)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	var user models.User
	dbErr := g.db.Where("email = ?", gUser.Email).First(&user).Error
	if dbErr == gorm.ErrRecordNotFound {
		user = models.User{
			Email: gUser.Email,
			IsVerified: true,
			Role: "user",
			CreatedAt: time.Now(),
		}
		g.db.Create(&user)
	}

	accessToken, err := utils.GenerateJWT(user.ID, 15*time.Minute, user.Role)
	refreshToken, err := utils.GenerateJWT(user.ID, 7*24*time.Hour, user.Role)
	if err != nil {
		return nil, fmt.Errorf("create refresh token failed: %w", err)
	}

	return &OAuthResponse{
		Status:       "success",
		Message:      "Logged in with Google",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
