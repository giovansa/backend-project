package internal

import (
	"crypto/rsa"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func GenerateJWTToken(user model.User, env string) (string, error) {
	privateKey, err := loadPrivateKey(env)
	if err != nil {
		return "", err
	}

	claims := &model.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
		Phone: user.Phone,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func loadPrivateKey(env string) (*rsa.PrivateKey, error) {
	if env == "unit_test" {
		keyDer, err := os.ReadFile("../private.pem")
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPrivateKeyFromPEM(keyDer)
	}
	keyDer, err := os.ReadFile("./private.pem")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(keyDer)
}

func LoadPublicKey() (*rsa.PublicKey, error) {
	keyDer, err := os.ReadFile("./public.pem")
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyDer)
}
