package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Search 以图搜图
func Search(c *gin.Context) {
	c.String(http.StatusOK, "%s", "search")
}
