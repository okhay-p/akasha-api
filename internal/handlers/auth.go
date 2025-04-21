package handlers

import (
	"akasha-api/pkg/config"
	"akasha-api/pkg/jwt"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"akasha-api/internal/gothic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "KFeiuh174yafbo33kfabab34knfaueh9r3ku8ef48dkbGWI"
	MaxAge = 86400 * 30
)

func NewAuth() {
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("DEV") == "false"

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "https://api.akashalearn.org/auth/google/callback"),
	)
}

func HandleLogin(c *gin.Context) {
	// TEMP LOG IN HANDLER MUST CHANGE AFTER GOOGLE OAUTH SETUP
	token, err := jwt.CreateToken(config.AkashaLearnUUID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	if config.Dev {
		c.IndentedJSON(http.StatusOK, gin.H{"auth": token})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("test", "test", 300, "/", "localhost", false, false)
	c.SetCookie("token", token, 86400, "/", "localhost", true, true)
	c.IndentedJSON(http.StatusOK, gin.H{"auth": "success"})
}

func HandleOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")

	c.Request = c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Println("ERROR:", err)
		fmt.Fprintln(c.Writer, c.Request)
		return
	}
	log.Println("Oauth callback:")
	log.Println(user.Name)
	fmt.Println(user)

	c.Redirect(http.StatusFound, "https://akashalearn.org/")
}

func HandleOAuth(c *gin.Context) {

	provider := c.Param("provider")

	c.Request = c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	log.Println(provider)

	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		log.Println("Oauth handler:")
		log.Println(gothUser.Name)
	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func HandleOAuthLogout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)
	c.Writer.Header().Set("Location", "/")
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)
}
