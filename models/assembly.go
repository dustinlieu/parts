package models

import (
	"github.com/jinzhu/gorm"
)

type Assembly struct {
	GenericPart
}

func (a Assembly) GetURL(db *gorm.DB) string {
	project, err := a.GetParentProject(db)

	if err != nil {
		return ""
	}

	return project.GetURL(db) + "/assembly/" + a.PartNumber
}
