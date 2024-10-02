package main

import (
	_ "app/init"
	"app/router"
	"fmt"
	utils "framework/utils/common"
)

func main() {
	ver := utils.GetEnv("BUILD_VER", "")
	appname := utils.GetEnv("APP_NAME", "")

	fmt.Println("Service: ", appname)
	fmt.Println("Version: ", ver)
	router.Start()
}
