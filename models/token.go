package models

type Token struct {
	SessionID   string `gorm:"primary_key"`
	UserID      string
	LastUpdated int64
}
