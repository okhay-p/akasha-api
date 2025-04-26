package handlers

import (
	"akasha-api/pkg/config"
	"akasha-api/pkg/jwt"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"akasha-api/internal/gothic"
	"akasha-api/internal/model"
	"akasha-api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"gorm.io/gorm"
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

	store.Options.Domain = "akashalearn.org"
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = os.Getenv("DEV") == "false"

	gothic.Store = store

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "https://api.akashalearn.org/auth/google/callback"),
	)
}

func HandleLogin(c *gin.Context) {
	if !config.Dev {
		c.Status(http.StatusUnauthorized)
		c.Abort()
		log.Printf("TRYING TO ACCESS DEV LOG IN ROUTE")
		return
	}

	// TEMP LOG IN HANDLER MUST CHANGE AFTER GOOGLE OAUTH SETUP
	token, err := jwt.CreateToken(config.AkashaLearnUUID)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"auth": token})
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
	log.Println(user.Email)
	log.Println(user.UserID)

	userUUID, err := services.GetUserUUIDByGoogleSub(user.UserID)
	if err != nil && err == gorm.ErrRecordNotFound {
		var userB model.AlUser

		userB.Email = user.Email
		userB.GoogleSubID = user.UserID
		userB.Username = strings.Split(user.Email, "@")[0]
		userB.Status = 1
		userB.PictureURL = user.AvatarURL

		userUUID, err = services.InsertNewUser(&userB)

		if err != nil {
			log.Println("Error:", err)
		}
	}
	log.Println(userUUID)

	token, err := jwt.CreateToken(userUUID.String())
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, 86400, "/", "akashalearn.org", true, true)
	c.Redirect(http.StatusFound, "https://akashalearn.org/generate-lessons")
}

func HandleOAuth(c *gin.Context) {

	if config.Dev {
		log.Println("DEV MODE LOGIN")
		token, err := jwt.CreateToken(config.AkashaLearnUUID)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			c.Abort()
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"auth": token})
		return
	}

	provider := c.Param("provider")

	c.Request = c.Request.WithContext(context.WithValue(context.Background(), "provider", provider))

	log.Println(provider)

	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		log.Println("Oauth handler:")
		log.Println(gothUser.Email)
		log.Println(gothUser.UserID)

	} else {
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func HandleOAuthLogout(c *gin.Context) {

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "/", "akashalearn.org", true, true)
	c.Status(http.StatusOK)
}
