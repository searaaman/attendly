package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/searaaman/attendly/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func RegisterUser(db *gorm.DB, name string, email string, password string) (*domain.User, error) {
	var existingUser domain.User
	result := db.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AuthenticateUser returns the same error for an unknown email and a wrong
// password, so the login endpoint doesn't reveal which emails are registered.
func AuthenticateUser(db *gorm.DB, email string, password string) (*domain.User, error) {
	var user domain.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return &user, nil
}

func GenerateToken(user *domain.User, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"iat":     now.Unix(),
		"exp":     now.Add(ttl).Unix(),
	})
	return token.SignedString([]byte(secret))
}

// ParseToken validates the signature, the signing method and the expiry.
func ParseToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
