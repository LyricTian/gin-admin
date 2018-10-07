package router

import (
	"gin-admin/src/api"
	"gin-admin/src/context"

	"github.com/gin-gonic/gin"
)

// APIDemoRouter 注册/demos路由
func APIDemoRouter(g *gin.RouterGroup, demo *api.Demo) {
	g.GET("/demos", context.WrapContext(demo.QueryList, "查询示例列表"))
	g.GET("/demos/:id", context.WrapContext(demo.Get, "查询单条示例数据"))
	g.POST("/demos", context.WrapContext(demo.Create, "创建示例数据"))
	g.PUT("/demos/:id", context.WrapContext(demo.Update, "更新示例数据"))
	g.DELETE("/demos/:id", context.WrapContext(demo.Delete, "删除示例数据"))
}
