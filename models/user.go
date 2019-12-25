package models

const (
	Viewer uint = iota
	Editor
	Administrator
)

type User struct {
	ID         string `gorm:"primary_key"`
	Password   []byte
	Name       string
	Email      string
	Permission uint
}
