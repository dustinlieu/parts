package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type GenericPart struct {
	ID              uint `gorm:"primary_key"`
	Project         uint
	Parent          int
	PartNumber      string
	Name            string
	Status          int
	SecondaryStatus int
	Quantity        int
}

func (g GenericPart) GetParentProject(db *gorm.DB) (Project, error) {
	var project Project

	dbc := db.First(&project, "id = ?", g.Project)

	if dbc.Error != nil {
		return project, dbc.Error
	}

	return project, nil
}

func (g GenericPart) GetParentAssembly(db *gorm.DB) (Assembly, error) {
	var assembly Assembly

	if g.Parent == -1 {
		return assembly, nil
	}

	dbc := db.First(&assembly, "id = ?", g.Parent)

	if dbc.Error != nil {
		return assembly, dbc.Error
	}

	return assembly, nil
}

func (g GenericPart) GetFullPartNumber(db *gorm.DB) string {
	project, err := g.GetParentProject(db)

	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s-%s-%s", project.Prefix, project.Year, g.PartNumber)
}
