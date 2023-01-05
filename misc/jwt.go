package misc

import (
	"github.com/golang-jwt/jwt"
)

type JWT struct {
	SigningKey []byte
}

type EmailJwtClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewJWT(signingKey []byte) *JWT {
	return &JWT{
		SigningKey: signingKey,
	}
}

func (j *JWT) GenToken(method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	return token.SignedString(j.SigningKey)
}

func (j *JWT) Parse(tokenStr string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, err error) {
		return j.SigningKey, nil
	})

	return token, err
}
