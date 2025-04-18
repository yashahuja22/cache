package routes

import (
	"gcache/pkg/handlers"
	"gcache/pkg/store"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes() *gin.Engine {
	r := gin.Default()
	store := store.NewStorageManager()
	s := handlers.NewServer(store)

	r.POST("/set", s.SetHandler)
	r.GET("/get", s.GetHandler)
	r.DELETE("/del", s.DelHandler)

	return r
}
