package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
	"golang.org/x/crypto/bcrypt"
)

func GETRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", gin.H{
		"title": "Register",
	})
}

func POSTRegister(c *gin.Context) {
	id := c.PostForm("id")
	passwordStr := c.PostForm("password")

	password, err := bcrypt.GenerateFromPassword([]byte(passwordStr), 10)

	if err != nil {
		return
	}

	db.Database.Create(&models.User{
		ID:       id,
		Password: password,
	})
}
