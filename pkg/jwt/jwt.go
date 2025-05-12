package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	// MemberCode string `json:"member_code"`
	Channel   string `json:"channel"`
	Role      string `json:"role"`
	Workspace string `json:"workspace"`
	// Settings   DataSetting `json:"settings"`
	jwt.StandardClaims
}

type RegisterClaims struct {
	Token string `json:"token_register"`
	jwt.StandardClaims
}

func GenerateToken(secret string) (signedToken string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := &CustomClaims{
		Channel:   "asu",
		Role:      "asu",
		Workspace: "asu",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    "1",
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := rawToken.SignedString([]byte(secret))

	return token, err
}
