package middleware

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	// MemberCode string `json:"member_code"`
	Channel   string `json:"channel"`
	Role      string `json:"role"`
	Workspace string `json:"workspace"`
	// Settings   DataSetting `json:"settings"`
	jwt.StandardClaims
}

func VerifyJwtToken(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := &CustomClaims{}
		tokenAuth := ctx.Request.Header.Get("Authorization")

		// Parse token
		_, err := jwt.ParseWithClaims(tokenAuth, claims, func(t *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != t.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			// Handling error dengan detail
			msg := "token is invalid"
			if mErr, ok := err.(*jwt.ValidationError); ok {
				if mErr.Errors == jwt.ValidationErrorExpired {
					msg = "token is expired"
				}
			}

			// Response error
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": msg,
			})
			ctx.Abort() // Hentikan eksekusi berikutnya
			return
		}

		// Set data yang bisa diakses di handler berikutnya
		ctx.Set("role", claims.Role)
		ctx.Set("workspace", claims.Workspace)

		ctx.Next()
	}
}
