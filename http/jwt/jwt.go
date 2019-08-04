package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chuabingquan/snippets"
	"github.com/dgrijalva/jwt-go"
)

// Authenticator implements the http.Authenticator interface
type Authenticator struct {
	SigningKey []byte
	ExpiryTime time.Duration
}

// GetAuthorizationInfo extracts and returns the authorization information of the owner
// of the given authentication token
func (a Authenticator) GetAuthorizationInfo(r *http.Request) (snippets.AuthorizationInfo, error) {
	authorizationInfo := snippets.AuthorizationInfo{}

	tokenString, err := getTokenFromHeader(r)
	if err != nil {
		return authorizationInfo, err
	}

	token, err := jwt.Parse(tokenString, a.keyGetter)
	if err != nil {
		return authorizationInfo, errors.New("Failed to extract token from string: " + err.Error())
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	authorizationInfo.UserID = claims["userId"].(string)

	return authorizationInfo, nil
}

// GenerateToken creates a new authentication token for the purpose of authorization
func (a Authenticator) GenerateToken(info snippets.AuthorizationInfo) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["userId"] = info.UserID
	claims["exp"] = time.Now().Add(a.ExpiryTime).Unix()

	tokenString, err := token.SignedString(a.SigningKey)
	if err != nil {
		return "", errors.New("Signing failed during token generation: " + err.Error())
	}

	return tokenString, nil
}

// Authenticate verifies the legitimacy and validity of an authentication token
func (a Authenticator) Authenticate(r *http.Request) (bool, error) {
	tokenString, err := getTokenFromHeader(r)
	if err != nil {
		return false, err
	}

	_, err = jwt.Parse(tokenString, a.keyGetter)
	if err != nil { // error occurs when token is invalid
		return false, nil
	}
	return true, nil
}

// getTokenFromHeader extracts and returns an authentication token from the request header,
// else, it returns an error should an invalid token format be supplied
func getTokenFromHeader(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	tokenParts := strings.Split(tokenString, " ")

	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("Invalid token format supplied")
	}

	return tokenParts[1], nil
}

// keyGetter is a callback function for jwt.Parse that allows custom logic/validation to
// be performed on a token before returning a signing key to verify the token's signature
func (a Authenticator) keyGetter(token *jwt.Token) (interface{}, error) {
	if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method != jwt.SigningMethodHS256 {
		return nil, fmt.Errorf("Unexpected signing method used: %v", token.Header["alg"])
	}
	return a.SigningKey, nil
}
