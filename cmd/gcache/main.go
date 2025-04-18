package main

import (
	"gcache/pkg/config"
	"gcache/pkg/routes"
	"os"
)

func main() {
	if !config.LoadCommon() {
		return
	}

	r := routes.SetUpRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	r.Run(":" + port)
}
