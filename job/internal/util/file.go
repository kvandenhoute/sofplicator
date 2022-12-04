package util

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kvandenhoute/sofplicator/job/internal/model"
	log "github.com/sirupsen/logrus"
)

func readFileToByteArray(path string) []byte {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Error(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()

	return byteValue
}

func ReadArtifacts(path string) []model.Artifact {
	byteValue := readFileToByteArray(path)
	var artifacts []model.Artifact

	json.Unmarshal(byteValue, &artifacts)

	return artifacts
}
