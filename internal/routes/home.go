package routes

import (
	"app/internal/db"
	"app/internal/templates"
	"log"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	user := db.GetUsers()[0]
	err := templates.Layout(templates.Home(user)).Render(c, c.Writer)

	if err != nil {
		log.Println("Error rendering template (index):", err.Error())
	}
}
