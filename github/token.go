package github

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const key = "abc123"

type customClaims struct {
	jwt.StandardClaims
	SessionID      string
	ExternalUserID string
	InternalUserID string
}

type githubData struct {
	ID      string `json:"id"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

func CreateToken(sessionID, internalUserID, externalUserID string) (string, error) {
	claims := &customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		SessionID:      sessionID,
		InternalUserID: internalUserID,
		ExternalUserID: externalUserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error while signing JWT: %w", err)
	}

	return signedToken, nil
}

func ParseTokenSessionID(token string) (string, error) {
	claims, err := parseToken(token)
	if err != nil {
		return "", err
	}

	return claims.SessionID, nil
}

func ParseTokenInternalUserID(token string) (string, error) {
	claims, err := parseToken(token)
	if err != nil {
		return "", err
	}

	return claims.InternalUserID, nil
}

func parseToken(token string) (*customClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while parsing JWT: %s", err)
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid auth")
	}

	claims, ok := parsedToken.Claims.(*customClaims)
	if !ok {
		return nil, fmt.Errorf("could not obtain auth claims: %T", parsedToken.Claims)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("auth expired")
	}

	return claims, nil
}
