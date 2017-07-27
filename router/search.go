package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/feature"
	"net/http"
)

// Search 以图搜图
func Search(c *gin.Context) {
	fmt.Println(GetImgFeature([]byte("lxm")))
	c.String(http.StatusOK, "%s", "search")
}
