package service

import (
	"github.com/IPreferWater/oyster-guardian/model"
	log "github.com/sirupsen/logrus"
)

func getTypeDetected(url string) model.TypeDetected {
	log.Infof("analysing img with url %s", url)
	// get type detected
	log.Infof("type detected is %s", model.HumanPack)

	return model.HumanPack
}
