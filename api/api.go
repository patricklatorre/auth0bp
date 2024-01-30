package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AuthConfig struct {
	Domain   string `json:"domain"`
	ClientId string `json:"clientId"`
	Audience string `json:"audience"`
}

var corsOrigin string
var authConfig *AuthConfig

func InitAuthConfig(config AuthConfig) {
	authConfig = &AuthConfig{
		Domain:   config.Domain,
		ClientId: config.ClientId,
		Audience: config.Audience,
	}
}

func InitCorsOrigin(origin string) {
	corsOrigin = origin
}

func AuthConfigHandler(c echo.Context) error {
	return c.JSON(200, *authConfig)
}

func HelloHandler(c echo.Context) error {
	accessToken, err := getAccessTokenFromContext(c)
	if err != nil {
		log.Println("Can't extract access token from context")
		return err
	}

	sub, err := unverifiedParseSubClaim(accessToken)
	if err != nil {
		log.Println("Can't decode access token:", err)
		return err
	}

	// Debug
	log.Println("sub (uid) from access token:", sub)

	c.Request().Header.Set("Access-Control-Allow-Origin", corsOrigin)
	c.Request().Header.Set("Access-Control-Allow-Credentials", "true")
	c.Request().Header.Set("Access-Control-Allow-Headers", "Authorization")

	return c.String(200, `{"msg": "Hello back!"}`)
}

////////////////////////////////////////////////////////////////////////////////

func getAccessTokenFromContext(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return authHeaderParts[1], nil
}

// Parses the sub claim from the jwt without validating the signature. This
// assumes that the token was already validated in a previous middleware.
func unverifiedParseSubClaim(tokenString string) (string, error) {
	var sub string

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub = fmt.Sprint(claims["sub"])
	}

	if sub == "" {
		return "", fmt.Errorf("invalid token payload. Token %v", token)
	}

	return sub, nil
}
