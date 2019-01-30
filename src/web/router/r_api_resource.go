package router

import (
	"github.com/LyricTian/gin-admin/src/web/ctl"
	"github.com/gin-gonic/gin"
)

// APIResourceRouter 注册/resources路由
func APIResourceRouter(g *gin.RouterGroup, c *ctl.Resource) {
	GET(g, "/v1/resources", c.Query, SetName("查询资源数据"))
	GET(g, "/v1/resources/:id", c.Get, SetName("查询指定资源数据"))
	POST(g, "/v1/resources", c.Create, SetName("创建资源数据"))
	PUT(g, "/v1/resources/:id", c.Update, SetName("更新资源数据"))
	DELETE(g, "/v1/resources/:id", c.Delete, SetName("删除资源数据"))
	DELETE(g, "/v1/resources", c.DeleteMany, SetName("删除多条资源数据"))
}
