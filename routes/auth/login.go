package auth

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

func GETLogin(c *gin.Context) {
	_, exists := c.Get("user")

	if exists {
		c.Redirect(http.StatusFound, "/")
	}

	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":   "Login",
		"hideNav": true,
	})
}

func POSTLogin(c *gin.Context) {
	id := c.PostForm("id")
	password := c.PostForm("password")

	var user models.User
	dbc := db.Database.First(&user, "id = ?", id)

	if dbc.Error != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	token := generateToken(user.ID)
	db.Database.Create(token)

	cookie := http.Cookie{
		Name:    "_SID",
		Value:   token.SessionID,
		Expires: time.Now().Add(365 * 24 * time.Hour),
	}
	http.SetCookie(c.Writer, &cookie)

	c.Redirect(http.StatusFound, "/")
}
