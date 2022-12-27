package service

import (
	"encoding/json"

	"github.com/IPreferWater/oyster-guardian/db"
	"github.com/IPreferWater/oyster-guardian/model"
	log "github.com/sirupsen/logrus"
)

func PayloadFromInput(payload string) {
	var detected model.Detected
	if err := json.Unmarshal([]byte(payload), &detected); err != nil {
		log.Error(err)
		return 
	}

	log.Infof("payload %s received", payload)
	db.Repo.InsertDetected(detected)
}
