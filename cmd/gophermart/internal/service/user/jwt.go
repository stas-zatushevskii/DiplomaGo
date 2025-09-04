package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint
}

func (u *ServiceUser) BuildJwt(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(u.config.JWT.Expire) * time.Hour)),
			Issuer:    "gophermart",
		},

		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(u.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("%w: %v", customErrors.ErrJWTBuild, err)
	}

	return tokenString, nil
}

func (u *ServiceUser) GetUserIDFromJwt(tokenString string) (uint, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(u.config.JWT.Secret), nil
		})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	return claims.UserID, nil
}
