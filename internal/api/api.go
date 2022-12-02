package api

import (
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

	jobId := replicator.StartReplication(replication)

	c.JSON(200, gin.H{"jobId": jobId})
}

func StartGlobalReplication(c *gin.Context) {
	var replication replicator.Replication

	if err := c.BindJSON(&replication); err != nil {
		log.Error(err)
		return
	}

	jobIds := replicator.StartGlobalReplication(replication)

	c.JSON(200, gin.H{"jobIds": jobIds})
}
