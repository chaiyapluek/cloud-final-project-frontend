package entity

import "github.com/golang-jwt/jwt/v5"

type AccessToken struct {
	*jwt.RegisteredClaims
	UserId string `json:"userId"`
}
