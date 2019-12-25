package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/team968/Parts/models"
)

var Database *gorm.DB

func Connect() {
	var err error
	Database, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/parts")

	if err != nil {
		defer Database.Close()
		panic("Unable to connect to databse")
	}

	Database.AutoMigrate(&models.Order{})
	Database.AutoMigrate(&models.Assembly{})
	Database.AutoMigrate(&models.Part{})
	Database.AutoMigrate(&models.Project{})
	Database.AutoMigrate(&models.Token{})
	Database.AutoMigrate(&models.User{})
}
