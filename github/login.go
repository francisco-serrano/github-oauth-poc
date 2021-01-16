package github

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func Login(ctx *gin.Context) {
	randomState := uuid.New().String()

	loginAttempts[randomState] = time.Now().Add(time.Hour)

	externalLoginURL := oauthConfig.AuthCodeURL(randomState, oauth2.AccessTypeOffline)

	ctx.Redirect(http.StatusSeeOther, externalLoginURL)
}
