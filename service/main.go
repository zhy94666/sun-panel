package main

import (
	"log"
	"os"
	"sun-panel/global"
	"sun-panel/initialize"
	"sun-panel/router"
)

func main() {
	err := initialize.InitApp()
	if err != nil {
		log.Println("初始化错误:", err.Error())
		panic(err)
	}

	ssl_data, env_status := os.LookupEnv("HTTPS")

	if env_status {
		if ssl_data == "true" || ssl_data == "True"{
			// Https部分
			httpsPort := global.Config.GetValueStringOrDefault("base", "https_port")
			if err := router.InitSSLRouters(":" + httpsPort); err != nil {
				panic(err)
			}
		} else {
			httpPort := global.Config.GetValueStringOrDefault("base", "http_port")
			if err := router.InitRouters(":" + httpPort); err != nil {
				panic(err)
			}
		}
	} else {
		httpPort := global.Config.GetValueStringOrDefault("base", "http_port")
		if err := router.InitRouters(":" + httpPort); err != nil {
			panic(err)
		}
	}

}
