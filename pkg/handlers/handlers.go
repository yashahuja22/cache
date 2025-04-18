package handlers

import (
	"encoding/json"
	"gcache/pkg/config"
	"gcache/pkg/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server interface {
	SetHandler(c *gin.Context)
	GetHandler(c *gin.Context)
	DelHandler(c *gin.Context)
}

type handlers struct {
	store store.DataStore
}

func (s *handlers) SetHandler(c *gin.Context) {
	var req struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if !json.Valid(req.Value) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON value"})
		return
	}

	s.store.JSONSet(req.Key, req.Value)

	config.Logger.Log.Sugar().Infof("data stored successfully, key: %s, value: %s", req.Key, string(req.Value))

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *handlers) GetHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing key"})
		return
	}
	value, ok := s.store.JSONGet(key)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	c.Data(http.StatusOK, "application/json", value)
}

func (s *handlers) DelHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing key"})
		return
	}
	if s.store.JSONDel(key) {
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
	}
}

func NewServer(store store.DataStore) Server {
	return &handlers{
		store: store,
	}
}
