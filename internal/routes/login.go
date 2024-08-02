package routes

import (
	oauth "app/internal/auth"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	oauth.OAuthGoogleLogin(c)
}
