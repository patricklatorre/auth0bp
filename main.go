package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/patricklatorre/auth0bp/api"
	"github.com/patricklatorre/auth0bp/middleware"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	api.InitCorsOrigin("http://" + addr)
	api.InitAuthConfig(api.AuthConfig{
		Domain:   os.Getenv("AUTH0_DOMAIN"),
		ClientId: os.Getenv("AUTH0_CLIENT_ID"),
		Audience: os.Getenv("AUTH0_AUDIENCE"),
	})

	e := echo.New()

	publicApi := e.Group("/api/public/v0")
	publicApi.GET("/auth_config", api.AuthConfigHandler)

	privateApi := e.Group("/api/private/v0")
	privateApi.Use(echo.WrapMiddleware(middleware.EnsureValidToken()))

	err = e.Start(addr)
	if err != nil {
		e.Logger.Fatal(err)
	}
}
