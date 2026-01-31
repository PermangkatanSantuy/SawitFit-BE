package auth

import (
	"log"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	jwks *keyfunc.JWKS
}

func NewVerifier(jwksURL string) (*Verifier, error) {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Printf("JWKS refresh error: %v", err)
		},
	})

	if err != nil {
		return nil, err
	}

	return &Verifier{jwks: jwks}, nil
}

func (v *Verifier) Verify(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, v.jwks.Keyfunc)
}