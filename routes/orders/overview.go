package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GETOverview(c *gin.Context) {
	c.HTML(http.StatusOK, "orders_overview.tmpl", gin.H{
		"title": "Orders Overview",
	})
}
