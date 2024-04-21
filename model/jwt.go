package model

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	Phone string `json:"phone"` // Add custom claims as needed
}
