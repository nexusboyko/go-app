package handlers

import (
	"app/internal/db"
	"app/internal/middleware"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var envLoaded = godotenv.Load()

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/callback",
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"openid https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, Path: "/"}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, []byte, error) {
	// fetch token and user profile information

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()

	// get auth profile and auth token JSON string

	authProfile, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	authToken, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("failed marshaling token: %s", err.Error())
	}

	return authProfile, authToken, nil
}

func OAuthGoogleLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c.Writer)

	authURL := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)

	c.Writer.Header().Set("HX-Redirect", authURL)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func OAuthGoogleCallback(c *gin.Context) {
	oauthState, err := c.Request.Cookie("oauthstate")

	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if c.Query("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	profile, token, err := getUserDataFromGoogle(c.Query("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	var newUserProfile db.AuthProfile
	json.Unmarshal(profile, &newUserProfile)
	var newUserToken oauth2.Token
	json.Unmarshal(token, &newUserToken)

	db.AddUser(db.User{
		Profile:   newUserProfile,
		AuthToken: newUserToken,
	})

	middleware.SetSession(c)

	c.Redirect(http.StatusPermanentRedirect, "/home")
}
