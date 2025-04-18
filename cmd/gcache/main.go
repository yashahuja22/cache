package main

import (
	"fmt"
	"gcache/pkg/config"
	"gcache/pkg/routes"
)

func main() {
	fmt.Println("Starting gcache..")

	if !config.LoadCommon() {
		return
	}

	config.Logger.Log.Info("gcache started..")

	r := routes.SetUpRoutes()

	r.Run(":" + config.Port)
}
