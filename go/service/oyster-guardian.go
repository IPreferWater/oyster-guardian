package service

import (
	"encoding/json"
	"time"

	"github.com/IPreferWater/oyster-guardian/model"
	log "github.com/sirupsen/logrus"
)

func HandleTopicDetected(payload string) {
	var detected model.Detected
	if err := json.Unmarshal([]byte(payload), &detected); err != nil {
		log.Error(err)
		return
	}

	typeDetected := getTypeDetected(detected.ImageUrl)
	threatPoints, err := getThreatPoints(detected, typeDetected)
	if err != nil {
		log.Error(err)
		return
	}
	if threatPoints > 20 {
		Stream.PublishTopicThreat(payload)
	}
	Repository.InsertDetected(detected)
}

func HandleTopicThreat(payload string) {
	log.Infof("ALERT ! There might be a threat on your farm on the location => %s", payload)
}

func getThreatPoints(detected model.Detected, typeDetected model.TypeDetected) (int, error) {
	points := 0
	switch typeDetected {
	case model.Car, model.Truck:
		//TODO vehicules
		break
	case model.Human, model.HumanPack:
		t, err := time.Parse("2006-01-02 15:04:05", detected.Timestamp)
		if err != nil {
			return points, err
		}

		// Get the hour from the time.Time value.
		hour := t.Hour()
		// having people in this time slot is strange
		if hour > 22 || hour < 6 {
			points += 100
		}else {
		// people without abilitation is strange, but could be hikers at this time slot
			points += 30
		}
	}

	return points, nil
}
