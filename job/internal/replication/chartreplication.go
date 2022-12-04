package replication

import (
	"github.com/kvandenhoute/sofplicator/job/internal/model"
	log "github.com/sirupsen/logrus"
)

func StartChartReplication(replicationInfo *model.ReplicationInfo) {
	log.Debug(replicationInfo)

	// TODO https://pkg.go.dev/helm.sh/helm/v3#section-readme

}
