package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/team968/Parts/db"
	"github.com/team968/Parts/middlewares"
	"github.com/team968/Parts/routes/auth"
	"github.com/team968/Parts/routes/orders"
	"github.com/team968/Parts/routes/parts"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "assets/")

	router.Use(middlewares.AuthBefore)

	router.POST("/login", auth.POSTLogin)
	router.GET("/login", auth.GETLogin)

	router.Use(middlewares.AuthAfter)

	router.POST("/register", auth.POSTRegister)
	router.GET("/register", auth.GETRegister)

	router.GET("/logout", auth.GETLogout)

	router.Static("/pdf", "GrabCAD/2018 FRC Robot/PDFS/")
	router.Static("/part", "GrabCAD/2018 FRC Robot/")

	router.GET("/overview", parts.GETOverview)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/overview")
	})

	projectGroup := router.Group("/project")
	{
		specifiedProjectGroup := projectGroup.Group("/:prefix/:year")
		{
			partGroup := specifiedProjectGroup.Group("/part")
			{
				partGroup.GET("/:partNumber/edit", parts.GETEditPart)
				partGroup.POST("/:partNumber/edit", parts.POSTEditPart)

				partGroup.GET("/:partNumber/move-status-up", parts.GETMovePartUp)
				partGroup.GET("/:partNumber/move-status-down", parts.GETMovePartDown)

				partGroup.GET("/:partNumber/delete", parts.GETDeletePart)

				partGroup.GET("/:partNumber", parts.GETPart)
			}

			assemblyGroup := specifiedProjectGroup.Group("/assembly")
			{
				assemblyGroup.GET("/:partNumber/edit", parts.GETEditAssembly)
				assemblyGroup.POST("/:partNumber/edit", parts.POSTEditAssembly)

				assemblyGroup.GET("/:partNumber/delete", parts.GETDeleteAssembly)

				assemblyGroup.GET("/:partNumber", parts.GETAssembly)
			}

			specifiedProjectGroup.GET("/new-part", parts.GETNewPart)
			specifiedProjectGroup.POST("/new-part", parts.POSTNewPart)

			specifiedProjectGroup.GET("/new-assembly", parts.GETNewAssembly)
			specifiedProjectGroup.POST("/new-assembly", parts.POSTNewAssembly)

			specifiedProjectGroup.GET("/", parts.GETProject)
		}
	}

	ordersGroup := router.Group("/orders")
	{
		ordersGroup.GET("/", orders.GETOverview)
		ordersGroup.POST("/edit/:id", orders.POSTEditOrder)
		ordersGroup.GET("/edit/:id", orders.GETEditOrder)
	}

	tvGroup := router.Group("/tv")
	{
		tvGroup.GET("/left", parts.GETTVLeft)
		tvGroup.GET("/right", parts.GETTVRight)
	}

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"title": "Page not found",
		})
	})

	db.Connect()

	router.Run(":8888")
}
