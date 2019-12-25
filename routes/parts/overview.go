package parts

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

func GETOverview(c *gin.Context) {
	var projects []models.Project

	dbc := db.Database.Find(&projects)

	if dbc.Error != nil {
		return
	}

	stats := make([][]int, len(projects))

	design := 0
	manufacturing := 0
	router := 0

	for i, project := range projects {
		parts := getAllParts(project.Prefix, project.Year)

		stats[i] = make([]int, len(models.StatusMap))

		for _, part := range parts {
			stats[i][part.Status]++

			if part.Status >= 0 && part.Status <= 2 {
				design++
			}

			if part.Status >= 5 && part.Status <= 8 {
				manufacturing++
			}

			if part.NeedsRouter && part.Status >= 5 && part.Status <= 8 {
				router++
			}
		}
	}

	c.HTML(http.StatusOK, "parts_overview.tmpl", gin.H{
		"title":         "Projects",
		"projects":      projects,
		"stats":         stats,
		"design":        design,
		"manufacturing": manufacturing,
		"router":        router,
	})
}
