package middleware

import (
	"errors"
	"net/http"
	"time"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type jwtMiddleware struct {
	accesstokenKey string
}

func NewJWTMiddleware(accessTokenKey string) *jwtMiddleware {
	return &jwtMiddleware{
		accesstokenKey: accessTokenKey,
	}
}

func (m *jwtMiddleware) validateToken(accessToken string) (*entity.AccessToken, error) {
	var payload entity.AccessToken
	token, err := jwt.ParseWithClaims(accessToken, &payload, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.accesstokenKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &payload, nil
}

func (m *jwtMiddleware) clearCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

func (m *jwtMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/static*" {
			return next(c)
		}
		cookie, err := c.Cookie("accessToken")
		if err == nil {
			accessToken := cookie.Value
			payload, err := m.validateToken(accessToken)
			if err == nil {
				if payload.UserId == "" {
					m.clearCookie(c)
				} else {
					c.Set("userId", payload.UserId)
					c.Set("isLogin", true)
				}
			} else {
				m.clearCookie(c)
			}
		}
		return next(c)
	}
}
