package router

import (
	"github.com/gin-gonic/gin"
	"github.com/kvandenhoute/sofplicator/internal/api"
)

func ServeRouter() {

	router := gin.Default()
	router.POST("/startReplication", func(c *gin.Context) {
		api.StartReplication(c)
	})
	router.POST("/startGlobalReplication", func(c *gin.Context) {
		api.StartGlobalReplication(c)
	})
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "I'm alive!"})
	})

	router.Run("0.0.0.0:8080")
}
