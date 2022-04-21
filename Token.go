package appstoreconnect

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/dgrijalva/jwt-go"
	errortools "github.com/leapforce-libraries/go_errortools"
)

type JwtToken struct {
	Token  string
	Expiry time.Time
}

func (service *Service) getToken() (string, *errortools.Error) {
	now := time.Now()

	if service.jwtToken != nil {
		if service.jwtToken.Token != "" {
			if service.jwtToken.Expiry.After(now.Add(time.Minute)) {
				return service.jwtToken.Token, nil
			}
		}
	}

	expirationTime := now.Add(time.Duration(jwtTokenExpiryMinutes) * time.Minute)

	claims := &jwt.StandardClaims{
		Audience:  service.audience,
		IssuedAt:  now.Unix(),
		Issuer:    service.issuerId,
		ExpiresAt: expirationTime.Unix(),
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = service.keyId

	// Create the JWT string
	block, _ := pem.Decode([]byte(service.privateKey))
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errortools.ErrorMessage(err)
	}

	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return "", errortools.ErrorMessage(err)
	}

	service.jwtToken = &JwtToken{
		Token:  tokenString,
		Expiry: expirationTime,
	}

	return service.jwtToken.Token, nil
}
