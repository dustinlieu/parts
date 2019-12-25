package parts

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/bradfitz/slice"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

func resolvePart(c *gin.Context) (models.Part, models.Project, error) {
	prefix := c.Param("prefix")
	year := c.Param("year")
	partNumber := c.Param("partNumber")

	var project models.Project
	var part models.Part

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return part, project, dbc.Error
	}

	dbc = db.Database.First(&part, "project = ? AND part_number = ?", project.ID, partNumber)

	if dbc.Error != nil {
		return part, project, dbc.Error
	}

	return part, project, nil
}

func resolveProject(c *gin.Context) (models.Project, error) {
	prefix := c.Param("prefix")
	year := c.Param("year")

	var project models.Project

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return project, dbc.Error
	}

	return project, nil
}

func getParents(parent int) []models.Assembly {
	if parent == -1 {
		return []models.Assembly{}
	}

	sorted := make([]models.Assembly, 0, 0)

	for {
		var assembly models.Assembly

		var dbc *gorm.DB

		if len(sorted) == 0 {
			dbc = db.Database.First(&assembly, "id = ?", parent)
		} else {
			dbc = db.Database.First(&assembly, "id = ?", sorted[len(sorted)-1].Parent)
		}

		if dbc.Error != nil {
			break
		}

		sorted = append(sorted, assembly)

		if sorted[len(sorted)-1].Parent == -1 {
			break
		}

		for _, a := range sorted {
			if assembly.Parent == int(a.ID) {
				break
			}
		}
	}

	flipped := make([]models.Assembly, len(sorted), len(sorted))

	for i := len(sorted) - 1; i >= 0; i-- {
		flipped[len(sorted)-i-1] = sorted[i]
	}

	return flipped
}

func getAllParts(prefix string, year string) []models.Part {
	var project models.Project
	var parts []models.Part

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return parts
	}

	db.Database.Find(&parts, "project = ?", project.ID)

	slice.Sort(parts[:], func(i int, j int) bool {
		return sort.StringsAreSorted([]string{parts[i].PartNumber, parts[j].PartNumber})
	})

	return parts
}

func getPartsUnder(prefix string, year string, parent uint) []models.Part {
	var project models.Project
	var parts []models.Part

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return parts
	}

	db.Database.Find(&parts, "project = ? AND parent = ?", project.ID, parent)

	slice.Sort(parts[:], func(i int, j int) bool {
		return sort.StringsAreSorted([]string{parts[i].PartNumber, parts[j].PartNumber})
	})

	return parts
}

func GETPart(c *gin.Context) {
	part, project, err := resolvePart(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})

		return
	}

	parents := getParents(part.Parent)

	c.HTML(http.StatusOK, "parts_part.tmpl", gin.H{
		"title":           part.GetFullPartNumber(db.Database),
		"subtitle":        part.Name,
		"project":         project,
		"parents":         parents,
		"lastParentIndex": len(parents) - 1,
		"part":            part,
		"statusMap":       models.StatusMap,
		"priorityMap":     models.PriorityMap,
		"materialMap":     models.MaterialMap,
	})
}

func GETEditPart(c *gin.Context) {
	part, project, err := resolvePart(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		assemblies := getAllAssemblies(c.Param("prefix"), c.Param("year"))

		c.HTML(http.StatusOK, "parts_part_edit.tmpl", gin.H{
			"title":       "Editing Part",
			"project":     project,
			"assemblies":  assemblies,
			"part":        part,
			"statusMap":   models.StatusMap,
			"materialMap": models.MaterialMap,
		})
	}
}

func POSTEditPart(c *gin.Context) {
	part, _, err := resolvePart(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		part.PartNumber = c.DefaultPostForm("part_number", part.PartNumber)
		part.Name = c.DefaultPostForm("part_name", part.Name)

		parent, err := strconv.Atoi(c.DefaultPostForm("parent", strconv.Itoa(part.Parent)))
		if err == nil {
			part.Parent = parent
		}

		status, err := strconv.Atoi(c.DefaultPostForm("status", strconv.Itoa(part.Status)))
		if err == nil {
			part.Status = status
		}

		material, err := strconv.Atoi(c.DefaultPostForm("material", strconv.Itoa(part.Material)))
		if err == nil {
			part.Material = material
		}

		part.HaveMaterial = c.DefaultPostForm("have_material", "0") == "on"
		part.NeedsRouter = c.DefaultPostForm("needs_router", "0") == "on"

		materialCutLength, err := strconv.ParseFloat(c.DefaultPostForm("material_cut_length", strconv.FormatFloat(part.MaterialCutLength, 'f', -1, 32)), 10)
		if err == nil {
			part.MaterialCutLength = materialCutLength
		}

		quantity, err := strconv.Atoi(c.DefaultPostForm("quantity", strconv.Itoa(part.Quantity)))
		if err == nil {
			part.Quantity = quantity
		}

		priority, err := strconv.Atoi(c.DefaultPostForm("priority", strconv.Itoa(part.Priority)))
		if err == nil {
			part.Priority = priority
		}

		db.Database.Save(&part)

		c.Redirect(http.StatusFound, part.GetURL(db.Database))
	}
}

func GETNewPart(c *gin.Context) {
	project, err := resolveProject(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		var assemblies []models.Assembly

		dbc := db.Database.Find(&assemblies)

		if dbc.Error != nil {
			c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
				"title": "Page not found",
			})

			return
		}

		c.HTML(http.StatusOK, "parts_part_new.tmpl", gin.H{
			"title":       "New Part",
			"subtitle":    project.Name,
			"project":     project,
			"assemblies":  assemblies,
			"statusMap":   models.StatusMap,
			"materialMap": models.MaterialMap,
		})
	}
}

func POSTNewPart(c *gin.Context) {
	project, err := resolveProject(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		var part models.Part

		part.Project = project.ID

		part.PartNumber = c.PostForm("part_number")
		part.Name = c.PostForm("part_name")

		parent, err := strconv.Atoi(c.PostForm("parent"))
		if err == nil {
			part.Parent = parent
		}

		status, err := strconv.Atoi(c.PostForm("status"))
		if err == nil {
			part.Status = status
		}

		material, err := strconv.Atoi(c.PostForm("material"))
		if err == nil {
			part.Material = material
		}

		part.HaveMaterial = c.DefaultPostForm("have_material", "0") == "on"
		part.NeedsRouter = c.DefaultPostForm("needs_router", "0") == "on"

		materialCutLength, err := strconv.ParseFloat(c.PostForm("material_cut_length"), 10)
		if err == nil {
			part.MaterialCutLength = materialCutLength
		}

		quantity, err := strconv.Atoi(c.PostForm("quantity"))
		if err == nil {
			part.Quantity = quantity
		}

		priority, err := strconv.Atoi(c.PostForm("priority"))
		if err == nil {
			part.Priority = priority
		}

		db.Database.Create(&part)

		c.Redirect(http.StatusFound, part.GetURL(db.Database))
	}
}

func GETDeletePart(c *gin.Context) {
	part, project, err := resolvePart(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		db.Database.Delete(&part)

		c.Redirect(http.StatusFound, project.GetURL(db.Database))
	}
}

func GETMovePartUp(c *gin.Context) {
	part, project, err := resolvePart(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	}

	if part.SecondaryStatus != 0 {
		part.SecondaryStatus--
		db.Database.Save(&part)
	}

	c.Redirect(http.StatusFound, project.GetURL(db.Database)+"?filter-by=mill")
}

func GETMovePartDown(c *gin.Context) {
	part, project, err := resolvePart(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	}

	if part.SecondaryStatus < len(models.SecondaryStatusMap)-1 {
		part.SecondaryStatus++
		db.Database.Save(&part)
	}

	c.Redirect(http.StatusFound, project.GetURL(db.Database)+"?filter-by=mill")
}
