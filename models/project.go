package models

import (
	"github.com/jinzhu/gorm"
)

type Project struct {
	ID     uint `gorm:"primary_key"`
	Prefix string
	Year   string
	Name   string
}

func (p Project) GetURL(db *gorm.DB) string {
	return "/project/" + p.Prefix + "/" + p.Year
}
