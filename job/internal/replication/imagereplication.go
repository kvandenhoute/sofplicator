package replication

import (
	"fmt"

	"github.com/kvandenhoute/sofplicator/job/internal/model"
	"github.com/kvandenhoute/sofplicator/job/internal/util"
	log "github.com/sirupsen/logrus"
)

func StartImageReplication(replicationInfo *model.ReplicationInfo) {
	log.Debug(replicationInfo)

	artifacts := util.ReadArtifacts(replicationInfo.Location.Images)

	for _, artifact := range artifacts {
		log.Info(artifact)
		err := replicate(&replicationInfo.Source, &replicationInfo.Target, &artifact)
		if !replicationInfo.ContinueOnError && err != nil {
			log.Fatal(err)
		}
	}
}

func replicate(source *model.Source, target *model.Target, artifact *model.Artifact) error {
	log.Info(source)

	// TODO https://github.com/containers/image

	return fmt.Errorf("Could not replicate %v from %v to %v", artifact, source, target)
}
