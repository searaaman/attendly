package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/searaaman/attendly/internal/domain"
)

func TestTokenRoundTrip(t *testing.T) {
	user := &domain.User{Email: "host@example.com"}
	user.ID = 7

	token, err := GenerateToken(user, "secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseToken(token, "secret")
	if err != nil {
		t.Fatal(err)
	}
	if claims["email"] != "host@example.com" || claims["user_id"] != float64(7) {
		t.Errorf("unexpected claims %v", claims)
	}
}

func TestParseTokenRejectsWrongSecret(t *testing.T) {
	token, _ := GenerateToken(&domain.User{}, "secret", time.Hour)

	if _, err := ParseToken(token, "other-secret"); err == nil {
		t.Error("expected error for wrong secret")
	}
}

func TestParseTokenRejectsExpired(t *testing.T) {
	token, _ := GenerateToken(&domain.User{}, "secret", -time.Minute)

	if _, err := ParseToken(token, "secret"); err == nil {
		t.Error("expected error for expired token")
	}
}

func TestParseTokenRejectsMissingExpiry(t *testing.T) {
	// Tokens issued before expiry was added have no exp claim and must be rejected.
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1}).SignedString([]byte("secret"))

	if _, err := ParseToken(token, "secret"); err == nil {
		t.Error("expected error for token without exp")
	}
}
