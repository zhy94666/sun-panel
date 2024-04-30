package router

import (
	"sun-panel/global"
	// "sun-panel/router/admin"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	"sun-panel/router/system"
)

// 初始化总路由
func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

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

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.Run(addr)
}

func InitSSLRouters(addr string) error {
	router := gin.Default()
	router.Use(httpsHandler())
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")
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

	global.Logger.Info("Sun-Panel is Started.  Listening and serving HTTP on ", addr)
	return router.RunTLS(addr, "crt", "key")
}

func httpsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddle := secure.New(secure.Options{
			SSLRedirect: true, //只允许https请求
			//SSLHost:"" //http到https的重定向
			STSSeconds:           1536000, //Strict-Transport-Security header的时效:1年
			STSIncludeSubdomains: true,    //includeSubdomains will be appended to the Strict-Transport-Security header
			STSPreload:           true,    //STS Preload(预加载)
			FrameDeny:            true,    //X-Frame-Options 有三个值:DENY（表示该页面不允许在 frame 中展示，即便是在相同域名的页面中嵌套也不允许）、SAMEORIGIN、ALLOW-FROM uri
			ContentTypeNosniff:   true,    //禁用浏览器的类型猜测行为,防止基于 MIME 类型混淆的攻击
			BrowserXssFilter:     true,    //启用XSS保护,并在检查到XSS攻击时，停止渲染页面
			//IsDevelopment:true,  //开发模式
		})
		err := secureMiddle.Process(context.Writer, context.Request)
		// 如果不安全，终止.
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, "数据不安全")
			return
		}
		// 如果是重定向，终止
		if status := context.Writer.Status(); status > 300 && status < 399 {
			context.Abort()
			return
		}
		context.Next()
	}
}
