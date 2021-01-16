package main

import (
	"github.com/francisco-serrano/github-oauth-poc/github"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	app := gin.Default()

	app.GET("/", github.Index)
	app.GET("/oauth/github/receive", github.ReceiveAuthCode)
	app.POST("/oauth/github/login", github.Login)
	app.GET("/partial-register", github.PartialRegister)
	app.GET("/register", github.Register)

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
