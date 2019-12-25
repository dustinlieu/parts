package parts

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

type item struct {
	PartNumber      string
	Type            string
	Name            string
	Parent          string
	Status          int
	SecondaryStatus int
	Material        string
	CutLength       float64
	NeedsRouter     bool
	Quantity        int
}

var FilterMap = map[string]string{
	"design":          "Design",
	"manufacturing":   "Manufacturing",
	"in_concepts":     "In concepts",
	"in_design":       "Design in progress",
	"drawings_needed": "Needs drawings",
	"cut":             "Ready to cut",
	"mill":            "Ready for mill",
	"lathe":           "Ready for lathe",
	"post":            "Ready for post processes",
	"router_needed":   "Needs router",
}

func GETProject(c *gin.Context) {
	prefix := c.Param("prefix")
	year := c.Param("year")
	filterOption := strings.ToLower(c.DefaultQuery("filter-by", "none"))

	var project models.Project

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})

		return
	}

	var wantedStatusLower int = 0
	var wantedStatusUpper int = len(models.StatusMap) - 1
	if filterOption == "design" {
		wantedStatusLower = 0
		wantedStatusUpper = 2
	} else if filterOption == "manufacturing" {
		wantedStatusLower = 5
		wantedStatusUpper = 8
	} else if filterOption == "in_concepts" {
		wantedStatusLower = 0
		wantedStatusUpper = 0
	} else if filterOption == "in_design" {
		wantedStatusLower = 1
		wantedStatusUpper = 1
	} else if filterOption == "drawings_needed" {
		wantedStatusLower = 2
		wantedStatusUpper = 2
	} else if filterOption == "cut" {
		wantedStatusLower = 5
		wantedStatusUpper = 5
	} else if filterOption == "mill" {
		wantedStatusLower = 6
		wantedStatusUpper = 6
	} else if filterOption == "lathe" {
		wantedStatusLower = 7
		wantedStatusUpper = 7
	} else if filterOption == "post" {
		wantedStatusLower = 8
		wantedStatusUpper = 8
	} else if filterOption == "router_needed" {
		wantedStatusLower = 5
		wantedStatusUpper = 8
	}

	parts := getAllParts(project.Prefix, project.Year)
	assemblies := getAllAssemblies(project.Prefix, project.Year)
	items := make([]item, 0)

	for _, assembly := range assemblies {
		if assembly.Status >= wantedStatusLower && assembly.Status <= wantedStatusUpper {
			items = append(items, item{
				PartNumber:      assembly.PartNumber,
				Type:            "Assembly",
				Name:            assembly.Name,
				Status:          assembly.Status,
				SecondaryStatus: assembly.SecondaryStatus,
			})

			for _, a := range assemblies {
				if int(a.ID) == assembly.Parent {
					items[len(items)-1].Parent = a.PartNumber
					break
				}
			}
		}
	}

	for _, part := range parts {
		if part.Status >= wantedStatusLower && part.Status <= wantedStatusUpper {
			items = append(items, item{
				PartNumber:      part.PartNumber,
				Type:            "Part",
				Name:            part.Name,
				Status:          part.Status,
				SecondaryStatus: part.SecondaryStatus,
				Material:        models.MaterialMap[part.Material],
				CutLength:       part.MaterialCutLength,
				NeedsRouter:     part.NeedsRouter,
				Quantity:        part.Quantity,
			})

			for _, a := range assemblies {
				if int(a.ID) == part.Parent {
					items[len(items)-1].Parent = a.PartNumber
					break
				}
			}
		}
	}

	c.HTML(http.StatusOK, "parts_project.tmpl", gin.H{
		"title":                     project.Name,
		"subtitle":                  "Parts and Assemblies",
		"project":                   project,
		"statusMap":                 models.StatusMap,
		"secondaryStatusMap":        models.SecondaryStatusMap,
		"finalSecondaryStatusIndex": len(models.SecondaryStatusMap) - 1,
		"items":                     items,
		"filterOption":              filterOption,
		"filterMap":                 FilterMap,
		"wantedStatusLower":         wantedStatusLower,
		"wantedStatusUpper":         wantedStatusUpper,
	})
}

func GETNewProject(c *gin.Context) {
	c.HTML(http.StatusOK, "parts_project_new.tmpl", gin.H{
		"title": "New Project",
	})
}
