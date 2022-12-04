package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	replicator "github.com/kvandenhoute/sofplicator/internal/replicator"
	log "github.com/sirupsen/logrus"
)

func StartReplication(c *gin.Context) {
	var replication replicator.Replication

	if err := c.BindJSON(&replication); err != nil {
		log.Error(err)
		return
	}

	jobId, err := replicator.StartReplication(replication)
	if err != nil {
		c.JSON(400, gin.H{"bad request": err})
	}

	c.JSON(200, gin.H{"jobId": jobId})
}

func StartGlobalReplication(c *gin.Context) {
	var replication replicator.Replication

	if err := c.BindJSON(&replication); err != nil {
		log.Error(err)
		return
	}

	jobIds, err := replicator.StartGlobalReplication(replication)
	if err != nil {
		c.JSON(400, gin.H{"bad request": err})
	}

	c.JSON(200, gin.H{"jobIds": jobIds})
}

func CleanReplication(uuid string, c *gin.Context) {
	err := replicator.CleanReplication(uuid)
	if err != nil {
		c.JSON(400, gin.H{"bad request": fmt.Sprint(err)})
		return
	}
	c.JSON(200, gin.H{"successful": true})
}
