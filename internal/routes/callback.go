package routes

import (
	oauth "app/internal/auth"

	"github.com/gin-gonic/gin"
)

func Callback(c *gin.Context) {
	oauth.OAuthGoogleCallback(c)
}
