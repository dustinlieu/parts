package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

func AuthBefore(c *gin.Context) {
	defer c.Next()

	cookie, err := c.Cookie("_SID")

	if err != nil {
		return
	}

	var token models.Token
	dbc := db.Database.First(&token, "session_id = ?", cookie)

	if dbc.Error != nil {
		return
	}

	if token.LastUpdated+1800 < time.Now().Unix() {
		db.Database.Delete(&token)

		return
	}

	token.LastUpdated = time.Now().Unix()
	db.Database.Save(&token)

	var user models.User
	dbc = db.Database.First(&user, token.UserID)

	if dbc.Error != nil {
		return
	}

	c.Set("user", user)
}

func AuthAfter(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var token models.Token
	dbc := db.Database.First(&token, "user_id = ?", user.(models.User).ID)

	if dbc.Error != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Next()
}
