package parts

import (
	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
	"net/http"
	"time"
	"math"
	"strconv"
)

type assemblyStats struct {
	PartNumber         string
	Name               string
	DesignCount        int
	ManufacturingCount int
}

func floatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(math.Floor(input_num), 'f', 0, 64)
}

func GETTVLeft(c *gin.Context) {
	var project models.Project

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", "968", "18")

	if dbc.Error != nil {
		return
	}

	assemblies := getAllAssemblies("968", "18")

	var items []assemblyStats

	for _, assembly := range assemblies {
		parts := getPartsUnder("968", "18", assembly.ID)

		stats := make([]int, len(models.StatusMap))

		for _, part := range parts {
			stats[part.Status]++
		}

		items = append(items, assemblyStats{
			PartNumber:         assembly.GetFullPartNumber(db.Database),
			Name:               assembly.Name,
			DesignCount:        stats[0]+stats[1],
			ManufacturingCount: stats[4]+stats[5]+stats[6]+stats[7],
		})
	}

	t, err := time.Parse(time.UnixDate, "Tue Feb 20 23:59:00 EST 2018")

	if err != nil {
		t = time.Now()
	}

	duration := t.Sub(time.Now())

	c.HTML(http.StatusOK, "tv_left.tmpl", gin.H{
		"title":      "Editing Part",
		"hideHeader": true,
		"prefix":     project.Prefix + "-" + project.Year,
		"assemblies": assemblies,
		"items":      items,
		"daysLeft": floatToString(duration.Hours()/24) + ":" + floatToString(math.Mod(duration.Hours(), 24)) + ":" + floatToString(math.Mod(duration.Minutes(), 60)) + ":" + floatToString(math.Mod(duration.Seconds(), 60)),
	})
}

func GETTVRight(c *gin.Context) {
	c.HTML(http.StatusOK, "tv_right.tmpl", gin.H{})
}
