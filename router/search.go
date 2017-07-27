package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/feature"
	. "github.com/orange-jacky/albums/util"
	"net/http"
)

// Search 以图搜图
func Search(c *gin.Context) {
	//host
	conf := Configure("")
	hostport := fmt.Sprintf("%s:%s", conf.Feature.Host, conf.Feature.Port)
	fmt.Println(GetImgFeature([]byte("lxm"), hostport))
	c.String(http.StatusOK, "%s", "search")
}
