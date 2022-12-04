package replication

import (
	"github.com/kvandenhoute/sofplicator/job/internal/model"
	log "github.com/sirupsen/logrus"
)

func Start(replicationInfo *model.ReplicationInfo) {
	log.Debug(replicationInfo)
	StartImageReplication(replicationInfo)
	StartChartReplication(replicationInfo)
}
