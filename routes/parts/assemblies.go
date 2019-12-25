package parts

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/bradfitz/slice"
	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/models"
)

func resolveAssembly(c *gin.Context) (models.Assembly, models.Project, error) {
	prefix := c.Param("prefix")
	year := c.Param("year")
	partNumber := c.Param("partNumber")

	var project models.Project
	var assembly models.Assembly

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return assembly, project, dbc.Error
	}

	dbc = db.Database.First(&assembly, "project = ? AND part_number = ?", project.ID, partNumber)

	if dbc.Error != nil {
		return assembly, project, dbc.Error
	}

	return assembly, project, nil
}

func getAllAssemblies(prefix string, year string) []models.Assembly {
	var project models.Project
	var assemblies []models.Assembly

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return assemblies
	}

	db.Database.Find(&assemblies, "project = ?", project.ID)

	slice.Sort(assemblies[:], func(i int, j int) bool {
		return sort.StringsAreSorted([]string{assemblies[i].PartNumber, assemblies[j].PartNumber})
	})

	return assemblies
}

func getAssembliesUnder(prefix string, year string, parent uint) []models.Assembly {
	var project models.Project
	var assemblies []models.Assembly

	dbc := db.Database.First(&project, "prefix = ? AND year = ?", prefix, year)

	if dbc.Error != nil {
		return assemblies
	}

	db.Database.Find(&assemblies, "project = ? AND parent = ?", project.ID, parent)

	slice.Sort(assemblies[:], func(i int, j int) bool {
		return sort.StringsAreSorted([]string{assemblies[i].PartNumber, assemblies[j].PartNumber})
	})

	return assemblies
}

func filterList(vs []item, f func(item) bool) []item {
	vsf := make([]item, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func GETAssembly(c *gin.Context) {
	assembly, project, err := resolveAssembly(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		parts := getPartsUnder(project.Prefix, project.Year, assembly.ID)
		assemblies := getAssembliesUnder(project.Prefix, project.Year, assembly.ID)

		items := make([]item, len(parts)+len(assemblies))

		for i, a := range assemblies {
			items[i] = item{
				PartNumber: a.PartNumber,
				Type:       "Assembly",
				Name:       a.Name,
				Status:     a.Status,
			}
		}

		for i, p := range parts {
			items[len(assemblies)+i] = item{
				PartNumber: p.PartNumber,
				Type:       "Part",
				Name:       p.Name,
				Status:     p.Status,
			}
		}

		c.HTML(http.StatusOK, "parts_assembly.tmpl", gin.H{
			"title":     assembly.GetFullPartNumber(db.Database),
			"subtitle":  assembly.Name,
			"project":   project,
			"assembly":  assembly,
			"parents":   getParents(assembly.Parent),
			"statusMap": models.StatusMap,
			"items":     items,
		})
	}
}

func GETEditAssembly(c *gin.Context) {
	assembly, project, err := resolveAssembly(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		assemblies := getAllAssemblies(project.Prefix, project.Year)

		c.HTML(http.StatusOK, "parts_assembly_edit.tmpl", gin.H{
			"title":      "Editing Assembly",
			"subtitle":   assembly.GetFullPartNumber(db.Database),
			"project":    project,
			"assemblies": assemblies,
			"assembly":   assembly,
			"statusMap":  models.StatusMap,
		})
	}
}

func POSTEditAssembly(c *gin.Context) {
	assembly, _, err := resolveAssembly(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		assembly.PartNumber = c.DefaultPostForm("assembly_number", assembly.PartNumber)
		assembly.Name = c.DefaultPostForm("assembly_name", assembly.Name)

		parent, err := strconv.Atoi(c.DefaultPostForm("parent", strconv.Itoa(assembly.Parent)))
		if err == nil {
			assembly.Parent = parent
		}

		status, err := strconv.Atoi(c.DefaultPostForm("status", strconv.Itoa(assembly.Status)))
		if err == nil {
			assembly.Status = status
		}

		quantity, err := strconv.Atoi(c.DefaultPostForm("quantity", strconv.Itoa(assembly.Quantity)))
		if err == nil {
			assembly.Quantity = quantity
		}

		db.Database.Save(&assembly)

		c.Redirect(http.StatusFound, assembly.GetURL(db.Database))
	}
}

func GETNewAssembly(c *gin.Context) {
	project, err := resolveProject(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		assemblies := getAllAssemblies(project.Prefix, project.Year)

		c.HTML(http.StatusOK, "parts_assembly_new.tmpl", gin.H{
			"title":       "New Assembly",
			"subtitle":    project.Name,
			"project":     project,
			"assemblies":  assemblies,
			"statusMap":   models.StatusMap,
			"materialMap": models.MaterialMap,
		})
	}
}

func POSTNewAssembly(c *gin.Context) {
	project, err := resolveProject(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		var assembly models.Assembly

		assembly.Project = project.ID

		assembly.PartNumber = c.PostForm("part_number")
		assembly.Name = c.PostForm("assembly_name")

		parent, err := strconv.Atoi(c.PostForm("parent"))
		if err == nil {
			assembly.Parent = parent
		}

		status, err := strconv.Atoi(c.PostForm("status"))
		if err == nil {
			assembly.Status = status
		}

		quantity, err := strconv.Atoi(c.PostForm("quantity"))
		if err == nil {
			assembly.Quantity = quantity
		}

		db.Database.Create(&assembly)

		c.Redirect(http.StatusFound, assembly.GetURL(db.Database))
	}
}

func GETDeleteAssembly(c *gin.Context) {
	assembly, project, err := resolveAssembly(c)

	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	} else {
		var assemblyChanges models.Assembly

		assemblyChanges.Parent = -1

		db.Database.Model(models.Part{}).Where("parent = ?", assembly.ID).Updates(assemblyChanges)
		db.Database.Model(models.Assembly{}).Where("parent = ?", assembly.ID).Updates(assemblyChanges)

		db.Database.Delete(&assembly)

		c.Redirect(http.StatusFound, project.GetURL(db.Database))
	}
}
