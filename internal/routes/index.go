package routes

import (
	m "app/internal/middleware"
	"app/internal/templates"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	if m.IsLoggedIn(c) {
		c.Redirect(http.StatusPermanentRedirect, "/home")
	}

	err := templates.Layout(templates.Index()).Render(c, c.Writer)

	if err != nil {
		log.Println("Error rendering template (index):", err.Error())
	}
}
