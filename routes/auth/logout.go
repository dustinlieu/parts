package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
	"net/http"
)

func GETLogout(c *gin.Context) {
	user, _ := c.Get("user")

	var token models.Token
	dbc := db.Database.First(&token, "user_id = ?", user.(models.User).ID)

	if dbc.Error != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db.Database.Delete(&token)

	c.Redirect(http.StatusFound, "/login")
}
