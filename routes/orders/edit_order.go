package orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GETEditOrder(c *gin.Context) {
	c.HTML(http.StatusOK, "orders_edit.tmpl", gin.H{
		"title": "Edit Order",
	})
}

func POSTEditOrder(c *gin.Context) {

}
