package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/deckforge/backend/internal/config"
	"github.com/deckforge/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTClaims holds data we embed inside JWT tokens.
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// AuthService handles registration, login, and token validation.
type AuthService struct {
	db     *gorm.DB
	secret []byte
	expiry time.Duration
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		secret: []byte(cfg.JWTSecret),
		expiry: time.Duration(cfg.JWTExpiryHrs) * time.Hour,
	}
}

// Register creates a new user with a hashed password.
func (s *AuthService) Register(email, password, fullName string) (*models.User, string, error) {
	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, "", errors.New("email already registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := models.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}
	return &user, token, nil
}

// Login verifies credentials and returns a JWT.
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid email or password")
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := s.generateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}
	return &user, token, nil
}

func (s *AuthService) generateToken(userID uuid.UUID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken parses and verifies a JWT, returning claims if valid.
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
