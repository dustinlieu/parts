package models

import "github.com/jinzhu/gorm"

var StatusMap = map[int]string{
	0: "In concepts",
	1: "Design in progress",
	2: "Needs drawings",
	3: "Material needs to be ordered",
	4: "Waiting for materials",
	5: "Ready to cut",
	6: "Ready for mill",
	7: "Ready for lathe",
	8: "Ready for post processes",
	9: "Done",
}

var SecondaryStatusMap = [...]string{
	"Needs facing",
	"Faced",
	"Milling in progress",
}

var PriorityMap = [...]string{
	"Low",
	"Normal",
	"High",
}

var MaterialMap = [...]string{
	"Aluminium",
	"Steel",
	"Polycarbonate",
	"Delrin",
	"Plywood",
	"Modified COTS",
	"3DP - PLA",
	"3DP - ABS",
	"3DP - TPS",
}

type Part struct {
	GenericPart
	Material          int
	HaveMaterial      bool
	MaterialCutLength float64
	NeedsRouter       bool
	Priority          int
}

func (p Part) GetURL(db *gorm.DB) string {
	project, err := p.GetParentProject(db)

	if err != nil {
		return ""
	}

	return project.GetURL(db) + "/part/" + p.PartNumber
}
