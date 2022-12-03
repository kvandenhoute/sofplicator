package replication

import (
	log "github.com/sirupsen/logrus"
)

func Start(replicationInfo *ReplicationInfo) {
	log.Debug(replicationInfo)
	StartImageReplication(replicationInfo)
	StartChartReplication(replicationInfo)
}
