package router

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	// "sun-panel/router/admin"
	_ "net/http"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	"sun-panel/router/system"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// http协议初始化
func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")
	public_router(addr, router, routerGroup)
	return router.Run(addr)
}

// https协议初始化
func InitSSLRouters(addr string) error {
	router := gin.Default()
	router.Use(httpsHandler())
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")
	public_router(addr, router, routerGroup)
	return router.RunTLS(addr, "./conf/server.crt", "./conf/server.key")
}

// 公共路由部分
func public_router(addr string, router *gin.Engine, routerGroup *gin.RouterGroup) {
	// 接口
	system.Init(routerGroup)
	panel.Init(routerGroup)
	openness.Init(routerGroup)

	// WEB文件服务
	{
		webPath := "./web"
		router.StaticFile("/", webPath+"/index.html")
		router.Static("/assets", webPath+"/assets")
		router.Static("/custom", webPath+"/custom")
		router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		router.StaticFile("/favicon.svg", webPath+"/favicon.svg")
	}

	// 上传的文件
	sourcePath := global.Config.GetValueString("base", "source_path")
	router.Static(sourcePath[1:], sourcePath)

	global.Logger.Info("Sun-Panel is Started.", addr)
}

func httpsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddle := secure.New(secure.Options{
			SSLRedirect: true, //只允许https请求
			//SSLHost:"" //http到https的重定向
			STSSeconds:           31536000,
			STSIncludeSubdomains: true,
			STSPreload:           true,
			FrameDeny:            true,
			ContentTypeNosniff:   true,
			BrowserXssFilter:     true,
			//IsDevelopment:true,  //开发模式
		})
		err := secureMiddle.Process(context.Writer, context.Request)
		// 如果不安全，终止.
		if err != nil {
			apiReturn.ErrorByCode(context, 1400)
			return
		}
		// 如果是重定向，终止
		if status := context.Writer.Status(); status > 300 && status < 399 {
			apiReturn.ErrorByCode(context, 1400)
			context.Abort()
			return
		}
		context.Next()
	}
}
