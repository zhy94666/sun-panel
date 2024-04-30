package main

import (
	"log"
	_ "sun-panel/global"
	"sun-panel/initialize"
	"sun-panel/router"
)

func main() {
	err := initialize.InitApp()
	if err != nil {
		log.Println("初始化错误:", err.Error())
		panic(err)
	}
	// httpPort := global.Config.GetValueStringOrDefault("base", "http_port")

	// if err := router.InitRouters(":" + httpPort); err != nil {
	// 	panic(err)
	// }
	// Https部分
	if err := router.InitSSLRouters(":443"); err != nil {
		panic(err)
	}
}
