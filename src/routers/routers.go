package routers

import (
	"com.phpstu/webapp/src/routers/approuter"
	"com.phpstu/webapp/src/routers/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Load(g *gin.Engine) {
	//g.Use(gin.Recovery())
	g.Use(middlewares.GinLogger(), middlewares.GinRecovery(false))

	g.Use(middlewares.SetRequestId)
	//不缓存
	g.Use(middlewares.NoCache)
	//options请求返回
	g.Use(middlewares.Options)
	//安全
	g.Use(middlewares.Secure)
	//路由找不到提示
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})
	//健康监测
	g.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	approuter.InitAppRouter(g)

}
