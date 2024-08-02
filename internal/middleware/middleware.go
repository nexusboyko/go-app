package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", `Content-Type, Content-Length,
			Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With`)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

/* AUTH MIDDLEWARE */

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsLoggedIn(c) {
			c.Redirect(http.StatusMovedPermanently, "/")
			c.Abort()
		}
	}
}

func IsLoggedIn(c *gin.Context) bool {
	cookie, err := c.Cookie("session")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	_, ok := store[cookie]

	return ok
}

/* SESSION MIDDLEWARE */

var store = make(map[string]bool)

func SetSession(c *gin.Context) {
	id := generateSessionId()

	store[id] = true
	c.SetCookie("session", id, int((time.Minute * 3).Seconds()), "/", "", true, true)

	c.Next()
}

func ClearSession(c *gin.Context) {
	id, err := c.Cookie("session")
	if err != nil {
		log.Println(err.Error())
		return
	}

	delete(store, id)
	c.SetCookie("session", "", -1, "/", "", true, true)

	c.Next()
}

func generateSessionId() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}
