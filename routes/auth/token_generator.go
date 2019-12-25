package auth

import (
	"math/rand"
	"time"

	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func generateToken(id string) models.Token {
	var sessionID string

	for {
		var token models.Token

		sessionID = generateRandomString(32)
		dbc := db.Database.First(&token, "session_id = ?", sessionID)

		if dbc.Error != nil {
			break
		}
	}

	return models.Token{
		SessionID:   sessionID,
		UserID:      id,
		LastUpdated: time.Now().Unix(),
	}
}
